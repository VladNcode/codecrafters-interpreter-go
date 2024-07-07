package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	exitCode := 0 // Default exit code

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	rawFileContents, err := os.ReadFile(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	tokenToType := map[string]string{
		"(":  "LEFT_PAREN",
		")":  "RIGHT_PAREN",
		"{":  "LEFT_BRACE",
		"}":  "RIGHT_BRACE",
		",":  "COMMA",
		".":  "DOT",
		"-":  "MINUS",
		"+":  "PLUS",
		";":  "SEMICOLON",
		"*":  "STAR",
		"/":  "SLASH",
		"=":  "EQUAL",
		"!":  "BANG",
		"<":  "LESS",
		">":  "GREATER",
		" ":  "SPACE",
		"\t": "TAB",
	}

	doublesToType := map[string]string{
		"==": "EQUAL_EQUAL",
		"!=": "BANG_EQUAL",
		"<=": "LESS_EQUAL",
		">=": "GREATER_EQUAL",
		"//": "Comment",
	}

	whiteSpace := map[string]bool{
		" ":  true,
		"\t": true,
	}

	operators := map[string]bool{
		"=": true,
		"!": true,
		"<": true,
		">": true,
		"/": true,
	}

	lines := strings.Split(string(rawFileContents), "\n")

	for lineNumber, line := range lines {
		shouldSkip := false

		for idx := 0; idx < len(line); idx++ {
			token := string(line[idx])

			if tokenType, ok := tokenToType[token]; !ok {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", lineNumber+1, token)
				exitCode = 65
			} else {
				if whiteSpace[token] {
					continue
				}

				if operators[token] && idx+2 <= len(line) && doublesToType[line[idx:idx+2]] != "" {
					token = line[idx : idx+2]
					tokenType = doublesToType[token]

					if tokenType == "Comment" {
						shouldSkip = true
						break
					}

					idx += 1
				}

				fmt.Printf("%s %s null\n", tokenType, token)
			}
		}

		if shouldSkip {
			continue
		}
	}

	fmt.Printf("EOF  null")
	os.Exit(exitCode)
}
