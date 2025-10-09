package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

const (
	BASE_API = "https://api.dictionaryapi.dev/api/v2/entries/en/"
	ITALIC   = "\033[3m"
	BOLD     = "\033[1m"
	RESET    = "\033[0m"
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

var rootCmd = &cobra.Command{
	Use:   "dicta",
	Short: "A CLI dictionary tool",
}

var meaningCmd = &cobra.Command{
	Use:   "meaning <word>",
	Short: "Get the meaning of a word",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		word := args[0]
		url := BASE_API + word
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal("Failed to fetch word:", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Fatalf("Word not found or API error: %s", resp.Status)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var wordData []Word
		if err := json.Unmarshal(data, &wordData); err != nil {
			log.Fatal("Failed to parse API response:", err)
		}

		w := wordData[0]
		fmt.Printf("%sWord:%s %s\n\n", BOLD, RESET, w.Word)

		if len(w.Phonetics) > 0 {
			fmt.Printf("%sPronunciation:%s\n", BOLD, RESET)
			for _, p := range w.Phonetics {
				if p.Text != "" {
					fmt.Printf("  %s\n", p.Text)
				}
				if p.Audio != "" {
					fmt.Printf("  Audio: %s\n", p.Audio)
				}
			}
			fmt.Println()
		}

		for i, m := range w.Meanings {
			fmt.Printf("%sMeaning %d (%s):%s\n", BOLD, i+1, m.PartOfSpeech, RESET)
			for j, d := range m.Definitions {
				fmt.Printf("  %d. %s\n", j+1, d.Def)
				if d.Example != "" {
					fmt.Printf("     Example: %s\n", d.Example)
				}
				if len(d.Synonyms) > 0 {
					fmt.Printf("     Synonyms: %s\n", strings.Join(d.Synonyms, ", "))
				}
				if len(d.Antonyms) > 0 {
					fmt.Printf("     Antonyms: %s\n", strings.Join(d.Antonyms, ", "))
				}
			}
			fmt.Println()
		}

		fmt.Printf("%sSourced from: api.dictionaryapi.dev%s\n", ITALIC, RESET)
	},
}

func main() {
	rootCmd.AddCommand(meaningCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
