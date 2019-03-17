package parser

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lukasmalkmus/spl/internal/app/spl/ast"
	"github.com/lukasmalkmus/spl/internal/app/spl/scanner"
	"github.com/lukasmalkmus/spl/internal/app/spl/token"
)

// Parser implements a parser for the simple programing language (SPL). The
// parser initializes a scanner itself which is used for the lexical analysis of
// the source code.
type Parser struct {
	scanner *scanner.Scanner
	errors  ErrorList

	// Current token
	tok token.Token
	lit string
	pos token.Position

	// Buffered token
	buf struct {
		tok token.Token
		lit string
		pos token.Position
		n   int
	}

	// Error recovery
	syncCnt int
	syncPos token.Position

	// Non-syntactic parser control
	exprLev int
	inRHS   bool

	// Ordinary identifier scopes
	pkgScope   *ast.Scope
	topScope   *ast.Scope
	unresolved []*ast.Ident
}

// New returns a new Parser which is initialized with the provided reader as its
// source.
func New(r io.Reader) *Parser {
	// Init Parser with EOF token. This ensures functions must read the first
	// token themselves.
	return &Parser{
		scanner: scanner.New(r),

		tok: token.EOF,
	}
}

// NewFileParser returns a new instance of Parser, but will exclusively take an
// *os.File as argument instead of the more general io.Reader interface.
// Therefore it will enhance token positions with the filename.
func NewFileParser(f *os.File) *Parser {
	// Init Parser with EOF token. This ensures the first token must be read
	// explicitly.
	return &Parser{
		scanner: scanner.NewFileScanner(f),

		tok: token.EOF,
		pos: token.Position{Filename: f.Name()},
	}
}

// ParseStatement parses a single SPL statement.
func ParseStatement(src string) (ast.Stmt, error) {
	p := New(strings.NewReader(src))
	return p.parseStmt(), p.errors
}

// Feed will provide the parser with a new scanner source, which effectively
// adds a new source of tokens. This preserves the previous parsing context
// when parsing new data.
func (p *Parser) Feed(r io.Reader) { p.scanner = scanner.New(r) }

// Parse parses the source the Parser is initialized with.
func (p *Parser) Parse() {
	// If scanning the first token fails, this is probably not a spl source
	// file.
	if p.next(); p.errors.Len() != 0 {
		return
	}

	p.openScope()
	p.pkgScope = p.topScope
	var decls []ast.Decl
	for p.tok != token.EOF {
		decls = append(decls, p.parseDecl(declStart))
	}
	p.closeScope()

	// Resolve global identifiers.
	i := 0
	for _, ident := range p.unresolved {
		ident.Obj = p.pkgScope.Lookup(ident.Name)
		if ident.Obj == nil {
			p.unresolved[i] = ident
			i++
		}
	}
	_ = decls
}

// -----------------------------------------------------------------------------
// Declarations

// parseDecl parses a declaration AST object.
func (p *Parser) parseDecl(sync map[token.Token]bool) ast.Decl {
	switch p.tok {
	case token.VAR:
		return p.parseVarDecl()
	case token.TYPE:
		return p.parseTypeDecl()
	case token.PROC:
		return p.parseProcDecl()
	}
	pos := p.pos
	p.advance(sync)
	p.errorExpected(pos, "declaration")
	return &ast.BadDecl{From: pos, To: p.pos}
}

// parseVarDecl parses a variable declaration AST object.
func (p *Parser) parseVarDecl() *ast.VarDecl {
	_ = p.expect(token.VAR)
	ident := p.parseIdent()
	_ = p.expect(token.COLON)
	typ := p.tryType()
	p.expectSemi()
	if typ == nil {
		p.error(ident.NamePos, "missing variable type")
	}

	decl := &ast.VarDecl{Name: ident, Type: typ}
	p.declare(decl, nil, p.topScope, ast.Var, ident)
	return decl
}

// parseTypeDecl parses a type declaration AST object.
func (p *Parser) parseTypeDecl() *ast.TypeDecl {
	_ = p.expect(token.TYPE)
	ident := p.parseIdent()
	decl := &ast.TypeDecl{Name: ident}
	p.declare(decl, nil, p.topScope, ast.Typ, ident)
	decl.Assign = p.expect(token.EQL)
	decl.Type = p.parseType()
	p.expectSemi()
	return decl
}

