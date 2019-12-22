package constant

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/lukasmalkmus/spl/internal/app/spl/token"
)

// Kind specifies the kind of value represented by a Value.
type Kind int

// All available Kinds.
const (
	Unknown Kind = iota
	Bool
	Int
)

// A Value represents the value of a SPL constant.
type Value interface {
	// Kind returns the value kind.
	Kind() Kind

	// String returns a short, quoted (human-readable) form of the value.
	String() string

	// Prevent external implementations.
	implementsValue()
}

// -----------------------------------------------------------------------------
// Implementations

type (
	unknownVal struct{}
	boolVal    bool
	intVal     int32
)

func (unknownVal) Kind() Kind { return Unknown }
func (boolVal) Kind() Kind    { return Bool }
func (intVal) Kind() Kind     { return Int }

func (unknownVal) String() string { return "unknown" }
func (x boolVal) String() string  { return strconv.FormatBool(bool(x)) }
func (x intVal) String() string   { return strconv.FormatInt(int64(x), 10) }

func (unknownVal) implementsValue() {}
func (boolVal) implementsValue()    {}
func (intVal) implementsValue()     {}

func newInt() *intVal { return new(intVal) }

// -----------------------------------------------------------------------------
// Factories

// MakeUnknown returns the Unknown value.
func MakeUnknown() Value { return unknownVal{} }

// MakeBool returns the Bool value for b.
func MakeBool(b bool) Value { return boolVal(b) }

// MakeInt returns the Int value for x.
func MakeInt(x int32) Value { return intVal(x) }

// MakeFromLiteral returns the corresponding integer for a SPL literal string.
// The tok value must be one of token.INT. If the literal string syntax is
// invalid, the result is an Unknown.
func MakeFromLiteral(lit string, tok token.Token, zero uint) Value {
	if tok != token.INT {
		return unknownVal{}
	}
	x, err := strconv.ParseInt(lit, 0, 32)
	if err != nil {
		return unknownVal{}
	}
	return intVal(x)
}

// -----------------------------------------------------------------------------
// Accessors
//
// For unknown arguments the result is the zero value for the respective
// accessor type.

// BoolVal returns the SPL boolean value of x, which must be a Bool or an
// Unknown. If x is Unknown, the result is false.
func BoolVal(x Value) (bool, error) {
	switch x := x.(type) {
	case boolVal:
		return bool(x), nil
	case unknownVal:
		return false, nil
	}
	return false, fmt.Errorf("%v not a Bool", x)
}

// IntVal returns the SPL integer value of x, which must be an Int or an
// Unknown. If x is Unknown, the result is 0.
func IntVal(x Value) (int, error) {
	switch x := x.(type) {
	case intVal:
		return int(x), nil
	case unknownVal:
		return 0, nil
	}
	return 0, fmt.Errorf("%v not an Int", x)
}

// -----------------------------------------------------------------------------
// Support for assembling/disassembling numeric values

const (
	// Compute the size of a Word in bytes.
	_m       = ^big.Word(0)
	_log     = _m>>8&1 + _m>>16&1 + _m>>32&1
	wordSize = 1 << _log
)

// Bytes returns the bytes for the absolute value of x in little-
// endian binary representation; x must be an Int.
func Bytes(x Value) ([]byte, error) {
	t, ok := x.(*intVal)
	if !ok {
		return nil, fmt.Errorf("%v not an Int", x)
	}

	words := strconv.FormatInt(int64(*t), 2)
	bytes := make([]byte, len(words)*wordSize)

	i := 0
	for _, w := range words {
		for j := 0; j < wordSize; j++ {
			bytes[i] = byte(w)
			w >>= 8
			i++
		}
	}
	for i > 0 && bytes[i-1] == 0 {
		i--
	}
	return bytes[:i], nil
}

// -----------------------------------------------------------------------------
// Operations

// is32bit reports whether x can be represented using 32 bits.
func is32bit(x int64) bool {
	const s = 32
	return -1<<(s-1) <= x && x <= 1<<(s-1)-1
}

