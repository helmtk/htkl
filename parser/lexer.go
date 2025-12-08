package parser

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenIdent
	TokenString
	TokenNumber
	TokenColon
	TokenComma
	TokenLBrace
	TokenRBrace
	TokenLBracket
	TokenRBracket
	TokenLParen
	TokenRParen
	TokenComment
	TokenNewline
	TokenIf
	TokenElse
	TokenFor
	TokenIn
	TokenWith
	TokenAs
	TokenDo
	TokenEnd
	TokenBreak
	TokenContinue
	TokenLet
	TokenDefine
	TokenInclude
	TokenSpread
	TokenTrue
	TokenFalse
	TokenNull
	TokenDot       // .
	TokenAssign    // =
	TokenPlus      // +
	TokenMinus     // -
	TokenMul       // *
	TokenDiv       // /
	TokenPipe      // |
	TokenAnd       // &&
	TokenOr        // ||
	TokenEq        // ==
	TokenNeq       // !=
	TokenNot       // !
	TokenLt        // <
	TokenLte       // <=
	TokenGt        // >
	TokenGte       // >=
)

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenIdent:
		return "identifier"
	case TokenString:
		return "string"
	case TokenNumber:
		return "number"
	case TokenColon:
		return "':'"
	case TokenComma:
		return "','"
	case TokenLBrace:
		return "'{'"
	case TokenRBrace:
		return "'}'"
	case TokenLBracket:
		return "'['"
	case TokenRBracket:
		return "']'"
	case TokenLParen:
		return "'('"
	case TokenRParen:
		return "')'"
	case TokenComment:
		return "comment"
	case TokenNewline:
		return "newline"
	case TokenIf:
		return "'if'"
	case TokenElse:
		return "'else'"
	case TokenFor:
		return "'for'"
	case TokenIn:
		return "'in'"
	case TokenWith:
		return "'with'"
	case TokenAs:
		return "'as'"
	case TokenDo:
		return "'do'"
	case TokenEnd:
		return "'end'"
	case TokenBreak:
		return "'break'"
	case TokenContinue:
		return "'continue'"
	case TokenLet:
		return "'let'"
	case TokenDefine:
		return "'define'"
	case TokenInclude:
		return "'include'"
	case TokenSpread:
		return "'spread'"
	case TokenTrue:
		return "'true'"
	case TokenFalse:
		return "'false'"
	case TokenNull:
		return "'null'"
	case TokenDot:
		return "'.'"
	case TokenAssign:
		return "'='"
	case TokenPlus:
		return "'+'"
	case TokenMinus:
		return "'-'"
	case TokenMul:
		return "'*'"
	case TokenDiv:
		return "'/'"
	case TokenPipe:
		return "'|'"
	case TokenAnd:
		return "'&&'"
	case TokenOr:
		return "'||'"
	case TokenEq:
		return "'=='"
	case TokenNeq:
		return "'!='"
	case TokenNot:
		return "'!'"
	case TokenLt:
		return "'<'"
	case TokenLte:
		return "'<='"
	case TokenGt:
		return "'>'"
	case TokenGte:
		return "'>='"
	default:
		return fmt.Sprintf("unknown(%d)", t)
	}
}

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

