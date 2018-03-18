/*
Copyright 2018 The Elasticshift Authors.
*/
package scanner

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"gitlab.com/conspico/elasticshift/pkg/shiftfile/token"
)

// eof represents a marker rune for the end of the reader.
const eof = rune(0)

// The New func takes error func as argument and passes
// the control to handle the error
type ErrorFunc func(pos token.Position, msg string)

// Scanner for ShiftFile to read tokens
// It takes io.Reader or []byte as source which can then be
// tokenized through scan function.
type Scanner struct {
	token         token.Token
	buffer        *bytes.Buffer
	src           []byte
	pos           token.Position
	prevPos       token.Position
	ch            rune
	Error         ErrorFunc
	incrementLine bool
}

// NewScanner...
// Creates a new scanner from a reader
func New(src []byte, err ErrorFunc) *Scanner {

	s := &Scanner{buffer: bytes.NewBuffer(src)}
	s.Error = err
	s.src = src
	s.ch = ' '
	s.pos.Offset = 0
	s.pos.Line = 1
	s.pos.Column = 0
	return s
}

// Returns the next rune from the scanner
func (s *Scanner) next() rune {

	s.prevPos = s.pos

	c, size, err := s.buffer.ReadRune()
	if err != nil {
		s.ch = eof
		s.pos.Column++
		s.pos.Offset += size
		return eof
	}

	s.ch = c
	s.pos.Column++
	s.pos.Offset += size

	if s.incrementLine && c == '\n' {
		s.pos.Column = 1
		s.pos.Line++
	}

	// If we see a null character with data left, then that is an error
	if c == '\x00' && len(s.src) > 0 {
		s.error("unexpected null character (0x00)")
	}

	// fmt.Println(fmt.Sprintf("Line %d, Column %d, Offset %d, Rune %s", s.pos.Line, s.pos.Column, s.pos.Offset, string(c)))
	return c
}

func (s *Scanner) unread() {

	err := s.buffer.UnreadRune()
	if err != nil {
		panic(err)
	}

	s.pos = s.prevPos
}

// Returns most recently parsed token
func (s *Scanner) Token() token.Token {
	return s.token
}

// HasMoreTokens indicate whether there are more token to scan, if true
// otherwise false
func (s *Scanner) HasMoreTokens() bool {
	return s.peek() != eof
}

// Scan consumes the next token
func (s *Scanner) Scan() token.Token {

	// scan's the first character
	s.next()

	// skip the whitespaces
	for isWhiteSpace(s.ch) {
		s.next()
	}

	var tok token.Token

	if s.ch == eof {
		tok.Type = token.EOF
		tok.Text = "EOF"
		return tok
	}

	s.incrementLine = true

	switch ch := s.ch; {
	case isLetter(ch):
		tok.Type, tok.Text = s.scanIdentifier()
	case isDigit(ch):
		tok.Type, tok.Text = s.scanNumber(false)
	default:
		switch ch {
		case '"':
			tok.Type, tok.Text = s.scanString()
		case '-':
			tok.Type, tok.Text = s.scanCommand()
		case '/':
			s.next()
			if ch == '/' && s.ch == '/' {
				tok.Type, tok.Text = token.HINT, "//"
			} else if ch == '/' && s.ch == '*' {
				tok.Type, tok.Text = token.LHINT, "/*"
			}
		case '*':
			s.next()
			if ch == '*' && s.ch == '/' {
				tok.Type, tok.Text = token.RHINT, "*/"
			}
		case '(':
			tok.Type = token.LPAREN
			tok.Text = "("
		case ')':
			tok.Type = token.RPAREN
			tok.Text = ")"
		case '{':
			tok.Type = token.LBRACE
			tok.Text = "{"
		case '}':
			tok.Type = token.RBRACE
			tok.Text = "}"
		case ':':
			tok.Type = token.HINT_DEL
			tok.Text = ":"
		case ',':
			tok.Type = token.COMMA
			tok.Text = ","
		case '#':
			tok.Type, tok.Text = s.scanComment(ch)
		case '`':
			tok.Type, tok.Text = s.scanMultilineString()
		}

		if ch == '\n' {
			s.pos.Column = 0
		}
	}

	tok.Position = s.pos

	s.token = tok
	return tok
}

func (s *Scanner) peek() rune {

	ru, _, err := s.buffer.ReadRune()
	if err != nil && strings.EqualFold("EOF", err.Error()) {
		return eof
	}

	err = s.buffer.UnreadRune()
	if err != nil {
		panic(err)
	}
	return ru
}

func (s *Scanner) unreadIfNotEOF(ch rune) {
	if ch != eof && ch >= 0 {
		s.unread()
	}
}

func (s *Scanner) scanWhitespace() (token.Type, string) {

	ofs := s.pos.Offset - 1
	for isWhiteSpace(s.ch) {
		s.next()
	}
	return token.WHITESPACE, string(s.src[ofs:s.pos.Offset])
}

