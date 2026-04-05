package parser

import (
	"errors"
)

type ReservedKeyword string

const (
	If         ReservedKeyword = "If"
	Else       ReservedKeyword = "Else"
	ElseIf     ReservedKeyword = "ElseIf"
	For        ReservedKeyword = "For"
	Return     ReservedKeyword = "Return"
	Function   ReservedKeyword = "Function"
	In         ReservedKeyword = "In"
	IntType    ReservedKeyword = "IntType"
	StringType ReservedKeyword = "StringType"
	BoolType   ReservedKeyword = "BoolType"
	FloatType  ReservedKeyword = "FloatType"
)

func getKeywords() map[string]ReservedKeyword {
	return map[string]ReservedKeyword{
		"if":     If,
		"else":   Else,
		"elseif": ElseIf,
		"for":    For,
		"return": Return,
		"fn":     Function,
		"in":     In,
		"int":    IntType,
		"string": StringType,
		"bool":   BoolType,
		"float":  FloatType,
	}
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// If the character is a whitespace character, it is used to separate tokens.
func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// Checks if a string is recognized as a keyword in the language.
func isKeyword(s string) bool {
	_, exists := getKeywords()[s]
	return exists
}

// Double quotes are used for string literals.
func isDoubleQuote(c byte) bool {
	return c == '"'
}

type Token struct {
	Type  string
	Value string
}

// Handles single-character operators.
func handleSingleOperators(c byte) Token {
	switch c {
	case '(':
		return Token{Type: "OpenParen", Value: "("}
	case ')':
		return Token{Type: "CloseParen", Value: ")"}
	case '{':
		return Token{Type: "OpenCurly", Value: "{"}
	case '}':
		return Token{Type: "CloseCurly", Value: "}"}
	case '[':
		return Token{Type: "OpenSquare", Value: "["}
	case ']':
		return Token{Type: "CloseSquare", Value: "]"}
	case '+':
		return Token{Type: "Plus", Value: "+"}
	case '-':
		return Token{Type: "Minus", Value: "-"}
	case '*':
		return Token{Type: "Asterisk", Value: "*"}
	case '/':
		return Token{Type: "Slash", Value: "/"}
	case '%':
		return Token{Type: "Percent", Value: "%"}
	case '<':
		return Token{Type: "LessThan", Value: "<"}
	case '>':
		return Token{Type: "GreaterThan", Value: ">"}
	case '=':
		return Token{Type: "Equals", Value: "="}
	case '!':
		return Token{Type: "ExclamationMark", Value: "!"}
	case '?':
		return Token{Type: "QuestionMark", Value: "?"}
	case ';':
		return Token{Type: "Semicolon", Value: ";"}
	case ':':
		return Token{Type: "Colon", Value: ":"}
	case ',':
		return Token{Type: "Comma", Value: ","}
	case '.':
		return Token{Type: "Period", Value: "."}
	}

	return Token{}
}

// Handles double-character operators like '==', '!=', '<=', '>=', '&&', '||', '++', '--'.
func handleDoubleOperators(c byte, next byte) Token {
	if c == '=' && next == '=' {
		return Token{Type: "EqualsEquals", Value: "=="}
	} else if c == '!' && next == '=' {
		return Token{Type: "NotEquals", Value: "!="}
	} else if c == '<' && next == '=' {
		return Token{Type: "LessThanOrEqual", Value: "<="}
	} else if c == '>' && next == '=' {
		return Token{Type: "GreaterThanOrEqual", Value: ">="}
	} else if c == '<' && next == '<' {
		return Token{Type: "LeftShift", Value: "<<"}
	} else if c == '>' && next == '>' {
		return Token{Type: "RightShift", Value: ">>"}
	} else if c == '&' && next == '&' {
		return Token{Type: "LogicalAnd", Value: "&&"}
	} else if c == '|' && next == '|' {
		return Token{Type: "LogicalOr", Value: "||"}
	} else if c == '+' && next == '+' {
		return Token{Type: "Increment", Value: "++"}
	} else if c == '-' && next == '-' {
		return Token{Type: "Decrement", Value: "--"}
	} else if c == '/' && next == '/' {
		return Token{Type: "LineComment", Value: "//"}
	} else if c == '/' && next == '*' {
		return Token{Type: "BlockCommentStart", Value: "/*"}
	} else if c == '*' && next == '/' {
		return Token{Type: "BlockCommentEnd", Value: "*/"}
	} else if c == '-' && next == '>' {
		return Token{Type: "ArrowRight", Value: "->"}
	} else if c == '<' && next == '-' {
		return Token{Type: "ArrowLeft", Value: "<-"}
	}

	return Token{}
}

func tokenizeInput(input string) ([]Token, error) {
	var tokens []Token

	var isInQuote bool = false
	var currentString string = ""
	var currentNumber string = ""
	var currentToken string = ""
	var isInLineComment bool = false
	var isInBlockComment bool = false

	flushCurrent := func() {
		if currentToken != "" {
			if isKeyword(currentToken) {
				tokens = append(tokens, Token{Type: string(getKeywords()[currentToken]), Value: currentToken})
			} else {
				tokens = append(tokens, Token{Type: "Identifier", Value: currentToken})
			}
			currentToken = ""
		} else if currentNumber != "" {
			tokens = append(tokens, Token{Type: "NumberLiteral", Value: currentNumber})
			currentNumber = ""
		} else if currentString != "" {
			tokens = append(tokens, Token{Type: "StringLiteral", Value: currentString})
			currentString = ""
		}
	}

	for i := 0; i < len(input); i++ {
		c := input[i]

		if isInLineComment {
			if c == '\n' {
				isInLineComment = false
			}
			continue
		}

		if isInBlockComment {
			if i+1 < len(input) {
				t := handleDoubleOperators(c, input[i+1])
				if t.Type == "BlockCommentEnd" {
					isInBlockComment = false
					tokens = append(tokens, t)
					i++ // Skip the next character since it's part of the block comment end operator.
					continue
				}
			}

			continue
		}

		// First we handle string literals because they can contain characters that would otherwise be treated as operators or identifiers.
		if isDoubleQuote(c) {
			isInQuote = !isInQuote

			if !isInQuote {
				flushCurrent()
				currentString = ""
			}
			continue
		}

		if isInQuote {
			currentString += string(c)
			continue
		}

		if i+1 < len(input) {
			next := input[i+1]
			t := handleDoubleOperators(c, next)
			if t.Type != "" {
				switch t.Type {
				case "LineComment":
					isInLineComment = true
				case "BlockCommentStart":
					isInBlockComment = true
				default:
					flushCurrent()
				}

				tokens = append(tokens, t)
				i++ // Skip the next character since it's part of the double operator.
				continue
			}
		}

		// Then, every other operator is handled.
		t := handleSingleOperators(c)
		if t.Type != "" {
			flushCurrent()
			tokens = append(tokens, t)
			continue
		}

		// Handling numbers easily by casting them into a string.
		if isDigit(c) {
			currentNumber += string(c)
			continue
		} else {
			if currentNumber != "" {
				flushCurrent()
			}
		}

		// Every other identifier is handled by checking if it's a keyword or not.
		if isAlpha(c) {
			currentToken += string(c)
			continue
		} else {
			flushCurrent()
		}

		if isWhitespace(c) {
			continue
		}

		return nil, errors.New("Unrecognized character: " + string(c))
	}

	if isInQuote {
		return nil, errors.New("Unterminated string literal")
	}

	flushCurrent()

	if len(tokens) == 0 {
		return nil, errors.New("No tokens found in input")
	}

	tokens = append(tokens, Token{Type: "EOF", Value: "EndOfFile"})

	return tokens, nil
}
