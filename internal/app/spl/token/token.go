package token

import "strconv"

// Token is a lexical token of the simple programming language.
type Token int

// All available tokens.
const (
	// Special tokens
	ILLEGAL Token = iota // Illegal character
	EOF                  // End of file/source
	COMMENT              // Source code comment

	// Identifiers and basic type literals (these tokens stand for classes of
	// literals)
	literalBeg
	IDENT // x, y, abc, foo_bar, fooBar, FooBar, main
	INT   // 12345, 0x12aBcD, 'a', '\n'
	literalEnd

	// Operators and delimiters
	operatorBeg
	ADD // +
	SUB // -
	MUL // *
	QUO // /

	EQL // =
	LSS // <
	GTR // >
	NOT // #

	LEQ    // <=
	GEQ    // >=
	ASSIGN // :=

	LPAREN // (
	LBRACK // [
	LBRACE // {
	RPAREN // )
	RBRACK // ]
	RBRACE // }

	COMMA     // ,
	COLON     // :
	SEMICOLON // ;
	operatorEnd

	// Keywords
	keywordBeg
	ARRAY // array
	ELSE  // else
	IF    // if
	OF    // of
	PROC  // proc
	REF   // ref
	TYPE  // type
	VAR   // var
	WHILE // while
	keywordEnd
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT: "IDENT",
	INT:   "INT",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	QUO: "/",

	EQL: "=",
	LSS: "<",
	GTR: ">",
	NOT: "#",

	LEQ:    "<=",
	GEQ:    ">=",
	ASSIGN: ":=",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",
	RPAREN: ")",
	RBRACK: "]",
	RBRACE: "}",

	COMMA:     ",",
	COLON:     ":",
	SEMICOLON: ";",

	ARRAY: "array",
	ELSE:  "else",
	IF:    "if",
	OF:    "of",
	PROC:  "proc",
	REF:   "ref",
	TYPE:  "type",
	VAR:   "var",
	WHILE: "while",
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keywordBeg + 1; i < keywordEnd; i++ {
		keywords[tokens[i]] = i
	}
}

// String returns the string representation of the token. For operators,
// delimiters, and keywords the string is the actual token character sequence
// for all other tokens the string corresponds to the token constant name.
func (t Token) String() string {
	if 0 <= t && t < Token(len(tokens)) {
		return tokens[t]
	}
	return "token(" + strconv.Itoa(int(t)) + ")"
}

// IsLiteral returns true for tokens corresponding to identifiers and basic type
// literals. It returns false otherwise.
func (t Token) IsLiteral() bool { return literalBeg < t && t < literalEnd }

// IsOperator returns true for tokens corresponding to operators and delimiters.
// It returns false otherwise.
func (t Token) IsOperator() bool { return operatorBeg < t && t < operatorEnd }

// IsKeyword returns true for tokens corresponding to keywords. It returns false
// otherwise.
func (t Token) IsKeyword() bool { return keywordBeg < t && t < keywordEnd }

// A set of constants for precedence-based expression parsing. Non-operators
// have the lowest precedence, followed by operators starting with precedence 1
// up to unary operators. The highest precedence serves as "catch-all"
// precedence for selector, indexing, and other operator and delimiter tokens.
const (
	LowestPrec  = 0
	UnaryPrec   = 4
	HighestPrec = 5
)

// Precedence returns the operator precedence of the token. If t is not a binary
// operator, the result is LowestPrecedence.
func (t Token) Precedence() int {
	switch t {
	case EQL, LSS, LEQ, GTR, GEQ:
		return 1
	case ADD, SUB:
		return 2
	case MUL, QUO:
		return 3
	}
	return LowestPrec
}

// Lookup returns the keyword token associated with a given string. It returns
// IDENT if no matching keyword is found.
func Lookup(ident string) Token {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
