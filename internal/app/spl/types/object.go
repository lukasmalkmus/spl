package types

import (
	"bytes"
	"fmt"
	"go/constant"

	"github.com/lukasmalkmus/spl/internal/app/spl/token"
)

// An Object describes a named language entity such as a constant, type,
// variable or procedure. All objects implement the Object interface.
type Object interface {
	// Parent returns the scope in which the object is declared.
	Parent() *Scope

	// Pos returns the declaration position of the object's identifier.
	Pos() token.Position

	// Name returns the object's name.
	Name() string

	// Type returns the object's type.
	Type() Type

	// String returns a human-readable string of the object.
	String() string

	// order reflects a package-level object's source order: If object a is
	// before object b in the source, then a.order() < b.order().
	order() uint32

	// scopePos returns the start position of the scope of this Object
	scopePos() token.Position

	// setOrder sets the order number of the object. It must be > 0.
	setOrder(uint32)

	// setParent sets the parent scope of the object.
	setParent(*Scope)

	// setScopePos sets the start position of the scope for this Object.
	setScopePos(pos token.Position)
}

// An object implements the common parts of an Object.
type object struct {
	parent    *Scope
	pos       token.Position
	name      string
	typ       Type
	order_    uint32
	scopePos_ token.Position
}

// Parent returns the scope in which the object is declared.
func (obj *object) Parent() *Scope { return obj.parent }

// Pos returns the declaration position of the object's identifier.
func (obj *object) Pos() token.Position { return obj.pos }

// Name returns the object's name.
func (obj *object) Name() string { return obj.name }

// Type returns the object's type.
func (obj *object) Type() Type { return obj.typ }

func (obj *object) order() uint32                  { return obj.order_ }
func (obj *object) scopePos() token.Position       { return obj.scopePos_ }
func (obj *object) setParent(parent *Scope)        { obj.parent = parent }
func (obj *object) setOrder(order uint32)          { obj.order_ = order }
func (obj *object) setScopePos(pos token.Position) { obj.scopePos_ = pos }

// A Const represents a declared constant.
type Const struct {
	object
	val constant.Value
}

// NewConst returns a new constant with value val.
func NewConst(pos token.Position, name string, typ Type, val constant.Value) *Const {
	return &Const{object{nil, pos, name, typ, 0, token.NoPos}, val}
}

// String returns a human-readable string of the constant.
func (obj *Const) String() string { return ObjectString(obj) }

// Val returns the constant's value.
func (obj *Const) Val() constant.Value { return obj.val }

// A TypeName represents a name for a defined type.
type TypeName struct {
	object
}

// NewTypeName returns a new type name denoting the given typ.
func NewTypeName(pos token.Position, name string, typ Type) *TypeName {
	return &TypeName{object{nil, pos, name, typ, 0, token.NoPos}}
}

// String returns a human-readable string of the type name.
func (obj *TypeName) String() string { return ObjectString(obj) }

// A Variable represents a declared variable (including procedure parameters).
type Var struct {
	object
	used bool
}

// NewVar returns a new variable.
func NewVar(pos token.Position, name string, typ Type) *Var {
	return &Var{object: object{nil, pos, name, typ, 0, token.NoPos}}
}

// NewParam returns a new variable representing a procedure parameter.
func NewParam(pos token.Position, name string, typ Type) *Var {
	return &Var{object: object{nil, pos, name, typ, 0, token.NoPos}, used: true}
}

// String returns a human-readable string of the variable.
func (obj *Var) String() string { return ObjectString(obj) }

// A Proc represents a declared procedure. Its Type() is always a *Signature.
type Proc struct {
	object
}

// NewProc returns a new procedure with the given signature, representing
// the procedure's type.
func NewProc(pos token.Position, name string, sig *Signature) *Proc {
	var typ Type
	if sig != nil {
		typ = sig
	}
	return &Proc{object{nil, pos, name, typ, 0, token.NoPos}}
}

// String returns a human-readable string of the procedure.
func (obj *Proc) String() string { return ObjectString(obj) }

// Scope returns the scope of the procedure's body block.
func (obj *Proc) Scope() *Scope { return obj.typ.(*Signature).scope }

// A Builtin represents a built-in procedure. Builtins don't have a valid type.
type Builtin struct {
	object
	id builtinId
}

func newBuiltin(id builtinId) *Builtin {
	return &Builtin{object{name: predeclaredFuncs[id].name, typ: &Invalid{}}, id}
}

// String returns a human-readable string of the built-in procedure.
func (obj *Builtin) String() string { return ObjectString(obj) }

// ObjectString returns the string form of obj.
func ObjectString(obj Object) string {
	var buf bytes.Buffer
	writeObject(&buf, obj)
	return buf.String()
}

func writeObject(buf *bytes.Buffer, obj Object) {
	var tname *TypeName
	typ := obj.Type()

	switch obj := obj.(type) {
	case *Const:
		buf.WriteString("const")
	case *TypeName:
		tname = obj
		buf.WriteString("type")
	case *Var:
		buf.WriteString("var")
	case *Proc:
		buf.WriteString("proc ")
		buf.WriteString(obj.name)
		if typ != nil {
			WriteSignature(buf, typ.(*Signature))
		}
		return
	case *Builtin:
		buf.WriteString("builtin")
		typ = nil
	default:
		panic(fmt.Sprintf("writeObject(%T)", obj))
	}

	buf.WriteByte(' ')
	buf.WriteString(obj.Name())

	if typ == nil {
		return
	}
	if tname != nil {
		// We have a type object: Don't print anything more for the integer type
		// since there's no more information.
		if _, ok := typ.(*Int); ok {
			return
		}
		typ = typ.Underlying()
	}

	buf.WriteByte(' ')
	WriteType(buf, typ)
}