func (p *Parser) parseProcDecl() *ast.ProcDecl {
	pos := p.expect(token.PROC)
	scope := ast.NewScope(p.topScope)
	ident := p.parseIdent()
	params := p.parseParameters(scope)
	body := p.parseBody(scope)

	decl := &ast.ProcDecl{
		Name:   ident,
		Proc:   pos,
		Params: params,
		Body:   body,
	}
	p.declare(decl, nil, p.pkgScope, ast.Pro, ident)
	return decl
}

// -----------------------------------------------------------------------------
// Identifiers

// parseIdent parses an identifier AST object.
func (p *Parser) parseIdent() *ast.Ident {
	pos := p.pos
	var name string
	if p.tok == token.IDENT {
		name = p.lit
		p.next()
	} else {
		_ = p.expect(token.IDENT)
	}
	return &ast.Ident{NamePos: pos, Name: name}
}

// ----------------------------------------------------------------------------
// Common productions

func (p *Parser) parseLHS() ast.Expr {
	old := p.inRHS
	p.inRHS = false
	x := p.checkExpr(p.parseExpr(true))
	p.resolve(x)
	p.inRHS = old
	return x
}

func (p *Parser) parseRHS() ast.Expr {
	old := p.inRHS
	p.inRHS = true
	x := p.checkExpr(p.parseExpr(false))
	p.inRHS = old
	return x
}

// -----------------------------------------------------------------------------
// Types

func (p *Parser) parseType() ast.Expr {
	typ := p.tryType()
	if typ == nil {
		pos := p.pos
		p.errorExpected(pos, "type")
		p.advance(exprEnd)
		return &ast.BadExpr{From: pos, To: p.pos}
	}
	return typ
}

func (p *Parser) tryType() ast.Expr {
	typ := p.tryIdentOrType()
	if typ != nil {
		p.resolve(typ)
	}
	return typ
}

func (p *Parser) parseParameters(scope *ast.Scope) *ast.FieldList {
	var params []*ast.Field
	lparen := p.expect(token.LPAREN)
	if p.tok != token.RPAREN {
		params = p.parseParameterList(scope)
	}
	rparen := p.expect(token.RPAREN)
	return &ast.FieldList{Opening: lparen, List: params, Closing: rparen}
}

func (p *Parser) parseParameterList(scope *ast.Scope) []*ast.Field {
	params := make([]*ast.Field, 0)
	for p.tok != token.RPAREN && p.tok != token.EOF {
		ref := p.optional(token.REF)
		ident := p.parseIdent()
		_ = p.expect(token.COLON)
		typ := p.parseVarType()
		field := &ast.Field{Ref: ref, Name: ident, Type: typ}
		params = append(params, field)
		p.declare(field, nil, scope, ast.Var, ident)
		p.resolve(typ)
		if !p.atComma("parameter list", token.RPAREN) {
			break
		}
		p.next()
	}
	return params
}

// If the result is an identifier, it is not resolved.
func (p *Parser) parseVarType() ast.Expr {
	typ := p.tryIdentOrType()
	if typ == nil {
		pos := p.pos
		p.errorExpected(pos, "type")
		p.next()
		typ = &ast.BadExpr{From: pos, To: p.pos}
	}
	return typ
}

// If the result is an identifier, it is not resolved.
func (p *Parser) tryIdentOrType() ast.Expr {
	switch p.tok {
	case token.IDENT:
		return p.parseIdent()
	case token.ARRAY:
		return p.parseArrayType()
	case token.LPAREN:
		lparen := p.pos
		p.next()
		typ := p.parseType()
		rparen := p.expect(token.RPAREN)
		return &ast.ParenExpr{Lparen: lparen, X: typ, Rparen: rparen}
	}
	return nil
}

func (p *Parser) parseArrayType() ast.Expr {
	array := p.expect(token.ARRAY)
	_ = p.expect(token.LBRACK)
	p.exprLev++
	len := p.parseRHS()
	p.exprLev--
	_ = p.expect(token.RBRACK)
	of := p.expect(token.OF)
	elt := p.parseType()
	return &ast.ArrayType{Array: array, Len: len, Of: of, Elt: elt}
}

// -----------------------------------------------------------------------------
// Blocks

func (p *Parser) parseStmtList() []ast.Stmt {
	list := make([]ast.Stmt, 0)
	for p.tok != token.RBRACE && p.tok != token.EOF {
		list = append(list, p.parseStmt())
	}
	return list
}

