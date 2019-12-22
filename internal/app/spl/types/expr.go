package types

// exprKind describes the kind of an expression. The kind determines if an
// expression is valid in 'statement context'.
type exprKind int

const (
	expression exprKind = iota
	statement
)
