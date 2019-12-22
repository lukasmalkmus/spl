package types

// The Universe scope contains all predeclared objects of SPL. It is the
// outermost scope of any chain of nested scopes.
var Universe *Scope

// A builtinId is the ID of a builtin function.
type builtinId uint8

const (
	printi builtinId = iota
	printc
	readi
	readc
	exit
	time
	clearAll
	setPixel
	drawLine
	drawCircle
)

var predeclaredFuncs = [...]struct {
	name  string
	nargs int
	kind  exprKind
}{
	printi:     {"printi", 1, statement},
	printc:     {"printc", 1, statement},
	readi:      {"readi", 1, statement},
	readc:      {"readc", 1, statement},
	exit:       {"exit", 0, statement},
	time:       {"time", 1, statement},
	clearAll:   {"clearAll", 1, statement},
	setPixel:   {"setPixel", 3, statement},
	drawLine:   {"drawLine", 5, statement},
	drawCircle: {"drawCircle", 4, statement},
}
