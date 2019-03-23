package ast

import "github.com/lukasmalkmus/spl/internal/app/spl/token"

// Node represents an element in the simple programming languages abstract
// syntax tree (AST).
type Node interface {
	// Pos returns the position of the first character belonging to the node.
	Pos() token.Position

	// End returns the position of the first character immediately after the
	// node.
	End() token.Position
}

// Expr is a simple programing language (SPL) expression.
type Expr interface {
	Node

	// exprNode is unexported to make sure implementations of Expr can only
	// originate in this package.
	exprNode()
}

func (*BadExpr) exprNode()    {}
func (*Ident) exprNode()      {}
func (*IntLit) exprNode()     {}
func (*ParenExpr) exprNode()  {}
func (*UnaryExpr) exprNode()  {}
func (*BinaryExpr) exprNode() {}
func (*IndexExpr) exprNode()  {}
func (*CallExpr) exprNode()   {}
func (*ArrayType) exprNode()  {}

// Stmt is a simple programing language (SPL) statement.
type Stmt interface {
	Node

	// stmtNode is unexported to make sure implementations of Stmt can only
	// originate in this package.
	stmtNode()
}

func (*BadStmt) stmtNode()    {}
func (*DeclStmt) stmtNode()   {}
func (*BlockStmt) stmtNode()  {}
func (*ExprStmt) stmtNode()   {}
func (*AssignStmt) stmtNode() {}
func (*WhileStmt) stmtNode()  {}
func (*IfStmt) stmtNode()     {}

// Decl is a simple programing language (SPL) declaration.
type Decl interface {
	Node

	// declNode is unexported to make sure implementations of Decl can only
	// originate in this package.
	declNode()
}

func (*BadDecl) declNode()  {}
func (*VarDecl) declNode()  {}
func (*TypeDecl) declNode() {}
func (*ProcDecl) declNode() {}

// -----------------------------------------------------------------------------
// Expressions and types

// Field represents a Field declaration list in a parameter declaration in a
// signature.
type Field struct {
	Ref  token.Position
	Name *Ident
	Type Expr
}

// Pos implements the Node interface.
func (f *Field) Pos() token.Position {
	if f.Ref != token.NoPos {
		return f.Ref
	}
	return f.Name.Pos()
}

// End implements the Node interface.
func (f *Field) End() token.Position { return f.Type.End() }

// FieldList represents a list of Fields, enclosed by parentheses or braces.
type FieldList struct {
	Opening token.Position
	List    []*Field
	Closing token.Position
}

// Pos implements the Node interface.
func (f *FieldList) Pos() token.Position { return f.Opening }

// End implements the Node interface.
func (f *FieldList) End() token.Position { return f.Closing }

// An expression is represented by a tree consisting of one or more of the
// following concrete expression nodes.
type (
	// BadExpr is a placeholder for expressions containing syntax errors for
	// which no correct expression nodes can be created.
	BadExpr struct {
		From token.Position
		To   token.Position
	}

	// Ident node represents an identifier.
	Ident struct {
		NamePos token.Position
		Name    string
		Obj     *Object
	}

	// IntLit represents a literal node of the integer type.
	IntLit struct {
		ValuePos token.Position
		Value    string
	}

	// ParenExpr represents a parenthesized expression node.
	ParenExpr struct {
		Lparen token.Position
		X      Expr
		Rparen token.Position
	}

	// UnaryExpr represents a unary expression node.
	UnaryExpr struct {
		OpPos token.Position
		Op    token.Token
		X     Expr
	}

	// BinaryExpr represents a binary expression node.
	BinaryExpr struct {
		OpPos token.Position
		Op    token.Token
		X     Expr
		Y     Expr
	}

	// IndexExpr represents an expression node followed by an index.
	IndexExpr struct {
		X      Expr
		Lbrack token.Position
		Index  Expr
		Rbrack token.Position
	}

	// CallExpr represents an expression node followed by an argument list.
	CallExpr struct {
		Pro    Expr
		Lparen token.Position
		Args   []Expr
		Rparen token.Position
	}
)

// Pos implements the Node interface.
func (x *BadExpr) Pos() token.Position { return x.From }

// End implements the Node interface.
func (x *BadExpr) End() token.Position { return x.To }

// Pos implements the Node interface.
func (x *Ident) Pos() token.Position { return x.NamePos }

// End implements the Node interface.
func (x *Ident) End() token.Position { pos := x.NamePos; pos.Line += len(x.Name); return pos }

// Pos implements the Node interface.
func (x *IntLit) Pos() token.Position { return x.ValuePos }

// End implements the Node interface.
func (x *IntLit) End() token.Position { pos := x.ValuePos; pos.Line += len(x.Value); return pos }

// Pos implements the Node interface.
func (x *ParenExpr) Pos() token.Position { return x.Lparen }

// End implements the Node interface.
func (x *ParenExpr) End() token.Position { return x.Rparen }

// Pos implements the Node interface.
func (x *UnaryExpr) Pos() token.Position { return x.OpPos }

// End implements the Node interface.
func (x *UnaryExpr) End() token.Position { return x.X.End() }

// Pos implements the Node interface.
func (x *BinaryExpr) Pos() token.Position { return x.X.Pos() }

// End implements the Node interface.
func (x *BinaryExpr) End() token.Position { return x.Y.End() }

// Pos implements the Node interface.
func (x *IndexExpr) Pos() token.Position { return x.Lbrack }

// End implements the Node interface.
func (x *IndexExpr) End() token.Position { return x.Rbrack }

