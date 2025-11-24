# lookup

A stylish command-line dictionary tool that fetches definitions, pronunciations, and examples.

This is a fork of [dicta](https://github.com/electr1fy0/dicta) by electr1fy0 - thanks for the original implementation!

## What's Different

- **Simplified usage** - just `lookup <word>` instead of `dicta meaning <word>`
- **Glamorous output** - styled with [Lipgloss](https://github.com/charmbracelet/lipgloss) featuring colors, borders, and clean formatting
- **Streamlined display** - limits to 3 definitions per part of speech, deduplicates phonetics

## Install

### From source

```bash
git clone https://github.com/mkaz/lookup.git
cd lookup
go build -o lookup
```

## Usage

```bash
lookup serendipity
lookup ephemeral
lookup hello
```

**Output includes:**
- Phonetic pronunciation
- Part of speech (noun, verb, etc.)
- Definitions with examples
- Synonyms and antonyms (when available)

## Notes

- Requires an internet connection to fetch definitions.
- Data is sourced from `api.dictionaryapi.dev`.
