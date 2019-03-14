package token_test

import (
	"go/token"
	"testing"
)

func TestPosition_String(t *testing.T) {
	tests := []struct {
		str string
		pos token.Position
	}{
		{"-", token.Position{}},
		{"-", token.Position{Column: 1}},
		{"1", token.Position{Line: 1}},
		{"a", token.Position{Filename: "a"}},
		{"a", token.Position{Filename: "a", Column: 1}},
		{"a:1", token.Position{Filename: "a", Line: 1}},
		{"1:1", token.Position{Filename: "", Line: 1, Column: 1}},
		{"a:1:2", token.Position{Filename: "a", Line: 1, Column: 2}},
	}
	for _, tt := range tests {
		_ = t.Run(tt.str, func(t *testing.T) {
			equals(t, tt.pos.String(), tt.str)
		})
	}
}
