package types

import (
	"sort"

	"github.com/lukasmalkmus/spl/internal/app/spl/token"
)

// A Scope maintains a set of objects and links to its parent and children
// scopes. Objects may be inserted and looked up by name. The zero value for
// Scope is a ready-to-use empty scope.
type Scope struct {
	parent   *Scope
	children []*Scope
	elems    map[string]Object
	pos, end token.Position
	isFunc   bool
}

// NewScope returns a new, empty scope contained in the given parent scope, if
// any.
func NewScope(parent *Scope, pos, end token.Position) *Scope {
	s := &Scope{parent, nil, nil, pos, end, false}
	if parent != nil && parent != Universe {
		parent.children = append(parent.children, s)
	}
	return s
}

// Parent returns the scope's parent scope.
func (s *Scope) Parent() *Scope { return s.parent }

// Len returns the number of scope elements.
func (s *Scope) Len() int { return len(s.elems) }

// Names returns the scope's element names in sorted order.
func (s *Scope) Names() []string {
	names := make([]string, len(s.elems))
	i := 0
	for name := range s.elems {
		names[i] = name
		i++
	}
	sort.Strings(names)
	return names
}

// NumChildren returns the number of scopes nested in s.
func (s *Scope) NumChildren() int { return len(s.children) }

// Child returns the i'th child scope for 0 <= i < NumChildren().
func (s *Scope) Child(i int) *Scope { return s.children[i] }

// Lookup returns the object in scope s with the given name if such an object
// exists, otherwise the result is nil.
func (s *Scope) Lookup(name string) Object {
	return s.elems[name]
}

// LookupParent follows the parent chain of scopes starting with s until it
// finds a scope where Lookup(name) returns a non-nil object, and then returns
// that scope and object. If a valid position pos is provided, only objects that
// were declared at or before pos are considered. If no such scope and object
// exists, the result is (nil, nil).
func (s *Scope) LookupParent(name string, pos token.Position) (*Scope, Object) {
	// for ; s != nil; s = s.parent {
	// 	if obj := s.elems[name]; obj != nil && (!pos.IsValid() || obj.scopePos() <= pos) {
	// 		return s, obj
	// 	}
	// }
	return nil, nil
}

// Insert attempts to insert an object obj into the scope s. If s already
// contains an alternative object with the same name, Insert leaves s unchanged
// and returns the alternative object instead. Otherwise it inserts obj, sets
// the object's parent scope if not already set, and returns nil.
func (s *Scope) Insert(obj Object) Object {
	name := obj.Name()
	if alt := s.elems[name]; alt != nil {
		return alt
	}
	if s.elems == nil {
		s.elems = make(map[string]Object)
	}
	s.elems[name] = obj
	if obj.Parent() == nil {
		obj.setParent(s)
	}
	return nil
}

// Pos describes the scope's source code extent. The result is guaranteed to be
// valid only if the type-checked AST has complete position information.
func (s *Scope) Pos() token.Position { return s.pos }

// End describes the scope's source code extent. The result is guaranteed to be
// valid only if the type-checked AST has complete position information.
func (s *Scope) End() token.Position { return s.end }

// Contains reports whether pos is within the scope's extent. The result is
// guaranteed to be valid only if the type-checked AST has complete position
// information.
func (s *Scope) Contains(pos token.Position) bool {
	return s.pos.Char <= pos.Char && pos.Char < s.end.Char
}

// Innermost returns the innermost child scope containing pos. If pos is not
// within any scope, the result is nil. The result is also nil for the Universe
// scope. The result is guaranteed to be valid only if the type-checked AST has
// complete position information.
func (s *Scope) Innermost(pos token.Position) *Scope {
	// Package scopes do not have extents since they may be
	// discontiguous, so iterate over the package's files.
	if s.parent == Universe {
		for _, s := range s.children {
			if inner := s.Innermost(pos); inner != nil {
				return inner
			}
		}
	}

	if s.Contains(pos) {
		for _, s := range s.children {
			if s.Contains(pos) {
				return s.Innermost(pos)
			}
		}
		return s
	}
	return nil
}
