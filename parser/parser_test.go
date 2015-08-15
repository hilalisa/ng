// Copyright 2015 The Numgrad Authors. All rights reserved.
// See the LICENSE file for rights to use this source code.

package parser

import (
	"math/big"
	"testing"
)

type parserTest struct {
	input string
	want  Expr
}

var parserTests = []parserTest{
	{"foo", &Ident{"foo"}},
	{"x + y", &BinaryExpr{Add, &Ident{"x"}, &Ident{"y"}}},
	{
		"x + y + 9",
		&BinaryExpr{
			Add,
			&BinaryExpr{Add, &Ident{"x"}, &Ident{"y"}},
			&BasicLiteral{big.NewInt(9)},
		},
	},
	{
		"x + (y + 7)",
		&BinaryExpr{
			Add,
			&Ident{"x"},
			&UnaryExpr{
				Op: LeftParen,
				Expr: &BinaryExpr{
					Add,
					&Ident{"y"},
					&BasicLiteral{big.NewInt(7)},
				},
			},
		},
	},
	{
		"x + y * z",
		&BinaryExpr{
			Add,
			&Ident{"x"},
			&BinaryExpr{Mul, &Ident{"y"}, &Ident{"z"}},
		},
	},
	{"y * /* comment */ z", &BinaryExpr{Mul, &Ident{"y"}, &Ident{"z"}}},
	{"y * z//comment", &BinaryExpr{Mul, &Ident{"y"}, &Ident{"z"}}},
	{
		"quit()",
		&CallExpr{Func: &Ident{Name: "quit"}},
	},
	{
		"foo(4)",
		&CallExpr{
			Func: &Ident{Name: "foo"},
			Args: []Expr{&BasicLiteral{Value: big.NewInt(4)}},
		},
	},
	{
		"min(1, 2)",
		&CallExpr{
			Func: &Ident{Name: "min"},
			Args: []Expr{
				&BasicLiteral{Value: big.NewInt(1)},
				&BasicLiteral{Value: big.NewInt(2)},
			},
		},
	},
}

func TestParseExpr(t *testing.T) {
	for _, test := range parserTests {
		got, err := ParseExpr([]byte(test.input))
		if err != nil {
			t.Errorf("ParseExpr(%q): error: %v", test.input, err)
			continue
		}
		if !EqualExpr(got, test.want) {
			t.Errorf("ParseExpr(%q):\n%v", test.input, Diff(test.want, got))
		}
	}
}