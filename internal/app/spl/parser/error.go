package parser

import (
	"fmt"
	"io"
	"sort"

	"github.com/lukasmalkmus/spl/internal/app/spl/token"
)

// Error represents an error in an ErrorList. The position Pos, if valid, points
// to the beginning of the offending token, and the error condition is described
// by Msg.
type Error struct {
	Pos token.Position
	Msg string
}

// Error implements the error interface.
func (e Error) Error() string {
	if e.Pos.Filename != "" || e.Pos.IsValid() {
		return e.Pos.String() + ": " + e.Msg
	}
	return e.Msg
}

// ErrorList is a list of *Errors. The zero value for an ErrorList is an empty
// ErrorList ready to use.
type ErrorList []*Error

// Add adds an Error with given position and error message to an ErrorList.
func (el *ErrorList) Add(pos token.Position, msg string) {
	*el = append(*el, &Error{pos, msg})
}

func (el ErrorList) Len() int { return len(el) }

func (el ErrorList) Swap(i, j int) { el[i], el[j] = el[j], el[i] }

func (el ErrorList) Less(i, j int) bool {
	e := &el[i].Pos
	f := &el[j].Pos
	if e.Filename != f.Filename {
		return e.Filename < f.Filename
	}
	if e.Line != f.Line {
		return e.Line < f.Line
	}
	if e.Column != f.Column {
		return e.Column < f.Column
	}
	return el[i].Msg < el[j].Msg
}

// Sort sorts an ErrorList. *Error entries are sorted by position, other errors
// are sorted by error message, and before any *Error entry.
func (el ErrorList) Sort() { sort.Sort(el) }

// An ErrorList implements the error interface.
func (el ErrorList) Error() string {
	switch len(el) {
	case 0:
		return "no errors"
	case 1:
		return el[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", el[0], len(el)-1)
}

// Err returns an error equivalent to this error list. If the list is empty, Err
// returns nil.
func (el ErrorList) Err() error {
	if len(el) == 0 {
		return nil
	}
	return el
}

// PrintError is a utility function that prints a list of errors to w, one error
// per line, if the err parameter is an ErrorList. Otherwise it prints the
// errors string.
func PrintError(w io.Writer, err error) {
	if list, ok := err.(ErrorList); ok {
		for _, e := range list {
			fmt.Fprintf(w, "%s\n", e)
		}
	} else if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	}
}
