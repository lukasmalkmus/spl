package scanner_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/lukasmalkmus/spl/internal/app/spl/scanner"
	"github.com/lukasmalkmus/spl/internal/app/spl/token"
)

func TestScanner_Scan(t *testing.T) {
	tests := []struct {
		str  string
		tok  token.Token
		lit  string
		line int
	}{
		// Special tokens
		{"!", token.ILLEGAL, "!", 1},
		{"_", token.ILLEGAL, "_", 1},
		{"_x", token.ILLEGAL, "_", 1},      // Underscore can't prefix identifier
		{"foo_", token.ILLEGAL, "foo_", 1}, // Underscore can't suffix identifier
		{"_123", token.ILLEGAL, "_", 1},    // Underscore can't prefix integer
		{"1foo", token.ILLEGAL, "1foo", 1}, // Integer can't suffix identifier
		{".", token.ILLEGAL, ".", 1},
		{".123", token.ILLEGAL, ".", 1},      // Dot can't prefix integer.
		{"123x", token.ILLEGAL, "123x", 1},   // Illegal integer (wrong hex representation)
		{"0xx08", token.ILLEGAL, "0xx08", 1}, // Illegal hex syntax
		{"''", token.ILLEGAL, "''", 1},       // Illegal integer (ASCII value) syntax
		{"'''", token.ILLEGAL, "'''", 1},     // Illegal integer (ASCII value) syntax
		{"'\\'", token.ILLEGAL, "'\\'", 1},   // Illegal integer (ASCII control character value) syntax
		{" x", token.IDENT, "x", 1},
		{"\nx", token.IDENT, "x", 2},
		{"", token.EOF, "", 1},
		{" ", token.EOF, "", 1},
		{"   ", token.EOF, "", 1},
		{"\t", token.EOF, "", 1},
		{"\n", token.EOF, "", 2},       // Single newline (LF)
		{"\r\n", token.EOF, "", 2},     // Single newline (CRLF)
		{"\n\n", token.EOF, "", 3},     // Double newline (LF + LF)
		{"\r\n\r\n", token.EOF, "", 3}, // Double newline (CRLF + CRLF)
		{"//", token.COMMENT, "//", 1},
		{"// This is a comment!", token.COMMENT, "// This is a comment!", 1},

		// Literals
		{"x", token.IDENT, "x", 1},
		{"foo ", token.IDENT, "foo", 1},
		{"foo_bar", token.IDENT, "foo_bar", 1},
		{"foo1", token.IDENT, "foo1", 1},
		{"foo_1", token.IDENT, "foo_1", 1},
		{"8", token.INT, "8", 1},
		{"64", token.INT, "64", 1},
		{"128", token.INT, "128", 1},
		{"1234", token.INT, "1234", 1},
		{"0x1a2f3F4e", token.INT, "0x1a2f3F4e", 1}, // Hex
		{"'a'", token.INT, "'a'", 1},               // ASCII character interpreted as number
		{`'\n'`, token.INT, `'\n'`, 1},             // ASCII control character interpreted as number

		// Operators and delimiters
		{"+", token.ADD, "+", 1},
		{"+4", token.ADD, "+", 1},
		{"-", token.SUB, "-", 1},
		{"-4", token.SUB, "-", 1},
		{"*", token.MUL, "*", 1},
		{"*4", token.MUL, "*", 1},
		{"/", token.QUO, "/", 1},
		{"/4", token.QUO, "/", 1},

		{"=", token.EQL, "=", 1},
		{"<", token.LSS, "<", 1},
		{">", token.GTR, ">", 1},
		{"#", token.NOT, "#", 1},

		{"<=", token.LEQ, "<=", 1},
		{">=", token.GEQ, ">=", 1},
		{":=", token.ASSIGN, ":=", 1},

		{"(", token.LPAREN, "(", 1},
		{"[", token.LBRACK, "[", 1},
		{"{", token.LBRACE, "{", 1},
		{")", token.RPAREN, ")", 1},
		{"]", token.RBRACK, "]", 1},
		{"}", token.RBRACE, "}", 1},

		{",", token.COMMA, ",", 1},
		{":", token.COLON, ":", 1},
		{";", token.SEMICOLON, ";", 1},

		// Keywords
		{"array", token.ARRAY, "array", 1},
		{"else", token.ELSE, "else", 1},
		{"if", token.IF, "if", 1},
		{"of", token.OF, "of", 1},
		{"proc", token.PROC, "proc", 1},
		{"ref", token.REF, "ref", 1},
		{"type", token.TYPE, "type", 1},
		{"var", token.VAR, "var", 1},
		{"while", token.WHILE, "while", 1},
	}

	for _, tt := range tests {
		_ = t.Run(tt.str, func(t *testing.T) {
			s := scanner.New(strings.NewReader(tt.str))
			tok, lit, pos := s.Scan()
			equals(t, tok.String(), tt.tok.String())
			equals(t, lit, tt.lit)
			equals(t, pos.Line, tt.line)
		})
	}
}

func TestScanner_ScanFullValidProgram(t *testing.T) {
	f, err := os.Open("../testdata/valid.spl")
	if err != nil {
		t.Fatal("failed to open testdata:", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal("failed to read testdata:", err)
	}

	s := scanner.New(bytes.NewReader(b))
	var count int
	for tok, _, _ := s.Scan(); tok != token.EOF; tok, _, _ = s.Scan() {
		count++
	}
	const expectedTokCount = 415
	equals(t, count, expectedTokCount)
}

// equals fails the test if got is not equal to want.
func equals(tb testing.TB, got, want interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(got, want) {
		tb.Errorf("\033[31m\n\n\tgot: %#v\n\n\twant: %#v\033[39m\n\n", got, want)
	}
}
