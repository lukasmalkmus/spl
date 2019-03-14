package ast

import (
	"github.com/lukasmalkmus/spl/internal/app/spl/token"
)

// Scope maintains the set of named language entities declared in the scope and
// a link to the immediately surrounding (outer) scope.
type Scope struct {
	Outer   *Scope
	Objects map[string]*Object
}

// NewScope creates a new scope nested in the outer scope.
func NewScope(outer *Scope) *Scope {
	return &Scope{outer, make(map[string]*Object)}
}

// Lookup returns the object with the given name if it is found in the scope,
// otherwise it returns nil. Outer scopes are ignored.
func (s *Scope) Lookup(name string) *Object {
	return s.Objects[name]
}

// Insert attempts to insert a named object into the scope. If the scope already
// contains an object with the same name, Insert leaves the scope unchanged and
// returns alt. Otherwise it inserts the object and returns nil.
func (s *Scope) Insert(obj *Object) (alt *Object) {
	if alt = s.Objects[obj.Name]; alt == nil {
		s.Objects[obj.Name] = obj
	}
	return
}

// -----------------------------------------------------------------------------
// Objects

// An Object describes a named language entity such as a type, variable or
// function.
//
// The Data field contains object-specific data:
//
//	Kind    Data type         Data value
//	Pkg     *Scope            package scope
type Object struct {
	Kind ObjKind
	Name string
	Decl interface{}
	Data interface{}
	Type interface{}
}

// NewObj creates a new object of a given kind and name.
func NewObj(kind ObjKind, name string) *Object { return &Object{Kind: kind, Name: name} }

// Pos computes the source position of the declaration of an object name. The
// result may be an invalid position if it cannot be computed (obj.Decl may be
// nil or not correct).
func (obj *Object) Pos() token.Position {
	name := obj.Name
	switch d := obj.Decl.(type) {
	case *Field:
		if d.Name.Name == name {
			return d.Name.Pos()
		}
	case *ProcDecl:
		if d.Name.Name == name {
			return d.Name.Pos()
		}
	}
	return token.NoPos
}

// ObjKind describes what an object represents.
type ObjKind int

// List of possible Object kinds.
const (
	Bad ObjKind = iota
	Typ
	Var
	Pro
)

var objKindStrings = [...]string{
	Bad: "bad",
	Typ: "type",
	Var: "var",
	Pro: "proc",
}

func (kind ObjKind) String() string { return objKindStrings[kind] }
