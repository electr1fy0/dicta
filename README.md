# dicta

A command-line dictionary tool that fetches definitions, pronunciations, and examples from the command line.

## Install

### From source

```bash
git clone https://github.com/your-username/dicta.git
cd dicta
go build -o dicta
# Move to PATH
sudo mv dicta /usr/local/bin/
```

## Usage

`dicta` pulls word definitions, pronunciations, examples, synonyms, and antonyms from the free Dictionary API.

### Look up a word

```bash
dicta meaning serendipity
```

**Output includes:**
- Phonetic pronunciation
- Part of speech (noun, verb, etc.)
- Definitions with examples
- Synonyms and antonyms (when available)
- Audio pronunciation links

## Notes

- `dicta` requires an internet connection to fetch definitions.
- Data is sourced from `api.dictionaryapi.dev`.
