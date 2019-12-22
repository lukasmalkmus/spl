package types

// Type represents a type of the Simple Programming Language.
type Type interface {
	// Underlying returns the underlying type of the type.
	Underlying() Type

	// String returns the string representation of the type.
	String() string
}

// Invalid is the primitive type for invalid types.
type Invalid struct{}

// Underlying returns the underlying type of the type.
func (i *Invalid) Underlying() Type { return i }

// String returns the string representation of the type.
func (i *Invalid) String() string { return TypeString(i) }

// Int is the primitive type for integers.
type Int struct {
	name string
}

// Underlying returns the underlying type of the type.
func (i *Int) Underlying() Type { return i }

// String returns the string representation of the type.
func (i *Int) String() string { return TypeString(i) }

// An Array represents an array type.
type Array struct {
	len  int64
	elem Type
}

// NewArray returns a new array type for the given element type and length. A
// negative length indicates an unknown length.
func NewArray(elem Type, len int64) *Array { return &Array{len, elem} }

// Underlying returns the underlying type of the type.
func (a *Array) Underlying() Type { return a }

// String returns the string representation of the type.
func (a *Array) String() string { return TypeString(a) }

// Len returns the length of array a. A negative result indicates an unknown
// length.
func (a *Array) Len() int64 { return a.len }

// Elem returns element type of array a.
func (a *Array) Elem() Type { return a.elem }

// A Tuple represents an ordered list of variables; a nil *Tuple is a valid
// (empty) tuple. Tuples are used as components of signatures, they are not
// first class types of SPL.
type Tuple struct {
	vars []*Var
}

// NewTuple returns a new tuple for the given variables.
func NewTuple(x ...*Var) *Tuple {
	if len(x) > 0 {
		return &Tuple{x}
	}
	return nil
}

// Underlying returns the underlying type of the type.
func (t *Tuple) Underlying() Type { return t }

// String returns the string representation of the type.
func (t *Tuple) String() string { return TypeString(t) }

// Len returns the number of variables of tuple t.
func (t *Tuple) Len() int {
	if t != nil {
		return len(t.vars)
	}
	return 0
}

// At returns the i'th variable of tuple t.
func (t *Tuple) At(i int) *Var { return t.vars[i] }

// A Signature represents a (non-builtin) procedure.
type Signature struct {
	scope  *Scope
	params *Tuple
}

// NewSignature returns a new procedure type for the given parameters.
func NewSignature(params *Tuple) *Signature {
	return &Signature{nil, params}
}

// Underlying returns the underlying type of the type.
func (s *Signature) Underlying() Type { return s }

// String returns the string representation of the type.
func (s *Signature) String() string { return TypeString(s) }

// Params returns the parameters of signature s, or nil.
func (s *Signature) Params() *Tuple { return s.params }

// A Named represents a named type.
type Named struct {
	obj        *TypeName
	underlying Type
}

// NewNamed returns a new named type for the given type name and underlying
// type. If the given type name obj doesn't have a type yet, its type is set to
// the returned named type. The underlying type must not be a *Named.
func NewNamed(obj *TypeName, underlying Type) *Named {
	if _, ok := underlying.(*Named); ok {
		panic("types.NewNamed: underlying type must not be *Named")
	}
	typ := &Named{obj: obj, underlying: underlying}
	if obj.typ == nil {
		obj.typ = typ
	}
	return typ
}

// Underlying returns the underlying type of the type.
func (t *Named) Underlying() Type { return t.underlying }

// String returns the string representation of the type.
func (t *Named) String() string { return TypeString(t) }

// Obj returns the type name for the named type t.
func (t *Named) Obj() *TypeName { return t.obj }

// SetUnderlying sets the underlying type and marks t as complete.
func (t *Named) SetUnderlying(underlying Type) {
	if underlying == nil {
		panic("types.Named.SetUnderlying: underlying type must not be nil")
	}
	if _, ok := underlying.(*Named); ok {
		panic("types.Named.SetUnderlying: underlying type must not be *Named")
	}
	t.underlying = underlying
}