func (p *Parser) parseBody(scope *ast.Scope) *ast.BlockStmt {
	lbrace := p.expect(token.LBRACE)
	p.topScope = scope
	list := p.parseStmtList()
	p.closeScope()
	rbrace := p.expect(token.RBRACE)
	return &ast.BlockStmt{Lbrace: lbrace, List: list, Rbrace: rbrace}
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	lbrace := p.expect(token.LBRACE)
	p.openScope()
	list := p.parseStmtList()
	p.closeScope()
	rbrace := p.expect(token.RBRACE)
	return &ast.BlockStmt{Lbrace: lbrace, List: list, Rbrace: rbrace}
}

// -----------------------------------------------------------------------------
// Expressions

// checkExpr checks that x is an expression (and not a type).
func (p *Parser) checkExpr(x ast.Expr) ast.Expr {
	switch unparen(x).(type) {
	case *ast.BadExpr:
	case *ast.Ident:
	case *ast.IntLit:
	case *ast.UnaryExpr:
	case *ast.BinaryExpr:
	case *ast.IndexExpr:
	case *ast.CallExpr:
	default:
		p.errorExpected(x.Pos(), "expression")
		return &ast.BadExpr{From: x.Pos(), To: x.End()}
	}
	return x
}

// If lhs is set and the result is an identifier, it is not resolved. The result
// may be a type and callers must check the result (using checkExpr or
// checkExprOrType).
func (p *Parser) parseExpr(lhs bool) ast.Expr {
	return p.parseBinaryExpr(lhs, token.LowestPrec+1)
}

// If lhs is set and the result is an identifier, it is not resolved.
func (p *Parser) parseBinaryExpr(lhs bool, prec1 int) ast.Expr {
	x := p.parseUnaryExpr(lhs)
	for {
		op, oprec := p.tokPrec()
		if oprec < prec1 {
			return x
		}
		pos := p.expect(op)
		if lhs {
			p.resolve(x)
			lhs = false
		}
		y := p.parseBinaryExpr(false, oprec+1)
		x = &ast.BinaryExpr{X: p.checkExpr(x), OpPos: pos, Op: op, Y: p.checkExpr(y)}
	}
}

// If lhs is set and the result is an identifier, it is not resolved.
func (p *Parser) parseUnaryExpr(lhs bool) ast.Expr {
	if p.tok == token.ADD || p.tok == token.SUB || p.tok == token.NOT || p.tok == token.MUL {
		pos, op := p.pos, p.tok
		p.next()
		x := p.parseUnaryExpr(false)
		return &ast.UnaryExpr{OpPos: pos, Op: op, X: p.checkExpr(x)}
	}
	return p.parsePrimaryExpr(lhs)
}

// If lhs is set and the result is an identifier, it is not resolved.
func (p *Parser) parsePrimaryExpr(lhs bool) ast.Expr {
	x := p.parseOperand(lhs)
L:
	for {
		switch p.tok {
		case token.LBRACK:
			if lhs {
				p.resolve(x)
			}
			x = p.parseIndex(p.checkExpr(x))
		case token.LPAREN:
			if lhs {
				p.resolve(x)
			}
			x = p.parseCall(p.checkExpr(x))
		default:
			break L
		}
		lhs = false
	}
	return x
}

// parseOperand may return an expression or a raw type (incl. array types).
// Callers must verify the result. If lhs is set and the result is an
// identifier, it is not resolved.
func (p *Parser) parseOperand(lhs bool) ast.Expr {
	switch p.tok {
	case token.IDENT:
		x := p.parseIdent()
		if !lhs {
			p.resolve(x)
		}
		return x
	case token.INT:
		x := &ast.IntLit{ValuePos: p.pos, Value: p.lit}
		p.next()
		return x
	case token.LPAREN:
		lparen := p.pos
		p.next()
		p.exprLev++
		x := p.parseRHS()
		p.exprLev--
		rparen := p.expect(token.RPAREN)
		return &ast.ParenExpr{Lparen: lparen, X: x, Rparen: rparen}
	}

	if typ := p.tryIdentOrType(); typ != nil {
		return typ
	}

	pos := p.pos
	p.errorExpected(pos, "operand")
	p.advance(stmtStart)
	return &ast.BadExpr{From: pos, To: p.pos}
}

func (p *Parser) parseIndex(x ast.Expr) ast.Expr {
	lbrack := p.expect(token.LBRACK)
	p.exprLev++
	index := p.parseRHS()
	p.exprLev--
	rbrack := p.expect(token.RBRACK)
	return &ast.IndexExpr{X: x, Lbrack: lbrack, Index: index, Rbrack: rbrack}
}

