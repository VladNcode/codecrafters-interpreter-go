package main

import (
	"fmt"
	"os"
	"strings"
)

func isAlpha(c string) bool {
	return (c >= "a" && c <= "z") || (c >= "A" && c <= "Z") || c == "_"
}

func isDigit(c string) bool {
	return c >= "0" && c <= "9"
}

func isStringStart(c string) bool {
	return c == "\""
}

func isIdentifier(c string) bool {
	return isDigit(c) || isAlpha(c)
}

func removeTrailingZeroes(s string) string {
	if !strings.HasSuffix(s, "0") {
		return s
	}

	return strings.TrimRight(s, "0") + "0"
}

func tokenize(rawFileContents []byte) {
	exitCode := 0 // Default exit code

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

	reservedWords := map[string]string{
		"and":    "AND",
		"class":  "CLASS",
		"else":   "ELSE",
		"false":  "FALSE",
		"for":    "FOR",
		"fun":    "FUN",
		"if":     "IF",
		"nil":    "NIL",
		"or":     "OR",
		"print":  "PRINT",
		"return": "RETURN",
		"super":  "SUPER",
		"this":   "THIS",
		"true":   "TRUE",
		"var":    "VAR",
		"while":  "WHILE",
	}

	lines := strings.Split(string(rawFileContents), "\n")

	for lineNumber, line := range lines {
		shouldSkip := false

		for idx := 0; idx < len(line); {
			token := string(line[idx])

			if isStringStart(token) {
				stringText := token
				idx += 1

				for idx < len(line) {
					token = string(line[idx])
					stringText += token

					idx += 1

					if token == "\"" {
						break
					}
				}

				if !strings.HasSuffix(stringText, "\"") {
					fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", lineNumber+1)
					exitCode = 65
					continue
				}

				fmt.Printf("STRING %s %s\n", stringText, strings.Trim(stringText, "\""))

				continue
			} else if isDigit(token) {
				periodCounter := int(0)
				numberText := ""

				for idx < len(line) && (isDigit(string(line[idx])) || string(line[idx]) == ".") && periodCounter <= 1 {
					if string(line[idx]) == "." {
						if periodCounter == int(1) || idx == len(line)-1 {
							break
						}

						periodCounter = 1
					}

					numberText += string(line[idx])
					idx += 1
				}

				if periodCounter > 0 {
					fmt.Printf("NUMBER %s %s\n", numberText, removeTrailingZeroes(numberText))
				} else {
					fmt.Printf("NUMBER %s %s.0\n", numberText, numberText)
				}

				continue
			} else if isAlpha(token) {
				identefierText := ""

				for idx < len(line) && isIdentifier(string(line[idx])) {
					identefierText += string(line[idx])
					idx += 1
				}

				reservedWord := reservedWords[identefierText]

				if reservedWord != "" {
					fmt.Printf("%s %s null\n", reservedWord, identefierText)
				} else {
					fmt.Printf("IDENTIFIER %s null\n", identefierText)
				}

				continue
			} else if whiteSpace[token] {
				idx++
				continue
			} else if operators[token] && idx+2 <= len(line) && doublesToType[line[idx:idx+2]] != "" {
				token = line[idx : idx+2]
				tokenType := doublesToType[token]

				if tokenType == "Comment" {
					shouldSkip = true
					break
				}

				fmt.Printf("%s %s null\n", tokenType, token)
				idx += 2
				continue
			} else if tokenToType[token] != "" && token != "" {
				fmt.Printf("%s %s null\n", tokenToType[token], token)
				idx++
				continue
			} else {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", lineNumber+1, token)
				exitCode = 65
				idx++
				continue
			}
		}

		if shouldSkip {
			continue
		}
	}

	fmt.Printf("EOF  null")
	os.Exit(exitCode)
}

func main() {
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

	tokenize(rawFileContents)
}
