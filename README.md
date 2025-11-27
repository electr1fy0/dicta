# dicta

A stylish command-line dictionary tool that fetches definitions, pronunciations, and examples.


## Install

### From source

```bash
git clone https://github.com/electr1fy0/dicta.git
cd dicta
go build -o dicta
```

## Usage

```bash
dicta serendipity
dicta ephemeral
dicta hello
```

**Output includes:**
- Phonetic pronunciation
- Part of speech (noun, verb, etc.)
- Definitions with examples
- Synonyms and antonyms (when available)

## Notes

- Requires an internet connection to fetch definitions.
- Data is sourced from `api.dictionaryapi.dev`.
