package parser

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/lukasmalkmus/spl/internal/app/spl/ast"
	"github.com/lukasmalkmus/spl/internal/app/spl/token"
)

func TestParser_ParseFullValidProgram(t *testing.T) {
	f, err := os.Open("../testdata/valid.spl")
	if err != nil {
		t.Fatal("failed to open testdata:", err)
	}
	p := NewFileParser(f)

	prog, _ := p.Parse()
	if p.errors.Len() > 0 {
		t.Errorf("expected no errors got %d: %s", p.errors.Len(), p.errors.Error())
	}
	if len(prog.Decls) == 0 {
		t.Errorf("didn't parse any top level declarations")
	}
	if prog.Name != "../testdata/valid.spl" {
		t.Errorf("invalid program name")
	}
}

func TestParser_ParseStatement(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    ast.Stmt
		wantErr bool
	}{
		{
			"assignment",
			"i := 0;",
			&ast.AssignStmt{
				Left:   &ast.Ident{NamePos: pos(1, 1), Name: "i"},
				Tok:    token.ASSIGN,
				TokPos: pos(1, 3),
				Right:  &ast.IntLit{ValuePos: pos(1, 6), Value: "0"},
			},
			false,
		},
		{
			"if",
			"if (i = 0) i + 1;",
			&ast.IfStmt{
				If: pos(1, 1),
				Cond: &ast.BinaryExpr{
					OpPos: pos(1, 7),
					Op:    token.EQL,
					X:     &ast.Ident{NamePos: pos(1, 5), Name: "i"},
					Y:     &ast.IntLit{ValuePos: pos(1, 9), Value: "0"},
				},
				Body: &ast.ExprStmt{X: &ast.BinaryExpr{
					OpPos: pos(1, 14),
					Op:    token.ADD,
					X:     &ast.Ident{NamePos: pos(1, 12), Name: "i"},
					Y:     &ast.IntLit{ValuePos: pos(1, 16), Value: "1"},
				}},
			},
			false,
		},
		{
			"while",
			"while (i > 0) i - 1;",
			&ast.WhileStmt{
				While: pos(1, 1),
				Cond: &ast.BinaryExpr{
					OpPos: pos(1, 10),
					Op:    token.GTR,
					X:     &ast.Ident{NamePos: pos(1, 8), Name: "i"},
					Y:     &ast.IntLit{ValuePos: pos(1, 12), Value: "0"},
				},
				Body: &ast.ExprStmt{X: &ast.BinaryExpr{
					OpPos: pos(1, 17),
					Op:    token.SUB,
					X:     &ast.Ident{NamePos: pos(1, 15), Name: "i"},
					Y:     &ast.IntLit{ValuePos: pos(1, 19), Value: "1"},
				}},
			},
			false,
		},
	}
	for _, tt := range tests {
		_ = t.Run(tt.name, func(t *testing.T) {
			got, perr := ParseStatement(tt.text)
			if err, _ := perr.(*ErrorList); (err != nil) != tt.wantErr {
				t.Errorf("parseDecl() error = %v, wantErr %v", err, tt.wantErr)
			}
			equals(t, got, tt.want)
		})
	}
}