// Pos implements the Node interface.
func (x *CallExpr) Pos() token.Position { return x.Pro.Pos() }

// End implements the Node interface.
func (x *CallExpr) End() token.Position { return x.Rparen }

func (x *Ident) String() string {
	if x != nil {
		return x.Name
	}
	return "<nil>"
}

// A type is represented by a tree consisting of one or more of the following
// type-specific expression nodes.
type (
	// ArrayType represents an array type node.
	ArrayType struct {
		Array token.Position
		Len   Expr
		Of    token.Position
		Elt   Expr
	}
)

// Pos implements the Node interface.
func (x *ArrayType) Pos() token.Position { return x.Array }

// End implements the Node interface.
func (x *ArrayType) End() token.Position { return x.Elt.Pos() }

// -----------------------------------------------------------------------------
// Statements

// A statement is represented by a tree consisting of one or more of the
// following concrete statement nodes.
type (
	// BadStmt node is a placeholder for statements containing syntax errors for
	// which no correct statement nodes can be created.
	BadStmt struct {
		From token.Position
		To   token.Position
	}

	// DeclStmt represents a declaration node in a statement list.
	DeclStmt struct {
		Decl Decl
	}

	// BlockStmt represents a braced statement list node.
	BlockStmt struct {
		Lbrace token.Position
		List   []Stmt
		Rbrace token.Position
	}

	// ExprStmt represents an (stand-alone) expression node in a statement list.
	ExprStmt struct {
		X Expr
	}

	// AssignStmt represents an assignment node.
	AssignStmt struct {
		Left   Expr
		TokPos token.Position
		Tok    token.Token
		Right  Expr
	}

	// WhileStmt represents a while node.
	WhileStmt struct {
		While token.Position
		Cond  Expr
		Body  Stmt
	}

	// IfStmt represents an if node.
	IfStmt struct {
		If   token.Position
		Cond Expr
		Body Stmt
		Else Stmt
	}
)

// Pos implements the Node interface.
func (s *BadStmt) Pos() token.Position { return s.From }

// End implements the Node interface.
func (s *BadStmt) End() token.Position { return s.To }

// Pos implements the Node interface.
func (s *DeclStmt) Pos() token.Position { return s.Decl.Pos() }

// End implements the Node interface.
func (s *DeclStmt) End() token.Position { return s.Decl.End() }

// Pos implements the Node interface.
func (s *BlockStmt) Pos() token.Position { return s.Lbrace }

// End implements the Node interface.
func (s *BlockStmt) End() token.Position { return s.Rbrace }

// Pos implements the Node interface.
func (s *ExprStmt) Pos() token.Position { return s.X.Pos() }

// End implements the Node interface.
func (s *ExprStmt) End() token.Position { return s.X.End() }

// Pos implements the Node interface.
func (s *AssignStmt) Pos() token.Position { return s.Left.Pos() }

// End implements the Node interface.
func (s *AssignStmt) End() token.Position { return s.Right.Pos() }

// Pos implements the Node interface.
func (s *WhileStmt) Pos() token.Position { return s.While }

// End implements the Node interface.
func (s *WhileStmt) End() token.Position { return s.Body.End() }

// Pos implements the Node interface.
func (s *IfStmt) Pos() token.Position { return s.If }

// End implements the Node interface.
func (s *IfStmt) End() token.Position {
	if s.Else != nil {
		return s.Else.End()
	}
	return s.Body.End()
}

// -----------------------------------------------------------------------------
// Declarations

// A declaration is represented by one of the following declaration nodes.
type (
	// BadDecl is a placeholder for declarations containing syntax errors for
	// which no correct declaration nodes can be created.
	BadDecl struct {
		From token.Position
		To   token.Position
	}

	// VarDecl represents a variable declaration node.
	VarDecl struct {
		Name *Ident
		Type Expr
	}

	// TypeDecl represents a type declaration node.
	TypeDecl struct {
		Name   *Ident
		Assign token.Position
		Type   Expr
	}

	// ProcDecl represents a procedure declaration node.
	ProcDecl struct {
		Name   *Ident
		Proc   token.Position
		Params *FieldList
		Body   *BlockStmt
	}
)

// Pos implements the Node interface.
func (d *BadDecl) Pos() token.Position { return d.From }

// End implements the Node interface.
func (d *BadDecl) End() token.Position { return d.To }

// Pos implements the Node interface.
func (d *VarDecl) Pos() token.Position { return d.Name.NamePos }

// End implements the Node interface.
func (d *VarDecl) End() token.Position { return d.Type.End() }

// Pos implements the Node interface.
func (d *TypeDecl) Pos() token.Position { return d.Name.NamePos }

// End implements the Node interface.
func (d *TypeDecl) End() token.Position { return d.Type.End() }

// Pos implements the Node interface.
func (d *ProcDecl) Pos() token.Position { return d.Proc }

// End implements the Node interface.
func (d *ProcDecl) End() token.Position { return d.Body.End() }

// -----------------------------------------------------------------------------
// Helpers

// Program represents a simple programing language (SPL) program AST.
type Program struct {
	Name       string
	Decls      []Decl
	Unresolved []*Ident
}

// Pos implements the Node interface.
func (p *Program) Pos() token.Position { return token.Position{Filename: p.Name, Line: 1, Column: 1} }

// End implements the Node interface.
func (p *Program) End() token.Position {
	if n := len(p.Decls); n > 0 {
		return p.Decls[n-1].End()
	}
	return token.Position{Filename: p.Name, Line: 1, Column: 1}
}
