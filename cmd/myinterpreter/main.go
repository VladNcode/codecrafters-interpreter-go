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

	runeToName := map[rune]string{
		'(': "LEFT_PAREN",
		')': "RIGHT_PAREN",
		'{': "LEFT_BRACE",
		'}': "RIGHT_BRACE",
		',': "COMMA",
		'.': "DOT",
		'-': "MINUS",
		'+': "PLUS",
		';': "SEMICOLON",
		'*': "STAR",
		'/': "SLASH",
		'=': "EQUAL",
	}

	lines := strings.Split(string(rawFileContents), "\n")

	for lineNumber, line := range lines {

		for idx := 0; idx < len(line); idx++ {
			current := rune(line[idx])

			if idx+2 <= len(line) && line[idx:idx+2] == "==" {
				fmt.Printf("EQUAL_EQUAL == null\n")
				idx += 1 // Skip the next character since its a part of '=='
				continue
			}

			if name, ok := runeToName[current]; ok {
				fmt.Printf("%s %c null\n", name, current)
			} else {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", lineNumber+1, current)
				exitCode = 65
			}

		}

	}

	fmt.Printf("EOF  null")
	os.Exit(exitCode)
}
