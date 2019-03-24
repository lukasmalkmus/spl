package token

import (
	"fmt"
)

// NoPos is the Position zero value and is not valid.
var NoPos Position

// Position describes an arbitrary source position including the file, line,
// column and total character count. A Position is valid if the line number is
// > 0.
type Position struct {
	Filename string
	Line     int
	Column   int
	Char     int
}

// IsValid reports whether the position is valid.
func (pos *Position) IsValid() bool { return pos.Line > 0 }

// String returns a string representation of the position.
func (pos Position) String() string {
	s := pos.Filename
	if pos.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d", pos.Line)
		if pos.Column != 0 {
			s += fmt.Sprintf(":%d", pos.Column)
		}
	}
	if s == "" {
		s = "-"
	}
	return s
}
