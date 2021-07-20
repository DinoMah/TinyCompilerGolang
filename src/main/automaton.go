package main

import (
	"regexp"
	"strings"
)

type info struct {
	tokenType string
	token     string
}

func analyze(line string) []info {
	var tokens []info
	state := ""
	actualToken := ""
	tokenType := ""
	actualPos := 0
	for state != "terminated" {
		actualToken, tokenType = getToken(line, &state, &actualPos)
		if isReserved(actualToken) {
			tokenType = "reserved word"
		}
		if strings.Compare(tokenType, "") != 0 && strings.Compare(actualToken, "") != 0 {
			tokens = append(tokens, info{actualToken, tokenType})
		}
	}
	return tokens
}

func isReserved(token string) bool {
	reservedWords := [...]string{"main", "if", "then", "else", "end", "do", "while", "repeat", "until", "read", "write", "float", "integer", "bool"}
	for i := 0; i < len(reservedWords); i++ {
		if strings.Compare(token, reservedWords[i]) == 0 {
			return true
		}
	}
	return false
}

func getToken(line string, state *string, actualPos *int) (string, string) {
	letterChar := regexp.MustCompile(`[\wñÑ]`)
	numChar := regexp.MustCompile(`[0-9]`)
	*state = "init"
	token := ""
	tokenType := ""
	actualLetter := ""
	discard := false
	for *state != "done" {
		actualLetter = getNextLetter(line, actualPos)
		if actualLetter == "" {
			discard = true
			tokenType = getLastToken(*state, token)
			*state = "terminated"
			break
		}
		switch *state {
		case "init":
			*state, discard = getActualState(actualLetter)
		case "id":
			if letterChar.MatchString(actualLetter) {
				discard = false
			} else {
				*state = "done"
				discard = true
				tokenType = "id"
				*actualPos--
			}
		case "num":
			if numChar.MatchString(actualLetter) {
				discard = false
			} else if actualLetter == "." {
				discard = false
				*state = "num real"
			} else if strings.Compare(actualLetter, " ") == 0 || strings.Compare(actualLetter, "\t") == 0 || strings.Compare(actualLetter, "\n") == 0 || isSymbol(actualLetter) {
				*state = "done"
				discard = true
				if strings.Contains(token, ".") {
					tokenType = "real"
				} else {
					tokenType = "integer"
				}
				*actualPos--
			} else {
				*state = "error"
				discard = false
			}
		case "num real":
			if numChar.MatchString(actualLetter) {
				discard = false
			} else if actualLetter == " " || actualLetter == "\n" || actualLetter == "\t" {
				discard = true
				*actualPos--
				*state = "done"
				tokenType = "real"
			} else {
				discard = false
				*state = "error"
			}
		case "string":
			if strings.Compare(actualLetter, "\"") == 0 {
				discard = false
				*state = "done"
				tokenType = "string"
			} else {
				discard = false
			}
		case "incompleteToken":
			if strings.Compare(actualLetter, "=") == 0 {
				discard = false
				*state = "done"
				if strings.Contains(token, ":") {
					tokenType = "asignment"
				} else if strings.Contains(token, "!") {
					tokenType = "different"
				} else {
					tokenType = "equal"
				}
			} else {
				discard = true
				*actualPos--
				*state = "error"
			}
		case "div":
			if strings.Compare(actualLetter, "/") == 0 {
				discard = false
				*state = "singleComment"
			} else if strings.Compare(actualLetter, "*") == 0 {
				discard = false
				*state = "multilineComment"
			} else {
				discard = true
				*actualPos--
				*state = "done"
				tokenType = "division"
			}
		case "singleComment":
			if strings.Compare(actualLetter, "\n") == 0 {
				discard = true
				*state = "done"
				tokenType = "single comment"
			}
		case "multilineComment":
			if strings.Compare(actualLetter, "*") == 0 {
				discard = false
				*state = "endingMultiComment"
			} else {
				discard = false
			}
		case "endingMultiComment":
			if strings.Compare(actualLetter, "/") == 0 {
				discard = false
				*state = "done"
				tokenType = "multi line comment"
			} else {
				discard = false
				*state = "multilineComment"
			}
		case "sum":
			if strings.Compare(actualLetter, "+") == 0 {
				discard = false
				*state = "done"
				tokenType = "increment by one"
			} else {
				discard = true
				*actualPos--
				*state = "done"
				tokenType = "sum"
			}
		case "res":
			if strings.Compare(actualLetter, "-") == 0 {
				discard = false
				*state = "done"
				tokenType = "decrement by one"
			} else {
				discard = true
				*actualPos--
				*state = "done"
				tokenType = "rest"
			}
		case "less":
			if strings.Compare(actualLetter, "=") == 0 {
				discard = false
				*state = "done"
				tokenType = "less or equal"
			} else {
				discard = true
				*actualPos--
				*state = "done"
				tokenType = "less"
			}
		case "greater":
			if strings.Compare(actualLetter, "=") == 0 {
				discard = false
				*state = "done"
				tokenType = "greater or equal"
			} else {
				discard = true
				*actualPos--
				*state = "done"
				tokenType = "greater"
			}
		case "other":
			discard = true
			*actualPos--
			*state = "done"
			tokenType = "symbol"
		case "blankSpace":
			*state = "done"
			discard = true
			*actualPos--
			tokenType = ""
		default:
			if actualLetter == " " || actualLetter == "\t" || actualLetter == "\n" {
				*state = "done"
				discard = true
				*actualPos--
				tokenType = "error"
			} else {
				discard = false
			}
		}
		if !discard {
			token += actualLetter
		}
	}
	return token, tokenType
}