type Lexer struct {
	input string
	pos   int
	line  int
	col   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
		pos:   0,
		line:  1,
		col:   1,
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF, Line: l.line, Col: l.col}
	}

	ch := l.current()

	// Newlines
	if ch == '\n' {
		token := Token{Type: TokenNewline, Line: l.line, Col: l.col, Value: "\n"}
		l.advance()
		return token
	}

	// Comments
	if ch == '#' {
		return l.readComment()
	}

	// Strings
	if ch == '"' {
		// Check for multiline string (triple quotes)
		if l.peek() == '"' && l.peekN(2) == '"' {
			return l.readMultilineString()
		}
		return l.readString()
	}

	// Numbers
	if unicode.IsDigit(rune(ch)) || (ch == '-' && l.peek() != 0 && unicode.IsDigit(rune(l.peek()))) {
		return l.readNumber()
	}

	// Identifiers and keywords
	if unicode.IsLetter(rune(ch)) || ch == '_' {
		return l.readIdentifier()
	}

	// Operators and single-character tokens
	token := Token{Line: l.line, Col: l.col}
	switch ch {
	case ':':
		token.Type = TokenColon
		token.Value = ":"
		l.advance()
	case ',':
		token.Type = TokenComma
		token.Value = ","
		l.advance()
	case '{':
		token.Type = TokenLBrace
		token.Value = "{"
		l.advance()
	case '}':
		token.Type = TokenRBrace
		token.Value = "}"
		l.advance()
	case '[':
		token.Type = TokenLBracket
		token.Value = "["
		l.advance()
	case ']':
		token.Type = TokenRBracket
		token.Value = "]"
		l.advance()
	case '(':
		token.Type = TokenLParen
		token.Value = "("
		l.advance()
	case ')':
		token.Type = TokenRParen
		token.Value = ")"
		l.advance()
	case '.':
		token.Type = TokenDot
		token.Value = "."
		l.advance()
	case '&':
		if l.peek() == '&' {
			token.Type = TokenAnd
			token.Value = "&&"
			l.advance()
			l.advance()
		} else {
			token.Type = TokenEOF
			token.Value = string(ch)
			l.advance()
		}
	case '|':
		if l.peek() == '|' {
			token.Type = TokenOr
			token.Value = "||"
			l.advance()
			l.advance()
		} else {
			token.Type = TokenPipe
			token.Value = "|"
			l.advance()
		}
	case '=':
		if l.peek() == '=' {
			token.Type = TokenEq
			token.Value = "=="
			l.advance()
			l.advance()
		} else {
			token.Type = TokenAssign
			token.Value = "="
			l.advance()
		}
	case '!':
		if l.peek() == '=' {
			token.Type = TokenNeq
			token.Value = "!="
			l.advance()
			l.advance()
		} else {
			token.Type = TokenNot
			token.Value = "!"
			l.advance()
		}
	case '<':
		if l.peek() == '=' {
			token.Type = TokenLte
			token.Value = "<="
			l.advance()
			l.advance()
		} else {
			token.Type = TokenLt
			token.Value = "<"
			l.advance()
		}
	case '>':
		if l.peek() == '=' {
			token.Type = TokenGte
			token.Value = ">="
			l.advance()
			l.advance()
		} else {
			token.Type = TokenGt
			token.Value = ">"
			l.advance()
		}
	case '+':
		token.Type = TokenPlus
		token.Value = "+"
		l.advance()
	case '-':
		token.Type = TokenMinus
		token.Value = "-"
		l.advance()
	case '*':
		token.Type = TokenMul
		token.Value = "*"
		l.advance()
	case '/':
		token.Type = TokenDiv
		token.Value = "/"
		l.advance()
	default:
		token.Type = TokenEOF
		token.Value = string(ch)
		l.advance()
	}

	return token
}

func (l *Lexer) current() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) peek() byte {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

func (l *Lexer) peekN(n int) byte {
	if l.pos+n >= len(l.input) {
		return 0
	}
	return l.input[l.pos+n]
}

func (l *Lexer) advance() {
	if l.pos < len(l.input) {
		if l.input[l.pos] == '\n' {
			l.line++
			l.col = 1
		} else {
			l.col++
		}
		l.pos++
	}
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		ch := l.current()
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.advance()
		} else {
			break
		}
	}
}

func (l *Lexer) readComment() Token {
	start := l.pos
	startCol := l.col
	l.advance() // skip #

	for l.pos < len(l.input) && l.current() != '\n' {
		l.advance()
	}

	return Token{
		Type:  TokenComment,
		Value: l.input[start:l.pos],
		Line:  l.line,
		Col:   startCol,
	}
}

