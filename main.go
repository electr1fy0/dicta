package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

const BASE_API = "https://api.dictionaryapi.dev/api/v2/entries/en/"
const ITALIC = "\033[3m"
const RESET = "\033[0m"

var rootCmd = &cobra.Command{
	Short: "a cli dictionary tool",
	Use:   "ditto <word>",
}

type Word struct {
	Word     string    `json:"word"`
	Meanings []Meaning `json:"meanings"`
}

type Meaning struct {
	PartOfSpeech string       `json:"partOfSpeech"`
	Definitions  []Definition `json:"definitions"`
}

type Definition struct {
	Def string `json:"definition"`
	Syn string `json:"synonyms"`
}

var meaningCmd = &cobra.Command{
	Short: "Get the meaning of a word",
	Use:   "meaning <word>",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		word := args[0]
		url := BASE_API + word
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		data, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}
		var wordData []Word
		json.Unmarshal(data, &wordData)
		fmt.Println("word:", wordData[0].Word)
		fmt.Println("part of speech:", wordData[0].Meanings[0].PartOfSpeech)
		fmt.Println("definition:", wordData[0].Meanings[0].Definitions[0].Def)

		fmt.Printf("\n%ssourced from: api.dictionaryapi.dev%s\n", ITALIC, RESET)
	},
}

func main() {
	rootCmd.AddCommand(meaningCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
