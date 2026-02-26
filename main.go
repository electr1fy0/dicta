package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	apiURL         = "https://api.dictionaryapi.dev/api/v2/entries/en/"
	maxDefinitions = 3
	maxTerms       = 5
	maxWidth       = 100
	minBoxWidth    = 40
	requestTimeout = 8 * time.Second
)

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

type APIError struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: "233", Dark: "212"}).
			Background(lipgloss.AdaptiveColor{Light: "153", Dark: "57"}).
			Padding(0, 2).
			MarginBottom(1)

	phoneticStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "239", Dark: "246"}).
			Italic(true)

	partOfSpeechStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.AdaptiveColor{Light: "125", Dark: "205"}).
				MarginTop(1)

	definitionNumStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "25", Dark: "39"}).
				Bold(true)

	definitionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "234", Dark: "255"})

	exampleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "241", Dark: "243"}).
			Italic(true).
			PaddingLeft(3)

	synonymStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "28", Dark: "114"})

	antonymStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "160", Dark: "203"})

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "238", Dark: "240"}).
			Bold(true)

	footerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "240", Dark: "241"}).
			Italic(true).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "124", Dark: "196"}).
			Bold(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "239", Dark: "246"})

	moreStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "238", Dark: "240"})

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.AdaptiveColor{Light: "67", Dark: "63"}).
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

func inputWord(args []string) string {
	return strings.TrimSpace(strings.Join(args, " "))
}

func fetchWordData(client *http.Client, word string) (Word, error) {
	reqURL := apiURL + url.QueryEscape(word)
	resp, err := client.Get(reqURL)
	if err != nil {
		return Word{}, fmt.Errorf("failed to fetch word: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Word{}, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(data, &apiErr); err == nil && apiErr.Message != "" {
			return Word{}, errors.New(apiErr.Message)
		}
		if resp.StatusCode == http.StatusNotFound {
			return Word{}, errors.New("word not found")
		}
		return Word{}, fmt.Errorf("dictionary API returned status %d", resp.StatusCode)
	}

	var wordData []Word
	if err := json.Unmarshal(data, &wordData); err != nil || len(wordData) == 0 {
		return Word{}, errors.New("failed to parse API response")
	}
	return wordData[0], nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "dicta <word>",
		Short: "A stylish command-line dictionary",
		Long:  "A stylish command-line dictionary that fetches definitions, pronunciations, and examples.",
		Example: strings.Join([]string{
			"dicta serendipity",
			"dicta ephemeral",
			`dicta "ice cream"`,
		}, "\n"),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			word := inputWord(args)
			if word == "" {
				return errors.New("please provide a word")
			}

			httpClient := &http.Client{Timeout: requestTimeout}
			w, err := fetchWordData(httpClient, word)
			if err != nil {
				return err
			}
			renderWord(w)
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.SetErrPrefix("")
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, errorStyle.Render(err.Error()))
		fmt.Fprintln(os.Stderr, hintStyle.Render("Example: dicta serendipity"))
		os.Exit(1)
	}
}

func renderWord(w Word) {
	var content strings.Builder

	titleCaser := cases.Title(language.English)
	content.WriteString(titleStyle.Render(titleCaser.String(w.Word)))
	content.WriteString("\n")

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
	if width < minBoxWidth {
		width = minBoxWidth
	}

	styledBox := boxStyle.Width(width - 6)
	fmt.Println(styledBox.Render(content.String()))
}