// UnaryOp returns the result of the unary expression op y. The operation must
// be defined for the operand. If y is Unknown, the result is Unknown.
func UnaryOp(op token.Token, y Value) (Value, error) {
	switch op {
	case token.ADD:
		switch y.(type) {
		case unknownVal, intVal:
			return y, nil
		}
	case token.SUB:
		switch y := y.(type) {
		case unknownVal:
			return y, nil
		case intVal:
			return MakeInt(-1 * int32(y)), nil
		}
	case token.NEQ:
		switch y := y.(type) {
		case unknownVal:
			return y, nil
		}
	}
	return nil, fmt.Errorf("invalid unary operation %s%v", op, y)
}

func ord(x Value) int {
	switch x.(type) {
	default:
		return -1
	case unknownVal:
		return 0
	case boolVal:
		return 1
	case intVal:
		return 2
	}
}

// match returns the matching representation (same type) with the
// smallest complexity for two values x and y. If one of them is
// numeric, both of them must be numeric. If one of them is Unknown
// or invalid (say, nil) both results are that value.
//
func match(x, y Value) (_, _ Value) {
	if ord(x) > ord(y) {
		y, x = match(y, x)
		return x, y
	}

	switch x := x.(type) {
	case boolVal:
		return x, y
	case intVal:
		switch y := y.(type) {
		case intVal:
			return x, y
		}
	}

	return x, x
}

// BinaryOp returns the result of the binary expression x op y.
// The operation must be defined for the operands. If one of the
// operands is Unknown, the result is Unknown.
// BinaryOp doesn't handle comparisons or shifts; use Compare
// or Shift instead.
//
// To force integer division of Int operands, use op == token.QUO_ASSIGN
// instead of token.QUO; the result is guaranteed to be Int in this case.
// Division by zero leads to a run-time panic.
//
func BinaryOp(xV Value, op token.Token, yV Value) Value {
	x, y := match(xV, yV)

	switch x := x.(type) {
	case unknownVal:
		return x

	case intVal:
		a := int32(x)
		b := int32(y.(intVal))
		var c int32
		switch op {
		case token.ADD:
			c = a + b
		case token.SUB:
			c = a - b
		case token.MUL:
			c = a * b
		case token.QUO:
			c = a / b
		default:
			goto Error
		}
		return MakeInt(c)
	}

Error:
	panic(fmt.Sprintf("invalid binary operation %v %s %v", xV, op, yV))
}

func add(x, y Value) Value { return BinaryOp(x, token.ADD, y) }
func sub(x, y Value) Value { return BinaryOp(x, token.SUB, y) }
func mul(x, y Value) Value { return BinaryOp(x, token.MUL, y) }
func quo(x, y Value) Value { return BinaryOp(x, token.QUO, y) }

func cmpZero(x int, op token.Token) bool {
	switch op {
	case token.EQL:
		return x == 0
	case token.NEQ:
		return x != 0
	case token.LSS:
		return x < 0
	case token.LEQ:
		return x <= 0
	case token.GTR:
		return x > 0
	case token.GEQ:
		return x >= 0
	}
	panic(fmt.Sprintf("invalid comparison %v %s 0", x, op))
}

// Compare returns the result of the comparison x op y.
// The comparison must be defined for the operands.
// If one of the operands is Unknown, the result is
// false.
//
func Compare(xV Value, op token.Token, yV Value) bool {
	x, y := match(xV, yV)

	switch x := x.(type) {
	case unknownVal:
		return false

	case boolVal:
		y := y.(boolVal)
		switch op {
		case token.EQL:
			return x == y
		case token.NEQ:
			return x != y
		}

	case intVal:
		y := y.(intVal)
		switch op {
		case token.EQL:
			return x == y
		case token.NEQ:
			return x != y
		case token.LSS:
			return x < y
		case token.LEQ:
			return x <= y
		case token.GTR:
			return x > y
		case token.GEQ:
			return x >= y
		}
	}

	panic(fmt.Sprintf("invalid comparison %v %s %v", xV, op, yV))
}
