package calculator

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType mengidentifikasi jenis token formula.
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNumber
	TokenIdent
	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash
	TokenPercent
	TokenLParen
	TokenRParen
)

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenNumber:
		return "number"
	case TokenIdent:
		return "identifier"
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenStar:
		return "*"
	case TokenSlash:
		return "/"
	case TokenPercent:
		return "%"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	default:
		return "unknown"
	}
}

// Token adalah satu unit leksikal hasil tokenisasi.
type Token struct {
	Type  TokenType
	Value string
	Pos   int // posisi byte di input asli (untuk pesan error yang jelas)
}

// Lexer men-tokenisasi string formula.
type Lexer struct {
	input string
	pos   int
}

// NewLexer membuat Lexer untuk input tertentu.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// NextToken mengambil token berikutnya dari input.
func (l *Lexer) NextToken() (Token, error) {
	for l.pos < len(l.input) {
		c := l.input[l.pos]
		switch {
		case c == ' ' || c == '\t' || c == '\n' || c == '\r':
			l.pos++
		case c == '+':
			l.pos++
			return Token{Type: TokenPlus, Value: "+", Pos: l.pos - 1}, nil
		case c == '-':
			l.pos++
			return Token{Type: TokenMinus, Value: "-", Pos: l.pos - 1}, nil
		case c == '*':
			l.pos++
			return Token{Type: TokenStar, Value: "*", Pos: l.pos - 1}, nil
		case c == '/':
			l.pos++
			return Token{Type: TokenSlash, Value: "/", Pos: l.pos - 1}, nil
		case c == '%':
			l.pos++
			return Token{Type: TokenPercent, Value: "%", Pos: l.pos - 1}, nil
		case c == '(':
			l.pos++
			return Token{Type: TokenLParen, Value: "(", Pos: l.pos - 1}, nil
		case c == ')':
			l.pos++
			return Token{Type: TokenRParen, Value: ")", Pos: l.pos - 1}, nil
		case isDigit(c):
			return l.lexNumber()
		case isIdentStart(c):
			return l.lexIdent()
		default:
			return Token{}, fmt.Errorf("karakter tidak dikenal %q pada posisi %d", string(c), l.pos)
		}
	}
	return Token{Type: TokenEOF, Pos: l.pos}, nil
}

// lexNumber membaca literal angka (integer atau desimal, mis. 500000, 2.5).
func (l *Lexer) lexNumber() (Token, error) {
	start := l.pos
	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
	}
	if l.pos < len(l.input) && l.input[l.pos] == '.' && l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1]) {
		l.pos++
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
	}
	return Token{Type: TokenNumber, Value: l.input[start:l.pos], Pos: start}, nil
}

// lexIdent membaca identifier variabel/komponen (huruf, angka, underscore).
func (l *Lexer) lexIdent() (Token, error) {
	start := l.pos
	for l.pos < len(l.input) && isIdentPart(l.input[l.pos]) {
		l.pos++
	}
	return Token{Type: TokenIdent, Value: l.input[start:l.pos], Pos: start}, nil
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }

func isIdentStart(c byte) bool {
	return c == '_' || unicode.IsLetter(rune(c))
}

func isIdentPart(c byte) bool {
	return c == '_' || unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c))
}

// Tokenize men-tokenisasi seluruh input menjadi slice token (untuk test/validasi).
func Tokenize(input string) ([]Token, error) {
	l := NewLexer(strings.TrimSpace(input))
	var tokens []Token
	for {
		tok, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			return tokens, nil
		}
	}
}
