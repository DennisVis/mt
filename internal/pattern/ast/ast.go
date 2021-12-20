// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ast

import "fmt"

type NodeKind int

const (
	NodeKindPattern NodeKind = iota
	NodeKindLiteral
	NodeKindOptional
	NodeKindCharGroup
	NodeKindLineCountExpression
	NodeKindOrExpression
)

type Node interface {
	Kind() NodeKind
	IndentedString(indent string) string
	String() string
}

type Pattern struct {
	Nodes []Node
}

func (p Pattern) Kind() NodeKind {
	return NodeKindPattern
}

func (p Pattern) IndentedString(indent string) string {
	s := indent + "Pattern:\n"

	for _, n := range p.Nodes {
		s += n.IndentedString(indent + "	")
	}

	return s
}

func (p Pattern) String() string {
	return p.IndentedString("")
}

type Literal struct {
	Value string
}

func (l Literal) Kind() NodeKind {
	return NodeKindLiteral
}

func (l Literal) IndentedString(indent string) string {
	return fmt.Sprintf("%sLiteral:\n%sValue: %q\n", indent, indent+"	", l.Value)
}

func (l Literal) String() string {
	return l.IndentedString("")
}

type Optional struct {
	Node Node
}

func (o Optional) Kind() NodeKind {
	return NodeKindOptional
}

func (o Optional) IndentedString(indent string) string {
	return fmt.Sprintf(
		"%sOptional:\n%sNode:\n%s",
		indent,
		indent+"	",
		o.Node.IndentedString(indent+"		"),
	)
}

func (o Optional) String() string {
	return o.IndentedString("")
}

type CharGroup struct {
	CharCount       int
	CharCountStrict bool
	CharSetKey      string
}

func (cg CharGroup) Kind() NodeKind {
	return NodeKindCharGroup
}

func (cg CharGroup) IndentedString(indent string) string {
	return fmt.Sprintf(
		"%sCharGroup:\n%sCharCount: %d\n%sCharCountStrict: %v\n%sCharSetKey: %s\n",
		indent,
		indent+"	",
		cg.CharCount,
		indent+"	",
		cg.CharCountStrict,
		indent+"	",
		cg.CharSetKey,
	)
}

func (cg CharGroup) String() string {
	return cg.IndentedString("")
}

type LineCountExpression struct {
	LineCount int
	Node      Node
}

func (lce LineCountExpression) Kind() NodeKind {
	return NodeKindLineCountExpression
}

func (lce LineCountExpression) IndentedString(indent string) string {
	return fmt.Sprintf(
		"%sLineCountExpression:\n%sLineCount: %d\n%sNode:\n%s",
		indent,
		indent+"	",
		lce.LineCount,
		indent+"	",
		lce.Node.IndentedString(indent+"		"),
	)
}

func (lce LineCountExpression) String() string {
	return lce.IndentedString("")
}

type OrExpression struct {
	Left  Node
	Right Node
}

func (oe OrExpression) Kind() NodeKind {
	return NodeKindOrExpression
}

func (oe OrExpression) IndentedString(indent string) string {
	var left string
	if oe.Left != nil {
		left = oe.Left.IndentedString(indent + "		")
	} else {
		left = indent + "		" + "EMPTY\n"
	}

	var right string
	if oe.Right != nil {
		right = oe.Right.IndentedString(indent + "		")
	} else {
		right = indent + "		" + "EMPTY\n"
	}

	return fmt.Sprintf(
		"%sOrExpression:\n%sLeft:\n%s%sRight:\n%s",
		indent,
		indent+"	",
		left,
		indent+"	",
		right,
	)
}

func (oe OrExpression) String() string {
	return oe.IndentedString("")
}
