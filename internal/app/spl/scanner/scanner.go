package scanner

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"

	"github.com/lukasmalkmus/spl/internal/app/spl/token"
)

var eof = rune(0)

// Scanner represents a lexical scanner which tokenizes source code.
type Scanner struct {
	r                *bufio.Reader
	pos              token.Position
	resetColumnCount bool
}

// New returns a new Scanner instance which reads from the given reader.
func New(r io.Reader) *Scanner {
	return &Scanner{
		r:   bufio.NewReader(r),
		pos: token.Position{Filename: "", Line: 1, Column: 0},
	}
}

// NewFileScanner returns a new Scanner instance, but will exclusively take an
// *os.File as argument instead of the more general io.Reader interface.
// Therefore it will enhance token positions with the filename.
func NewFileScanner(f *os.File) *Scanner {
	return &Scanner{
		r:   bufio.NewReader(f),
		pos: token.Position{Filename: f.Name(), Line: 1, Column: 0},
	}
}

// Scan scans the next token and returns the token itself, its literal and its
// position in the source code. The source end is indicated by token.EOF.
func (s *Scanner) Scan() (token.Token, string, token.Position) {
	s.skipWhitespace()
	ch, pos := s.read()

	// If we see a letter consume as an ident or reserved word.
	// If we see a digit consume as an integer.
	// If we see a "'" consume as an integer as well, but us a specialized
	// scanning method.
	if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	} else if isDigit(ch) {
		s.unread()
		return s.scanInteger()
	} else if ch == '\'' {
		s.unread()
		return s.scanSpecialInteger()
	}

	// Otherwise tokenize the individual characters. No match results in an
	// illegal token.
	switch ch {
	case eof:
		pos.Column--
		return token.EOF, "", pos
	case '+':
		return token.ADD, string(ch), pos
	case '-':
		return token.SUB, string(ch), pos
	case '*':
		return token.MUL, string(ch), pos
	case '/':
		if pch := s.peek(); pch == '/' {
			return s.scanComment()
		}
		return token.QUO, string(ch), pos
	case '=':
		return token.EQL, string(ch), pos
	case '<':
		if pch := s.peek(); pch == '=' {
			_, _ = s.read()
			return token.LEQ, string(ch) + string(pch), pos
		}
		return token.LSS, string(ch), pos
	case '>':
		if pch := s.peek(); pch == '=' {
			_, _ = s.read()
			return token.GEQ, string(ch) + string(pch), pos
		}
		return token.GTR, string(ch), pos
	case '#':
		return token.NOT, string(ch), pos
	case ':':
		if pch := s.peek(); pch == '=' {
			_, _ = s.read()
			return token.ASSIGN, string(ch) + string(pch), pos
		}
		return token.COLON, string(ch), pos
	case '(':
		return token.LPAREN, string(ch), pos
	case '[':
		return token.LBRACK, string(ch), pos
	case '{':
		return token.LBRACE, string(ch), pos
	case ')':
		return token.RPAREN, string(ch), pos
	case ']':
		return token.RBRACK, string(ch), pos
	case '}':
		return token.RBRACE, string(ch), pos
	case ',':
		return token.COMMA, string(ch), pos
	case ';':
		return token.SEMICOLON, string(ch), pos
	}
	return token.ILLEGAL, string(ch), pos
}