func getNextLetter(line string, i *int) string {
	if *i < len(line) {
		letter := string(line[*i])
		*i++
		return letter
	}
	return ""
}

func getActualState(actualLetter string) (string, bool) {
	letterChar := regexp.MustCompile(`[a-zA-ZñÑ]`)
	numChar := regexp.MustCompile(`[0-9]`)
	switch {
	case letterChar.MatchString(actualLetter):
		return "id", false
	case numChar.MatchString(actualLetter):
		return "num", false
	case strings.Compare(actualLetter, ".") == 0:
		return "num real", false
	case strings.Compare(actualLetter, "\"") == 0:
		return "string", false
	case actualLetter == " ", actualLetter == "\t", actualLetter == "\n":
		return "blankSpace", true
	case strings.Compare(actualLetter, ":") == 0, strings.Compare(actualLetter, "!") == 0, strings.Compare(actualLetter, "=") == 0:
		return "incompleteToken", false
	case strings.Compare(actualLetter, "/") == 0:
		return "div", false
	case strings.Compare(actualLetter, "+") == 0:
		return "sum", false
	case strings.Compare(actualLetter, "-") == 0:
		return "res", false
	case strings.Compare(actualLetter, "<") == 0:
		return "less", false
	case strings.Compare(actualLetter, ">") == 0:
		return "greater", false
	case actualLetter == "%", actualLetter == "(", actualLetter == ")", actualLetter == "{", actualLetter == "}", actualLetter == "*", actualLetter == ";", actualLetter == ",":
		return "other", false
	default:
		return "error", false
	}
}

func isSymbol(letter string) bool {
	switch letter {
	case "(", ")", "{", "}", "+", "-", "/", "%", "*", ">", "<", ">=", "<=", "!=", "==", ",", ";":
		return true
	default:
		return false
	}
}

func getLastToken(state string, token string) string {
	switch state {
	case "id":
		return "id"
	case "num":
		if strings.Contains(token, ".") {
			return "real"
		}
		return "integer"
	case "string":
		return "string"
	case "incompleteToken":
		return "error"
	case "singleComment":
		return "single comment"
	case "multilineComment":
		return "multi line comment"
	case "div":
		return "division"
	case "sum":
		return "sum"
	case "res":
		return "rest"
	case "less":
		return "less"
	case "greater":
		return "greater"
	case "other":
		return "symbol"
	default:
		return "error"
	}
}
