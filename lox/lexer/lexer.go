package lexer

type Lexer struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int

	error LexerError
}

func NewScanner(source string) *Lexer {
	return &Lexer{source: source, start: 0, current: 0, line: 1}
}

func (t *Lexer) ScanTokens() ([]Token, LexerError) {
	for !t.isAtEnd() {
		// We are at the beginning of the next lexeme
		t.start = t.current
		t.scanToken()
		if t.error.HasError {
			return t.tokens, t.error
		}
	}
	t.tokens = append(t.tokens, *NewToken(EOF, "", "", t.line))
	return t.tokens, t.error
}

func (t *Lexer) scanToken() {
	c := t.advance()
	switch c {
	case "(":
		t.addToken(LEFT_PAREN)
		break
	case ")":
		t.addToken(RIGHT_PAREN)
		break
	case "{":
		t.addToken(LEFT_BRACE)
		break
	case "}":
		t.addToken(RIGHT_BRACE)
		break
	case ",":
		t.addToken(COMMA)
		break
	case ".":
		t.addToken(DOT)
		break
	case "-":
		t.addToken(MINUS)
		break
	case "+":
		t.addToken(PLUS)
		break
	case ";":
		t.addToken(SEMICOLON)
		break
	case "*":
		t.addToken(STAR)
		break

	case "!":
		if t.match("=") {
			t.addToken(BANG_EQUAL)
		} else {
			t.addToken(BANG)
		}
		break
	case "=":
		if t.match("=") {
			t.addToken(EQUAL_EQUAL)
		} else {
			t.addToken(EQUAL)
		}
		break
	case "<":
		if t.match("=") {
			t.addToken(LESS_EQUAL)
		} else {
			t.addToken(LESS)
		}
		break
	case ">":
		if t.match("=") {
			t.addToken(GREATER_EQUAL)
		} else {
			t.addToken(GREATER)
		}
		break

	case "/":
		if t.match("/") {
			// A comment goes until the end of the line
			for t.peek() != "\n" && !t.isAtEnd() {
				t.advance()
			}
		} else if t.match("*") {
			for !t.isAtEnd() {
				if t.peek() == "*" {
					t.advance()
					if t.match("/") {
						break
					}
				} else {
					t.advance()
				}
			}
		} else {
			t.addToken(SLASH)
		}
		break

	case " ":
		break
	case "\r":
		break
	case "\t":
		// Ignore whitespace
		break
	case "\n":
		t.line++
		break

	case `"`:
		t.string()
		break

	default:
		if isDigit(c) {
			t.number()
		} else if isAlpha(c) {
			t.identifier()
		} else {
			t.error = LexerError{true, t.line, "Unexpected character."}
			return
		}
		break
	}

}

func (t *Lexer) advance() string {
	t.current++
	return string(t.source[t.current-1])
}

func (t *Lexer) addToken(type0 TokenType) {
	t.addTokenWithLiteral(type0, "")
}

func (t *Lexer) addTokenWithLiteral(type0 TokenType, literal string) {
	text := t.source[t.start:t.current]
	t.tokens = append(t.tokens, *NewToken(type0, text, literal, t.line))
}

func (t *Lexer) isAtEnd() bool {
	return t.current >= len(t.source)
}

func (t *Lexer) match(expected string) bool {
	if t.isAtEnd() {
		return false
	}
	if string(t.source[t.current]) != expected {
		return false
	}
	t.current++
	return true
}

func (t *Lexer) peek() string {
	if t.isAtEnd() {
		return "\\0"
	}
	return string(t.source[t.current])
}

func (t *Lexer) string() {
	for t.peek() != `"` && !t.isAtEnd() {
		if t.peek() == "\n" {
			t.line++
		}
		t.advance()
	}

	if t.isAtEnd() {
		t.error = LexerError{true, t.line, "Unterminated string"}
		return
	}

	// The closing ".
	t.advance()

	// Trim the surrounding quotes.
	value := t.source[t.start+1 : t.current-1]
	t.addTokenWithLiteral(STRING, value)
}

func (t *Lexer) number() {
	for isDigit(t.peek()) {
		t.advance()
	}
	if t.peek() == "." && isDigit(t.peekNext()) {
		// Consume the "."
		t.advance()
		for isDigit(t.peek()) {
			t.advance()
		}
	}
	//double, _ := strconv.ParseFloat(t.source[t.start:t.current], 64)
	t.addTokenWithLiteral(NUMBER, t.source[t.start:t.current])
}

func (t *Lexer) peekNext() string {
	if t.current+1 >= len(t.source) {
		return `\0`
	}

	return string(t.source[t.current+1])
}

func (t *Lexer) identifier() {
	for isAlphaNumeric(t.peek()) {
		t.advance()
	}
	text := t.source[t.start:t.current]
	type0, ok := KeyWords[text]
	if !ok {
		type0 = IDENTIFIER
	}
	t.addToken(type0)
}

func isDigit(c string) bool {
	return c >= "0" && c <= "9"
}

func isAlpha(c string) bool {
	return (c >= "a" && c <= "z") ||
		(c >= "A" && c <= "Z") ||
		c == "_"
}

func isAlphaNumeric(c string) bool {
	return isAlpha(c) || isDigit(c)
}
