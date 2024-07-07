package main

import (
	"fmt"
	"os"
	"regexp"
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

	isDigit := regexp.MustCompile(`[0-9]`)

	for lineNumber, line := range lines {
		shouldSkip := false
		stringMode := false
		stringText := ""
		numberMode := false
		numberText := ""
		float := false

		for idx := 0; idx < len(line); idx++ {
			token := string(line[idx])

			if tokenType, ok := tokenToType[token]; !ok && token != "\"" && !stringMode && !numberMode && !isDigit.MatchString(token) {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", lineNumber+1, token)
				exitCode = 65
			} else if idx == len(line)-1 && stringMode && token != "\"" {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", lineNumber+1)
				exitCode = 65
			} else {
				if token == "\"" {
					if !stringMode {
						stringMode = true
						continue
					}

					fmt.Printf("STRING \"%s\" %s\n", stringText, stringText)
					stringText = ""
					stringMode = false
					continue
				}

				if stringMode {
					stringText += token
					continue
				}

				if isDigit.MatchString(token) {
					if !numberMode {
						numberMode = true
					}

					numberText += token

					if idx != len(line)-1 {
						continue
					}
				}

				if numberMode {
					if token == "." {
						if float {
							fmt.Printf("NUMBER %s %s\n", numberText, numberText)
							numberMode = false
							numberText = ""
							float = false
							fmt.Printf("%s %s null\n", tokenType, token)
						} else {
							if idx != len(line)-1 {
								float = true
								numberText += token
							} else {
								fmt.Printf("NUMBER %s %s.0\n", numberText, numberText)
								fmt.Printf("%s %s null\n", tokenType, token)
							}
						}

						continue
					} else {
						if float {
							fmt.Printf("NUMBER %s %s\n", numberText, numberText)
						} else {
							fmt.Printf("NUMBER %s %s.0\n", numberText, numberText)
						}

						numberMode = false
						numberText = ""
						float = false

						if isDigit.MatchString(token) {
							continue
						}
					}

				}

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
