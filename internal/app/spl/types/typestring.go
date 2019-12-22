package types

import (
	"bytes"
	"fmt"
)

// TypeString returns the string representation of typ.
func TypeString(typ Type) string {
	var buf bytes.Buffer
	WriteType(&buf, typ)
	return buf.String()
}

// WriteType writes the string representation of typ to buf.
func WriteType(buf *bytes.Buffer, typ Type) {
	writeType(buf, typ, make([]Type, 0, 8))
}

func writeType(buf *bytes.Buffer, typ Type, visited []Type) {
	for _, t := range visited {
		if t == typ {
			fmt.Fprintf(buf, "○%T", typ)
			return
		}
	}
	visited = append(visited, typ)

	switch t := typ.(type) {
	case nil:
		buf.WriteString("<nil>")
	case *Invalid:
		buf.WriteString("<invalid>")
	case *Int:
		buf.WriteString(t.name)
	case *Array:
		fmt.Fprintf(buf, "[%d]", t.len)
		writeType(buf, t.elem, visited)
	case *Tuple:
		writeTuple(buf, t, visited)
	case *Signature:
		buf.WriteString("proc")
		writeSignature(buf, t, visited)
	case *Named:
		s := "<Named w/o object>"
		if obj := t.obj; obj != nil {
			s = obj.name
		}
		buf.WriteString(s)
	default:
		buf.WriteString(t.String())
	}
}

func writeTuple(buf *bytes.Buffer, tup *Tuple, visited []Type) {
	buf.WriteByte('(')
	if tup != nil {
		for i, v := range tup.vars {
			if i > 0 {
				buf.WriteString(", ")
			}
			if v.name != "" {
				buf.WriteString(v.name)
				buf.WriteByte(' ')
			}
			writeType(buf, v.typ, visited)
		}
	}
	buf.WriteByte(')')
}

// WriteSignature writes the representation of the signature sig to buf, without
// a leading "proc" keyword.
func WriteSignature(buf *bytes.Buffer, sig *Signature) {
	writeSignature(buf, sig, make([]Type, 0, 8))
}

func writeSignature(buf *bytes.Buffer, sig *Signature, visited []Type) {
	writeTuple(buf, sig.params, visited)
}
