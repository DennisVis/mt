// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ast_test

import (
	"testing"

	"github.com/DennisVis/mt/internal/pattern/ast"
)

func TestAST(t *testing.T) {
	t.Helper()

	for _, test := range []struct {
		name         string
		node         ast.Node
		expectedKind ast.NodeKind
		expectedStr  string
	}{
		{
			name: "Literal",
			node: &ast.Literal{
				Value: "/",
			},
			expectedKind: ast.NodeKindLiteral,
			expectedStr: `Literal:
	Value: "/"
`,
		},
		{
			name: "OptionalLiteral",
			node: &ast.Optional{
				Node: &ast.Literal{
					Value: "/",
				},
			},
			expectedKind: ast.NodeKindOptional,
			expectedStr: `Optional:
	Node:
		Literal:
			Value: "/"
`,
		},
		{
			name: "CharGroup",
			node: &ast.CharGroup{
				CharCount:       1,
				CharCountStrict: true,
				CharSetKey:      "n",
			},
			expectedKind: ast.NodeKindCharGroup,
			expectedStr: `CharGroup:
	CharCount: 1
	CharCountStrict: true
	CharSetKey: n
`,
		},
		{
			name: "LineCountExpression",
			node: &ast.LineCountExpression{
				LineCount: 2,
				Node: &ast.Literal{
					Value: "/",
				},
			},
			expectedKind: ast.NodeKindLineCountExpression,
			expectedStr: `LineCountExpression:
	LineCount: 2
	Node:
		Literal:
			Value: "/"
`,
		},
		{
			name: "OrExpression",
			node: &ast.OrExpression{
				Left: &ast.Literal{
					Value: "/",
				},
				Right: &ast.Literal{
					Value: `\`,
				},
			},
			expectedKind: ast.NodeKindOrExpression,
			expectedStr: `OrExpression:
	Left:
		Literal:
			Value: "/"
	Right:
		Literal:
			Value: "\\"
`,
		},
		{
			name: "OrExpressionEmpty",
			node: &ast.OrExpression{
				Left:  nil,
				Right: nil,
			},
			expectedKind: ast.NodeKindOrExpression,
			expectedStr: `OrExpression:
	Left:
		EMPTY
	Right:
		EMPTY
`,
		},
		{
			name: "Pattern",
			node: &ast.Pattern{
				Nodes: []ast.Node{
					&ast.Literal{
						Value: "/",
					},
					&ast.CharGroup{
						CharCount:       1,
						CharCountStrict: true,
						CharSetKey:      "n",
					},
				},
			},
			expectedKind: ast.NodeKindPattern,
			expectedStr: `Pattern:
	Literal:
		Value: "/"
	CharGroup:
		CharCount: 1
		CharCountStrict: true
		CharSetKey: n
`,
		},
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if test.node.Kind() != test.expectedKind {
				t.Errorf("expected kind %q, got %q", test.expectedKind, test.node.Kind())
			}

			actual := test.node.String()

			if actual != test.expectedStr {
				t.Errorf("expected %q, got %q", test.expectedStr, actual)
			}
		})
	}
}
