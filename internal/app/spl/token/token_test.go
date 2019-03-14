package token_test

import (
	"reflect"
	"testing"

	"github.com/lukasmalkmus/spl/internal/app/spl/token"
)

func TestToken(t *testing.T) {
	tests := []struct {
		str    string
		lit    string
		tok    token.Token
		isSpec bool
		isLit  bool
		isOp   bool
		isKey  bool
		isDir  bool
	}{
		// Special tokens
		{"ILLEGAL", "", token.ILLEGAL, true, false, false, false, false},
		{"EOF", "", token.EOF, true, false, false, false, false},
		{"COMMENT", "", token.COMMENT, true, false, false, false, false},

		// Identifiers and basic type literals
		{"IDENT", "", token.IDENT, false, true, false, false, false},
		{"INT", "", token.INT, false, true, false, false, false},

		// Operators and delimiters
		{"", "+", token.ADD, false, false, true, false, false},
		{"", "-", token.SUB, false, false, true, false, false},
		{"", "*", token.MUL, false, false, true, false, false},
		{"", "/", token.QUO, false, false, true, false, false},

		{"", "=", token.EQL, false, false, true, false, false},
		{"", "<", token.LSS, false, false, true, false, false},
		{"", ">", token.GTR, false, false, true, false, false},
		{"", "#", token.NOT, false, false, true, false, false},

		{"", "<=", token.LEQ, false, false, true, false, false},
		{"", ">=", token.GEQ, false, false, true, false, false},
		{"", ":=", token.ASSIGN, false, false, true, false, false},

		{"", "(", token.LPAREN, false, false, true, false, false},
		{"", "[", token.LBRACK, false, false, true, false, false},
		{"", "{", token.LBRACE, false, false, true, false, false},
		{"", ")", token.RPAREN, false, false, true, false, false},
		{"", "]", token.RBRACK, false, false, true, false, false},
		{"", "}", token.RBRACE, false, false, true, false, false},

		{"", ",", token.COMMA, false, false, true, false, false},
		{"", ":", token.COLON, false, false, true, false, false},
		{"", ";", token.SEMICOLON, false, false, true, false, false},

		// Keywords
		{"array", "array", token.ARRAY, false, false, false, true, false},
		{"else", "else", token.ELSE, false, false, false, true, false},
		{"if", "if", token.IF, false, false, false, true, false},
		{"of", "of", token.OF, false, false, false, true, false},
		{"proc", "proc", token.PROC, false, false, false, true, false},
		{"ref", "ref", token.REF, false, false, false, true, false},
		{"type", "type", token.TYPE, false, false, false, true, false},
		{"var", "var", token.VAR, false, false, false, true, false},
		{"while", "while", token.WHILE, false, false, false, true, false},
	}

	for _, tt := range tests {
		name := tt.str
		if tt.str == "" && tt.lit != "" {
			name = tt.lit
		}
		_ = t.Run(name, func(t *testing.T) {
			if tt.str != "" && tt.lit == "" {
				equals(t, tt.tok.String(), tt.str)
			} else if tt.str != "" && tt.lit != "" {
				equals(t, tt.tok.String(), tt.str)
			} else {
				equals(t, tt.tok.String(), tt.lit)

			}
			equals(t, tt.tok.IsKeyword(), tt.isKey)
			equals(t, tt.tok.IsLiteral(), tt.isLit)
			equals(t, tt.tok.IsOperator(), tt.isOp)
		})
	}
}

// TestLookup makes sure that Lookup returns either the right keyword or IDENT
// for non keywords, like directives or identifiers.
func TestLookup(t *testing.T) {
	tests := []struct {
		str   string
		isKey bool
		isDir bool
	}{
		// Identifiers and basic type literals
		{"abc", false, false},
		{"123", false, false},

		// Keywords
		{"array", true, false},
		{"else", true, false},
		{"if", true, false},
		{"of", true, false},
		{"proc", true, false},
		{"ref", true, false},
		{"type", true, false},
		{"var", true, false},
		{"while", true, false},
	}

	for _, tt := range tests {
		_ = t.Run(tt.str, func(t *testing.T) {
			tok := token.Lookup(tt.str)
			equals(t, tok.IsKeyword(), tt.isKey)
		})
	}
}

// equals fails the test if got is not equal to want.
func equals(tb testing.TB, got, want interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(got, want) {
		tb.Errorf("\033[31m\n\n\tgot: %#v\n\n\twant: %#v\033[39m\n\n", got, want)
	}
}