// unescapeString processes escape sequences in a string.
// Note: \$ is handled specially to prevent interpolation - it's kept as a
// marker (\x00$) that will be replaced by $ after interpolation processing.
func unescapeString(s string) string {
	if !strings.Contains(s, "\\") {
		return s
	}

	var result strings.Builder
	result.Grow(len(s))

	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				result.WriteByte('\n')
				i++
			case 't':
				result.WriteByte('\t')
				i++
			case 'r':
				result.WriteByte('\r')
				i++
			case '"':
				result.WriteByte('"')
				i++
			case '\\':
				result.WriteByte('\\')
				i++
			case '$':
				// Escaped $ - use null byte as marker to prevent interpolation
				result.WriteByte('\x00')
				result.WriteByte('$')
				i++
			default:
				// Unknown escape sequence, keep the backslash
				result.WriteByte('\\')
			}
		} else {
			result.WriteByte(s[i])
		}
	}

	return result.String()
}

func (l *Lexer) readString() Token {
	start := l.pos
	startCol := l.col
	l.advance() // skip opening "

	for l.pos < len(l.input) && l.current() != '"' {
		if l.current() == '\\' {
			l.advance() // skip escape
		}
		l.advance()
	}

	if l.pos < len(l.input) {
		l.advance() // skip closing "
	}

	// Remove quotes from value and unescape
	value := unescapeString(l.input[start+1 : l.pos-1])

	return Token{
		Type:  TokenString,
		Value: value,
		Line:  l.line,
		Col:   startCol,
	}
}

func (l *Lexer) readMultilineString() Token {
	start := l.pos
	startCol := l.col
	startLine := l.line

	// Skip opening """
	l.advance()
	l.advance()
	l.advance()

	// Read until we find closing """
	for l.pos < len(l.input) {
		if l.current() == '"' && l.peek() == '"' && l.peekN(2) == '"' {
			// Found closing """
			value := unescapeString(l.input[start+3 : l.pos])
			// Skip closing """
			l.advance()
			l.advance()
			l.advance()
			return Token{
				Type:  TokenString,
				Value: value,
				Line:  startLine,
				Col:   startCol,
			}
		}
		l.advance()
	}

	// If we get here, we didn't find closing """
	value := unescapeString(l.input[start+3:])
	return Token{
		Type:  TokenString,
		Value: value,
		Line:  startLine,
		Col:   startCol,
	}
}

func (l *Lexer) readNumber() Token {
	start := l.pos
	startCol := l.col

	if l.current() == '-' {
		l.advance()
	}

	for l.pos < len(l.input) && (unicode.IsDigit(rune(l.current())) || l.current() == '.') {
		l.advance()
	}

	return Token{
		Type:  TokenNumber,
		Value: l.input[start:l.pos],
		Line:  l.line,
		Col:   startCol,
	}
}

func (l *Lexer) readIdentifier() Token {
	start := l.pos
	startCol := l.col

	for l.pos < len(l.input) {
		ch := l.current()
		if unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' {
			l.advance()
		} else {
			break
		}
	}

	value := l.input[start:l.pos]
	tokenType := TokenIdent

	// Check for keywords
	switch value {
	case "if":
		tokenType = TokenIf
	case "else":
		tokenType = TokenElse
	case "for":
		tokenType = TokenFor
	case "in":
		tokenType = TokenIn
	case "with":
		tokenType = TokenWith
	case "as":
		tokenType = TokenAs
	case "do":
		tokenType = TokenDo
	case "end":
		tokenType = TokenEnd
	case "break":
		tokenType = TokenBreak
	case "continue":
		tokenType = TokenContinue
	case "let":
		tokenType = TokenLet
	case "define":
		tokenType = TokenDefine
	case "include":
		tokenType = TokenInclude
	case "spread":
		tokenType = TokenSpread
	case "true":
		tokenType = TokenTrue
	case "false":
		tokenType = TokenFalse
	case "null":
		tokenType = TokenNull
	}

	return Token{
		Type:  tokenType,
		Value: value,
		Line:  l.line,
		Col:   startCol,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%v(%q)", t.Type, t.Value)
}