// scanComment consumes the current rune and all contiguous comment runes.
func (s *Scanner) scanComment() (token.Token, string, token.Position) {
	// Create a buffer for the comments text. It is initially populated with a
	// slash which is the first slash of the comment token.
	var buf bytes.Buffer
	_ = buf.WriteByte('/')
	ch, pos := s.read()
	_, _ = buf.WriteRune(ch)

	// Read every subsequent character into the buffer. Newline or EOF will
	// cause the loop to exit.
	for {
		if ch, _ := s.read(); isNewline(ch) {
			s.unread()
			break
		} else if ch == eof {
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}
	return token.COMMENT, buf.String(), pos
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (token.Token, string, token.Position) {
	var buf bytes.Buffer
	ch, pos := s.read()
	_, _ = buf.WriteRune(ch)

	// Read every subsequent ident character into the buffer. Non-ident
	// characters and EOF will cause the loop to exit.
	for {
		if ch, _ := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// Make sure the last character is not an underscore, which is illegal.
	if ch := buf.Bytes()[buf.Len()-1]; ch == '_' {
		return token.ILLEGAL, buf.String(), pos
	}
	return token.Lookup(buf.String()), buf.String(), pos
}

// scanInteger consumes the current rune and all contiguous integer runes.
func (s *Scanner) scanInteger() (token.Token, string, token.Position) {
	var buf bytes.Buffer
	ch, pos := s.read()
	_, _ = buf.WriteRune(ch)

	// Read every subsequent character into the buffer. Non-integer characters
	// and EOF will cause the loop to exit.
	for {
		if ch, _ := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// Uppercase 'X' not allowed in hexadecimal representation.
	if bytes.ContainsRune(buf.Bytes(), 'X') {
		return token.ILLEGAL, buf.String(), pos
	} else if _, err := strconv.ParseInt(buf.String(), 0, 32); err != nil {
		return token.ILLEGAL, buf.String(), pos
	}
	return token.INT, buf.String(), pos
}

// scanSpecialInteger consumes the current rune and all contiguous special
// integer runes.
func (s *Scanner) scanSpecialInteger() (token.Token, string, token.Position) {
	var buf bytes.Buffer
	ch, pos := s.read()
	_, _ = buf.WriteRune(ch)

	// Read every subsequent character into the buffer. Non-integer characters,
	// multiple whitespaces and EOF will cause the loop to exit.
	var (
		charCount uint8 = 1 // We already saw one.
		sqmCount  uint8 = 1 // We already saw one.
		wsCount   uint8
	)
	for {
		if ch, _ := s.read(); ch == eof {
			break
		} else if (sqmCount == 2 && charCount >= 3) || sqmCount > 2 || wsCount > 2 {
			s.unread()
			break
		} else {
			charCount++
			if ch == '\'' {
				sqmCount++
			} else if ch == ' ' {
				wsCount++
			}
			_, _ = buf.WriteRune(ch)
		}
	}
	b := buf.Bytes()

	// The first character is a tick so the last one must be one, too.
	if l := len(b); l < 3 || b[l-1] != '\'' {
		return token.ILLEGAL, buf.String(), pos
	}

	// If the length of the input is three, the character encapsuled by the
	// single quotation marks must be a printable ASCII character.
	// If the length of the input is four, the characters encapsuled by the
	// single quotation marks must form an ASCII escape sequence by the first
	// character being a backslash and the second character being a letter.
	if len(b) == 3 && b[1] > 31 && b[1] < 127 {
		return token.INT, buf.String(), pos
	} else if len(b) == 4 && b[1] == '\\' && isLetter(rune(b[2])) {
		return token.INT, buf.String(), pos
	}
	return token.ILLEGAL, buf.String(), pos
}

// skipWhitespace consumes the current rune and all contiguous newline and
// whitespace. It keeps track of the token position.
func (s *Scanner) skipWhitespace() {
	var buf bytes.Buffer
	for ch, _ := s.read(); isNewline(ch) || isWhitespace(ch); ch, _ = s.read() {
		if isNewline(ch) {
			_, _ = buf.WriteRune(ch)
			s.resetColumnCount = true
		}
	}
	clean := stripCR(buf.Bytes())
	s.pos.Line += len(clean)
	s.unread()
}

// read reads the next rune from the buffered reader. Returns rune(0) if an
// error occurs (which can also be io.EOF returned from the underlying reader).
func (s *Scanner) read() (rune, token.Position) {
	if s.resetColumnCount {
		s.pos.Column = 0
		s.resetColumnCount = false
	}
	s.pos.Column++
	s.pos.Char++

	// Read from the underlying reader.
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof, s.pos
	}
	return ch, s.pos
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
	s.pos.Column--
	s.pos.Char--
}

// unread peeks for the next rune from the buffered reader.
func (s *Scanner) peek() rune {
	ch, _ := s.read()
	s.unread()
	return ch
}

// isWhitespace returns true if the rune is a space or tab.
func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' }

// isNewline returns true if the rune is a newline.
func isNewline(ch rune) bool { return ch == '\n' || ch == '\r' }

// isLetter returns true if the rune is a letter.
func isLetter(ch rune) bool { return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') }

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

// stripCR removes every carriage-return from a slice of bytes, effectively
// turning a CRLF into a LF.
func stripCR(b []byte) []byte {
	c := make([]byte, 0)
	for _, ch := range b {
		if ch == '\n' {
			c = append(c, ch)
		}
	}
	return c
}