func TestParser_parseDecl(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    ast.Decl
		wantErr bool
	}{
		{
			"variable",
			"var i: int;",
			&ast.VarDecl{
				Name: &ast.Ident{NamePos: pos(1, 5), Name: "i"},
				Type: &ast.Ident{NamePos: pos(1, 8), Name: "int"},
			},
			false,
		},
		{
			"type",
			"type myInt = int;",
			&ast.TypeDecl{
				Name:   &ast.Ident{NamePos: pos(1, 6), Name: "myInt"},
				Assign: pos(1, 12),
				Type:   &ast.Ident{NamePos: pos(1, 14), Name: "int"},
			},
			false,
		},
		{
			"type array",
			"type vector = array [5] of int;",
			&ast.TypeDecl{
				Name:   &ast.Ident{NamePos: pos(1, 6), Name: "vector"},
				Assign: pos(1, 13),
				Type: &ast.ArrayType{
					Array: pos(1, 15),
					Len:   &ast.IntLit{ValuePos: pos(1, 22), Value: "5"},
					Of:    pos(1, 25),
					Elt:   &ast.Ident{NamePos: pos(1, 28), Name: "int"},
				},
			},
			false,
		},
		{
			"type array double",
			"type matrix = array [3] of array [5] of int;",
			&ast.TypeDecl{
				Name:   &ast.Ident{NamePos: pos(1, 6), Name: "matrix"},
				Assign: pos(1, 13),
				Type: &ast.ArrayType{
					Array: pos(1, 15),
					Len:   &ast.IntLit{ValuePos: pos(1, 22), Value: "3"},
					Of:    pos(1, 25),
					Elt: &ast.ArrayType{
						Array: pos(1, 28),
						Len:   &ast.IntLit{ValuePos: pos(1, 35), Value: "5"},
						Of:    pos(1, 38),
						Elt:   &ast.Ident{NamePos: pos(1, 41), Name: "int"},
					},
				},
			},
			false,
		},
		{
			"procedure",
			"proc empty() {}",
			&ast.ProcDecl{
				Name: &ast.Ident{
					NamePos: pos(1, 6),
					Name:    "empty",
				},
				Proc:   pos(1, 1),
				Params: &ast.FieldList{Opening: pos(1, 11), Closing: pos(1, 12)},
				Body:   &ast.BlockStmt{Lbrace: pos(1, 14), Rbrace: pos(1, 15), List: []ast.Stmt{}},
			},
			false,
		},
		{
			"procedure one param",
			"proc one(a: int) {}",
			&ast.ProcDecl{
				Name: &ast.Ident{NamePos: pos(1, 6), Name: "one"},
				Proc: pos(1, 1),
				Params: &ast.FieldList{
					Opening: pos(1, 9),
					Closing: pos(1, 16),
					List: []*ast.Field{
						{
							Name: &ast.Ident{NamePos: pos(1, 10), Name: "a"},
							Type: &ast.Ident{NamePos: pos(1, 13), Name: "int"},
						},
					},
				},
				Body: &ast.BlockStmt{Lbrace: pos(1, 18), Rbrace: pos(1, 19), List: []ast.Stmt{}},
			},
			false,
		},
		{
			"procedure two params",
			"proc two(a: int, b: int) {}",
			&ast.ProcDecl{
				Name: &ast.Ident{NamePos: pos(1, 6), Name: "two"},
				Proc: pos(1, 1),
				Params: &ast.FieldList{
					Opening: pos(1, 9),
					Closing: pos(1, 24),
					List: []*ast.Field{
						{
							Name: &ast.Ident{NamePos: pos(1, 10), Name: "a"},
							Type: &ast.Ident{NamePos: pos(1, 13), Name: "int"},
						},
						{
							Name: &ast.Ident{NamePos: pos(1, 18), Name: "b"},
							Type: &ast.Ident{NamePos: pos(1, 21), Name: "int"},
						},
					},
				},
				Body: &ast.BlockStmt{Lbrace: pos(1, 26), Rbrace: pos(1, 27), List: []ast.Stmt{}},
			},
			false,
		},
		{
			"procedure two reference params",
			"proc swap(ref i: int, ref j: int) {}",
			&ast.ProcDecl{
				Name: &ast.Ident{NamePos: pos(1, 6), Name: "swap"},
				Proc: pos(1, 1),
				Params: &ast.FieldList{
					Opening: pos(1, 10),
					Closing: pos(1, 33),
					List: []*ast.Field{
						{
							Ref:  pos(1, 11),
							Name: &ast.Ident{NamePos: pos(1, 15), Name: "i"},
							Type: &ast.Ident{NamePos: pos(1, 18), Name: "int"},
						},
						{
							Ref:  pos(1, 23),
							Name: &ast.Ident{NamePos: pos(1, 27), Name: "j"},
							Type: &ast.Ident{NamePos: pos(1, 30), Name: "int"},
						},
					},
				},
				Body: &ast.BlockStmt{Lbrace: pos(1, 35), Rbrace: pos(1, 36), List: []ast.Stmt{}},
			},
			false,
		},
	}
	for _, tt := range tests {
		_ = t.Run(tt.name, func(t *testing.T) {
			p := New(strings.NewReader(tt.text))
			initParser(p)
			got, err := p.parseDecl(declStart), p.errors
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDecl() error = %v, wantErr %v", err, tt.wantErr)
			}
			equals(t, got, tt.want)
		})
	}
}

func pos(line, column int) token.Position {
	return token.Position{Filename: "", Line: line, Column: column}
}

func initParser(p *Parser) {
	p.openScope()
	p.pkgScope = p.topScope
	p.next()
}

// equals fails the test if got is not equal to want.
func equals(tb testing.TB, got, want interface{}) {
	tb.Helper()
	opts := cmp.Options{
		cmpopts.IgnoreTypes(&ast.Object{}),
		cmpopts.IgnoreFields(token.Position{}, "Char"),
	}
	if diff := cmp.Diff(got, want, opts...); diff != "" {
		tb.Errorf("\033[31m\n\n\tgot: %#+v\n\n\twant: %#+v\n\n\t%s\033[39m\n\n", got, want, diff)
	}
}
