// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pattern

import "github.com/DennisVis/mt/internal/pattern/ast"

type processor struct{}

func process(astPattern ast.Pattern) Pattern {
	p := &processor{}
	return p.astPatternToPattern(astPattern, 1)
}

func (p *processor) astLiteralToLiteral(l ast.Literal) Literal {
	return Literal{
		Chars: l.Value,
	}
}

func (p *processor) astOptionalToOptional(o ast.Optional, currLine int) Optional {
	return Optional{
		Pattern: p.astNodeToPartialValidator(o.Node, currLine),
	}
}

func (p *processor) astCharGroupToCharGroup(cg ast.CharGroup) CharGroup {
	return CharGroup{
		charSetKey:  cg.CharSetKey,
		CharSet:     charSets[cg.CharSetKey],
		Count:       cg.CharCount,
		CountStrict: cg.CharCountStrict,
	}
}

func (p *processor) astLineCountExpressionToLinePattern(lce ast.LineCountExpression, currLine int) LinePattern {
	endLine := currLine + lce.LineCount

	lp := LinePattern{
		InRange: func(line int) bool {
			return line >= currLine && line < endLine
		},
		Pattern: p.astNodeToPartialValidator(lce.Node, endLine),
	}

	return lp
}

func (p *processor) astOrExpressionToOrPattern(oe ast.OrExpression, currLine int) OrPattern {
	return OrPattern{
		Left:  p.astNodeToPartialValidator(oe.Left, currLine),
		Right: p.astNodeToPartialValidator(oe.Right, currLine),
	}
}

func (p *processor) astPatternToPattern(astPattern ast.Pattern, currLine int) Pattern {
	pattern := make(Pattern, len(astPattern.Nodes))

	for i, node := range astPattern.Nodes {
		pattern[i] = p.astNodeToPartialValidator(node, currLine)
	}

	return pattern
}

func (p *processor) astNodeToPartialValidator(node ast.Node, currLine int) ValidatesPartially {
	switch node.Kind() {
	case ast.NodeKindLiteral:
		return p.astLiteralToLiteral(node.(ast.Literal))
	case ast.NodeKindOptional:
		return p.astOptionalToOptional(node.(ast.Optional), currLine)
	case ast.NodeKindCharGroup:
		return p.astCharGroupToCharGroup(node.(ast.CharGroup))
	case ast.NodeKindLineCountExpression:
		return p.astLineCountExpressionToLinePattern(node.(ast.LineCountExpression), currLine)
	case ast.NodeKindOrExpression:
		return p.astOrExpressionToOrPattern(node.(ast.OrExpression), currLine)
	// astNodeKindPattern:
	default:
		return p.astPatternToPattern(node.(ast.Pattern), currLine)
	}
}