func (p *Parser) parseCall(pro ast.Expr) *ast.CallExpr {
	lparen := p.expect(token.LPAREN)
	p.exprLev++
	var list []ast.Expr
	for p.tok != token.RPAREN && p.tok != token.EOF {
		list = append(list, p.parseRHS())
		if !p.atComma("argument list", token.RPAREN) {
			break
		}
		p.next()
	}
	p.exprLev--
	rparen := p.expectClosing(token.RPAREN, "argument list")
	return &ast.CallExpr{Pro: pro, Lparen: lparen, Args: list, Rparen: rparen}
}

func (p *Parser) tokPrec() (token.Token, int) {
	tok := p.tok
	if p.inRHS && tok == token.EQL {
		tok = token.EQL
	}
	return tok, tok.Precedence()
}

// If x is of the form (T), unparen returns unparen(T), otherwise it returns x.
func unparen(x ast.Expr) ast.Expr {
	if p, isParen := x.(*ast.ParenExpr); isParen {
		x = unparen(p.X)
	}
	return x
}

// -----------------------------------------------------------------------------
// Statements

// parseStmt parses lexical tokens into a Statement AST object.
func (p *Parser) parseStmt() (stmt ast.Stmt) {
	switch p.tok {
	case token.VAR, token.TYPE:
		stmt = &ast.DeclStmt{Decl: p.parseDecl(stmtStart)}
	case token.IDENT, token.INT, token.LPAREN,
		token.LBRACK, token.ADD, token.SUB, token.MUL, token.NOT:
		stmt = p.parseSimpleStmt()
		p.expectSemi()
	case token.LBRACE:
		stmt = p.parseBlockStmt()
		p.expectSemi()
	case token.WHILE:
		stmt = p.parseWhileStmt()
	case token.IF:
		stmt = p.parseIfStmt()
	default:
		pos := p.pos
		p.errorExpected(pos, "statement")
		p.advance(stmtStart)
		stmt = &ast.BadStmt{From: pos, To: p.pos}
	}
	return stmt
}

func (p *Parser) parseSimpleStmt() ast.Stmt {
	x := p.parseLHS()
	if p.tok == token.ASSIGN {
		pos, tok := p.pos, p.tok
		p.next()
		y := p.parseRHS()
		return &ast.AssignStmt{Left: x, TokPos: pos, Tok: tok, Right: y}
	}
	return &ast.ExprStmt{X: x}
}

func (p *Parser) parseWhileStmt() ast.Stmt {
	pos := p.expect(token.WHILE)
	_ = p.expect(token.LPAREN)
	prevLev := p.exprLev
	p.exprLev = -1
	x := p.checkExpr(p.parseExpr(false))
	p.exprLev = prevLev
	_ = p.expect(token.RPAREN)
	body := p.parseBlockStmt()

	return &ast.WhileStmt{
		While: pos,
		Cond:  x,
		Body:  body,
	}
}

func (p *Parser) parseIfStmt() ast.Stmt {
	pos := p.expect(token.IF)
	_ = p.expect(token.LPAREN)
	prevLev := p.exprLev
	p.exprLev = -1
	x := p.checkExpr(p.parseExpr(false))
	p.exprLev = prevLev
	_ = p.expect(token.RPAREN)
	body := p.parseBlockStmt()

	var alt ast.Stmt
	if p.tok == token.ELSE {
		p.next()
		alt = p.parseBlockStmt()
	}

	return &ast.IfStmt{
		If:   pos,
		Cond: x,
		Body: body,
		Else: alt,
	}
}

// -----------------------------------------------------------------------------
// Scoping support

func (p *Parser) openScope() {
	p.topScope = ast.NewScope(p.topScope)
}

func (p *Parser) closeScope() {
	p.topScope = p.topScope.Outer
}

func (p *Parser) declare(decl, data interface{}, scope *ast.Scope, kind ast.ObjKind, idents ...*ast.Ident) {
	for _, ident := range idents {
		obj := ast.NewObj(kind, ident.Name)
		obj.Decl = decl
		obj.Data = data
		ident.Obj = obj
		if ident.Name != "_" {
			if alt := scope.Insert(obj); alt != nil {
				prevDecl := ""
				if pos := alt.Pos(); pos.IsValid() {
					prevDecl = fmt.Sprintf("\n\tprevious declaration at %s", pos)
				}
				p.error(ident.Pos(), fmt.Sprintf("%s redeclared in this block%s", ident.Name, prevDecl))
			}
		}
	}
}

// The unresolved object is a sentinel to mark identifiers that have been added
// to the list of unresolved identifiers. The sentinel is only used for
// verifying internal consistency.
var unresolved = new(ast.Object)