func (s *Scanner) scanMultilineString() (token.Type, string) {

	ofs := s.pos.Offset - 1

	for {

		ch := s.ch
		if ch < 0 {
			s.err("raw string not terminated")
			break
		}

		if ch == '`' {
			break
		}
	}

	return token.STRING, string(s.src[ofs:s.pos.Offset])
}

func (s *Scanner) scanCommand() (token.Type, string) {

	ofs := s.pos.Offset

	cont := true
	ch := s.ch
	for cont {

		ch = s.ch
		s.next()

		if ch == '\\' && s.ch == '\n' {
			cont = true
		} else {
			cont = s.ch != '\n'
		}
	}

	s.unreadIfNotEOF(ch)

	cmd := strings.TrimSpace(string(s.src[ofs:s.pos.Offset]))
	return token.COMMAND, cmd //s.stripEscapeChars(cmd)
}

func (s *Scanner) stripEscapeChars(param string) string {
	data := []byte(param)
	buf := make([]byte, len(data))
	i := 0
	for _, ch := range data {
		if ch == '\r' || ch == '\t' || ch == '\\' || ch == '\n' {
			continue
		}
		buf[i] = ch
		i++
	}
	return string(buf[:i])
}

func (s *Scanner) scanString() (token.Type, string) {

	ofs := s.pos.Offset

	for {

		ch := s.next()
		if ch == '\n' || ch < 0 {
			s.err("string not terminated")
		}

		if ch == '\\' {
			continue
		}

		if ch == '"' {
			break
		}

		// if ch == '\\' {
		// 	s.scanEscape('""')
		// }
	}

	return token.STRING, string(s.src[ofs : s.pos.Offset-1])
}

func (s *Scanner) scanComment(ch rune) (token.Type, string) {

	ofs := s.pos.Offset - 1

	for ch != '\n' && ch >= 0 && ch != eof {
		ch = s.next()
	}

	if ch != eof && ch >= 0 {
		s.unread()
	}

	// fmt.Println(fmt.Sprintf("Scancomment offset:%d, currentoffset:%d, value:%s", ofs, s.pos.Offset, string(s.src[ofs:s.pos.Offset])))

	return token.COMMENT, string(s.src[ofs:s.pos.Offset])
}

// Scan the number by reading the next character, it could be INT or FLOATi
func (s *Scanner) scanNumber(decimalPoint bool) (token.Type, string) {

	ofs := s.pos.Offset
	var tok token.Type
	tok = token.INT

	if decimalPoint {
		ofs--
		tok = token.FLOAT
		s.scanMantissa(10)
		s.scanExponent()
	}
	return tok, string(s.src[ofs:s.pos.Offset])
}

func (s *Scanner) scanMantissa(base int) {
	for digitVal(s.ch) < base {
		s.next()
	}
}

func (s *Scanner) scanExponent() {

	if s.ch == 'e' || s.ch == 'E' {
		s.next()
		if s.ch == '-' || s.ch == '+' {
			s.next()
		}
		if digitVal(s.ch) < 10 {
			s.scanMantissa(10)
		} else {
			s.error("illegal floating-point exponent")
		}
	}
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= ch && ch <= 'f':
		return int(ch - 'a' + 10)
	case 'A' <= ch && ch <= 'F':
		return int(ch - 'A' + 10)
	}
	return 16 // larger than any legal digit val
}

// Scan the identifier by reading next character in look until reach whitespace or endline
func (s *Scanner) scanIdentifier() (token.Type, string) {

	ofs := s.pos.Offset - 1
	for isLetter(s.ch) || isDigit(s.ch) {
		s.next()
	}

	// unread if identifer read any special characters
	s.unreadIfNotEOF(s.ch)

	iden := string(s.src[ofs:s.pos.Offset])
	if len(iden) > 1 {
		tok := token.Lookup(iden)
		if tok.IsKeyword() {
			return tok, iden
		}
	}
	return token.IDENTIFIER, iden

}

// Construct the error with the given message
func (s *Scanner) err(message string) {
	s.error(message)
	fmt.Sprintf("Error while parsing shiftfile %d:%d -> %s", s.pos.Line, s.pos.Offset, message)
}

func (s *Scanner) error(msg string) {
	s.Error(s.pos, msg)
}

// Returns true if the rune is a whitespace, otherwise false
func isWhiteSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// Returns true if rune is a character type, otherwise false
func isLetter(c rune) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' || c >= utf8.RuneSelf && unicode.IsLetter(c)
}

// Returns true if rune is a digit, otherwise false
func isDigit(c rune) bool {
	return '0' <= c && c <= '9' || c >= utf8.RuneSelf && unicode.IsDigit(c)
}

// Returns true if rune is a newline character, otherwise false
func isNewline(c rune) bool {
	return c == '\n'
}
