package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const (
	apiURL         = "https://api.dictionaryapi.dev/api/v2/entries/en/"
	maxDefinitions = 3
	maxTerms       = 5
	maxWidth       = 100
)

// API response types
type Word struct {
	Word      string     `json:"word"`
	Phonetics []Phonetic `json:"phonetics"`
	Meanings  []Meaning  `json:"meanings"`
}

type Phonetic struct {
	Text  string `json:"text"`
	Audio string `json:"audio"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
}

type Definition struct {
	Def      string   `json:"definition"`
	Example  string   `json:"example,omitempty"`
	Synonyms []string `json:"synonyms,omitempty"`
	Antonyms []string `json:"antonyms,omitempty"`
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Background(lipgloss.Color("57")).
			Padding(0, 2).
			MarginBottom(1)

	phoneticStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")).
			Italic(true)

	partOfSpeechStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205")).
				MarginTop(1)

	definitionNumStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Bold(true)

	definitionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	exampleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243")).
			Italic(true).
			PaddingLeft(3)

	synonymStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("114"))

	antonymStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203"))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Bold(true)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2)
)

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		return maxWidth
	}
	if width > maxWidth {
		return maxWidth
	}
	return width
}

func truncateSlice(s []string, max int) []string {
	if len(s) > max {
		return s[:max]
	}
	return s
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(errorStyle.Render("Usage: lookup <word>"))
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("246")).Render("Example: lookup serendipity"))
		os.Exit(1)
	}

	word := os.Args[1]

	resp, err := http.Get(apiURL + word)
	if err != nil {
		fmt.Println(errorStyle.Render("Failed to fetch word: " + err.Error()))
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(errorStyle.Render("Word not found"))
		os.Exit(1)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(errorStyle.Render(err.Error()))
		os.Exit(1)
	}

	var wordData []Word
	if err := json.Unmarshal(data, &wordData); err != nil || len(wordData) == 0 {
		fmt.Println(errorStyle.Render("Failed to parse API response"))
		os.Exit(1)
	}

	w := wordData[0]
	renderWord(w)
}

func renderWord(w Word) {
	var content strings.Builder

	// Title with word
	content.WriteString(titleStyle.Render(strings.Title(w.Word)))
	content.WriteString("\n")

	// Phonetics
	if len(w.Phonetics) > 0 {
		var phonetics []string
		seen := make(map[string]bool)
		for _, p := range w.Phonetics {
			if p.Text != "" && !seen[p.Text] {
				phonetics = append(phonetics, p.Text)
				seen[p.Text] = true
			}
		}
		if len(phonetics) > 0 {
			content.WriteString(phoneticStyle.Render(strings.Join(phonetics, " • ")))
			content.WriteString("\n")
		}
	}

	for _, m := range w.Meanings {
		content.WriteString("\n")
		content.WriteString(partOfSpeechStyle.Render("▸ " + m.PartOfSpeech))
		content.WriteString("\n")

		for j, d := range m.Definitions {
			if j >= maxDefinitions {
				remaining := len(m.Definitions) - maxDefinitions
				if remaining > 0 {
					moreStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
					content.WriteString(moreStyle.Render(fmt.Sprintf("  ... and %d more\n", remaining)))
				}
				break
			}

			num := definitionNumStyle.Render(fmt.Sprintf("%d.", j+1))
			def := definitionStyle.Render(d.Def)
			content.WriteString(fmt.Sprintf("  %s %s\n", num, def))

			if d.Example != "" {
				content.WriteString(exampleStyle.Render("\""+d.Example+"\"") + "\n")
			}

			if len(d.Synonyms) > 0 {
				syns := truncateSlice(d.Synonyms, maxTerms)
				content.WriteString(fmt.Sprintf("     %s %s\n",
					labelStyle.Render("syn:"),
					synonymStyle.Render(strings.Join(syns, ", "))))
			}

			if len(d.Antonyms) > 0 {
				ants := truncateSlice(d.Antonyms, maxTerms)
				content.WriteString(fmt.Sprintf("     %s %s\n",
					labelStyle.Render("ant:"),
					antonymStyle.Render(strings.Join(ants, ", "))))
			}
		}
	}

	content.WriteString(footerStyle.Render("source: api.dictionaryapi.dev"))

	width := getTerminalWidth()
	// Account for border (2) and padding (4)
	styledBox := boxStyle.Width(width - 6)
	fmt.Println(styledBox.Render(content.String()))
}