func (p *Parser) tryResolve(x ast.Expr, collectUnresolved bool) {
	// Nothing to do if x is not an identifier.
	ident, ok := x.(*ast.Ident)
	if !ok {
		return
	}

	for s := p.topScope; s != nil; s = s.Outer {
		if obj := s.Lookup(ident.Name); obj != nil {
			ident.Obj = obj
			return
		}
	}

	if collectUnresolved {
		ident.Obj = unresolved
		p.unresolved = append(p.unresolved, ident)
	}
}

func (p *Parser) resolve(x ast.Expr) {
	p.tryResolve(x, true)
}

// -----------------------------------------------------------------------------
// Parsing support

func (p *Parser) expect(tok token.Token) token.Position {
	pos := p.pos
	if p.tok != tok {
		p.errorExpected(pos, "'"+tok.String()+"'")
	}
	p.next()
	return pos
}

func (p *Parser) optional(tok token.Token) token.Position {
	pos := p.pos
	if p.tok != tok {
		return token.NoPos
	}
	p.next()
	return pos
}

// expectClosing is like expect but provides a better error message for the
// common case of a missing comma before a newline.
func (p *Parser) expectClosing(tok token.Token, context string) token.Position {
	if p.tok != tok && p.tok == token.SEMICOLON && p.lit == "\n" {
		p.error(p.pos, "missing ',' before newline in "+context)
		p.next()
	}
	return p.expect(tok)
}

func (p *Parser) expectSemi() {
	if p.tok != token.RPAREN && p.tok != token.RBRACE {
		switch p.tok {
		case token.COMMA:
			p.errorExpected(p.pos, "';'")
			fallthrough
		case token.SEMICOLON:
			p.next()
		default:
			p.errorExpected(p.pos, "';'")
			p.advance(stmtStart)
		}
	}
}

func (p *Parser) atComma(context string, follow token.Token) bool {
	if p.tok == token.COMMA {
		return true
	}
	if p.tok != follow {
		msg := "missing ','"
		if p.tok == token.SEMICOLON && p.lit == "\n" {
			msg += " before newline"
		}
		p.error(p.pos, msg+" in "+context)
		return true
	}
	return false
}

// next scans the next non-comment token.
func (p *Parser) next() {
	// TODO: Collect comments.
	p.scan()
	for p.tok == token.COMMENT {
		p.scan()
	}
}

// scan returns the next token from the underlying scanner. If a token has been
// unscanned read that one instead.
func (p *Parser) scan() {
	// If we have a token on the buffer return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		p.tok, p.lit, p.pos = p.buf.tok, p.buf.lit, p.buf.pos
		return
	}

	// Otherwise read the next token from the scanner and save it to the buffer
	// in case we unscan later.
	p.tok, p.lit, p.pos = p.scanner.Scan()
	p.buf.tok, p.buf.lit, p.buf.pos = p.tok, p.lit, p.pos
}

var (
	stmtStart = map[token.Token]bool{
		token.ARRAY: true,
		token.ELSE:  true,
		token.IF:    true,
		token.OF:    true,
		token.PROC:  true,
		token.REF:   true,
		token.TYPE:  true,
		token.VAR:   true,
		token.WHILE: true,
	}

	declStart = map[token.Token]bool{
		token.TYPE: true,
		token.VAR:  true,
	}

	exprEnd = map[token.Token]bool{
		token.RPAREN: true,
	}
)

// advance consumes tokens until the current token is in the provided set, or
// token.EOF.
func (p *Parser) advance(to map[token.Token]bool) {
	for ; p.tok != token.EOF; p.next() {
		if to[p.tok] {
			if p.pos == p.syncPos && p.syncCnt < 10 {
				p.syncCnt++
			} else if p.pos.Line > p.syncPos.Line && p.pos.Column > p.syncPos.Column {
				p.syncPos = p.pos
				p.syncCnt = 0
			}
		}
	}
}

// unscan pushes the previously read token back onto the buffer.
// func (p *Parser) unscan() { p.buf.n = 1 }

// -----------------------------------------------------------------------------
// Errors

func (p *Parser) error(pos token.Position, msg string) { p.errors.Add(pos, msg) }

func (p *Parser) errorExpected(pos token.Position, msg string) {
	msg = "expected " + msg
	if pos == p.pos {
		switch {
		case p.tok == token.SEMICOLON && p.lit == "\n":
			msg += ", found newline"
		case p.tok.IsLiteral():
			msg += ", found " + p.lit
		default:
			msg += ", found '" + p.tok.String() + "'"
		}
	}
	p.error(pos, msg)
}
