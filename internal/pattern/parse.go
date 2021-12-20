// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pattern

import (
	"fmt"
	"strconv"

	"github.com/DennisVis/mt/internal/pattern/ast"
)

type parser struct {
	tokens         []token
	currTokenIndex int
	currToken      token
}

func parse(tokensCh chan token) (ast.Pattern, error) {
	var astPattern ast.Pattern

	tokens := make([]token, 0)

	for token := range tokensCh {
		tokens = append(tokens, token)
	}

	p := &parser{
		currTokenIndex: -1,
		tokens:         tokens,
	}

	p.next()

	astPattern, err := p.parsePattern()
	if err != nil {
		return astPattern, fmt.Errorf("could not create pattern ast: %w", err)
	}

	return astPattern, nil
}

func (p *parser) next() {
	if p.currTokenIndex < len(p.tokens)-1 {
		p.currTokenIndex++
		p.currToken = p.tokens[p.currTokenIndex]
	}
}

func (p *parser) peek() token {
	return p.tokens[p.currTokenIndex+1]
}

func (p *parser) backup() {
	if p.currTokenIndex > 0 {
		p.currTokenIndex--
		p.currToken = p.tokens[p.currTokenIndex]
	}
}

func (p *parser) parseOptional() (ast.Optional, error) {
	node := ast.Optional{}
	wrappedNodes := make([]ast.Node, 0)

	leftMetaFound := false
	lineCountFound := false

	var errToReturn error

Loop:
	for {
		switch p.currToken.typ {
		case tokenEOF:
			errToReturn = fmt.Errorf("parse optional: unclosed optional expression")
			break Loop
		case tokenOptionalLeftMeta:
			if leftMetaFound {
				optional, err := p.parseOptional()
				if err != nil {
					errToReturn = fmt.Errorf("parse optional: %w", err)
					break Loop
				}

				wrappedNodes = append(wrappedNodes, optional)

				p.next()
			} else {
				leftMetaFound = true
				p.next()
			}
		case tokenOptionalRightMeta:
			node.Node = ast.Pattern{
				Nodes: wrappedNodes,
			}
			break Loop
		case tokenLiteral:
			wrappedNodes = append(wrappedNodes, ast.Literal{
				Value: p.currToken.val,
			})

			p.next()
		case tokenCharCount:
			charGroup := p.parseCharGroup()

			wrappedNodes = append(wrappedNodes, charGroup)

			p.next()
		case tokenLineCount:
			// if we encounter a line count after a pattern
			// we need to wrap the pattern in a single line count expression
			if !lineCountFound && len(wrappedNodes) > 0 {
				wrappedNodes = []ast.Node{
					ast.LineCountExpression{
						LineCount: 1,
						Node: ast.Pattern{
							Nodes: wrappedNodes,
						},
					},
				}
			}

			lineCount, err := p.parseLineCount()
			if err != nil {
				errToReturn = fmt.Errorf("parse optional: %w", err)
				break Loop
			}

			wrappedNodes = append(wrappedNodes, lineCount)

			lineCountFound = true

			p.next()
		case tokenOrPatternMeta:
			or, err := p.parseOr(ast.Pattern{
				Nodes: wrappedNodes,
			})
			if err != nil {
				errToReturn = fmt.Errorf("parse optional: %w", err)
				break Loop
			}

			wrappedNodes = []ast.Node{or}

			p.next()
		}
	}

	return node, errToReturn
}

func (p *parser) parseCharGroup() ast.CharGroup {
	node := ast.CharGroup{}

Loop:
	for {
		switch p.currToken.typ {
		case tokenCharCount:
			// we know the lexer will only return a valid number, safe to ingore the error
			//nolint
			charCount, _ := strconv.Atoi(p.currToken.val)
			node.CharCount = charCount
			p.next()
		case tokenCharCountStrictMeta:
			node.CharCountStrict = true
			p.next()
		case tokenCharSet:
			node.CharSetKey = p.currToken.val
			break Loop
		}
	}

	return node
}

func (p *parser) parseOr(left ast.Node) (ast.OrExpression, error) {
	node := ast.OrExpression{
		Left: left,
	}
	rightNodes := make([]ast.Node, 0)
	ownTokenFound := false

	var errToReturn error

Loop:
	for {
		switch p.currToken.typ {
		// when we've reached the end of the input, or optional without returning
		// it means we've been collecting a regular pattern
		case tokenOptionalRightMeta:
			p.backup()
			fallthrough
		case tokenEOF:
			node.Right = ast.Pattern{
				Nodes: rightNodes,
			}
			break Loop
		case tokenOrPatternMeta:
			// if we've found our own or token and find another one it means
			if ownTokenFound {
				or, err := p.parseOr(ast.Pattern{
					Nodes: rightNodes,
				})
				if err != nil {
					errToReturn = fmt.Errorf("parse or: %w", err)
					break Loop
				}

				node.Right = or

				break Loop
			} else {
				ownTokenFound = true

				p.next()
			}
		case tokenLineCount:
			// a line count immediately after an or expression signals it to be the alternative to the preceding line
			// expression
			// the preceding line expression becomes the left side and the new line expression becomes the right side
			if len(rightNodes) == 0 {
				lineCount, err := p.parseLineCount()
				if err != nil {
					errToReturn = fmt.Errorf("parse or expression: %w", err)
					break Loop
				}

				node.Right = lineCount
			} else {
				// a line count after a pattern signals it to be an alternative to the preceding pattern
				// the preceding pattern becomes the left side and the collected nodes become the right side pattern
				node.Right = ast.Pattern{
					Nodes: rightNodes,
				}

				p.backup()
			}

			break Loop
		case tokenLiteral:
			rightNodes = append(rightNodes, ast.Literal{
				Value: p.currToken.val,
			})

			p.next()
		case tokenOptionalLeftMeta:
			optional, err := p.parseOptional()
			if err != nil {
				errToReturn = fmt.Errorf("parse or expression: %w", err)
				break Loop
			}

			rightNodes = append(rightNodes, optional)

			p.next()
		case tokenCharCount:
			charGroup := p.parseCharGroup()

			rightNodes = append(rightNodes, charGroup)

			p.next()
		}
	}

	return node, errToReturn
}

func (p *parser) parseLineCount() (ast.Node, error) {
	node := ast.LineCountExpression{}

	wrappedNodes := make([]ast.Node, 0)

	countSet := false

	var errToReturn error

Loop:
	for {
		switch p.currToken.typ {
		case tokenEOF:
			node.Node = ast.Pattern{
				Nodes: wrappedNodes,
			}

			break Loop
		case tokenLineCount:
			// we've encountered a new line count, we return the current one and prepare to parse the next
			if countSet {
				node.Node = ast.Pattern{
					Nodes: wrappedNodes,
				}

				p.backup()

				break Loop
			} else {
				// we know the lexer will only return a valid number, safe to ingore the error
				//nolint
				lineCount, _ := strconv.Atoi(p.currToken.val)
				node.LineCount = lineCount

				countSet = true

				p.next()
			}
		case tokenLineCountMeta:
			p.next()
		case tokenOptionalLeftMeta:
			optional, err := p.parseOptional()
			if err != nil {
				return node, fmt.Errorf("parse line count: %w", err)
			}

			wrappedNodes = append(wrappedNodes, optional)

			p.next()
		case tokenCharCount:
			charGroup := p.parseCharGroup()

			wrappedNodes = append(wrappedNodes, charGroup)

			p.next()
		case tokenOrPatternMeta:
			nextToken := p.peek()

			// if the next token is a line count it means that line count is being or'ed against the current one
			// in this case we finish the current line count and make it the left side of the or expression
			if nextToken.typ == tokenLineCount {
				node.Node = ast.Pattern{
					Nodes: wrappedNodes,
				}
				or, err := p.parseOr(node)
				if err != nil {
					errToReturn = fmt.Errorf("parse line count: %w", err)
					break Loop
				}

				return or, nil
			} else {
				or, err := p.parseOr(ast.Pattern{
					Nodes: wrappedNodes,
				})
				if err != nil {
					errToReturn = fmt.Errorf("parse line count: %w", err)
					break Loop
				}

				wrappedNodes = []ast.Node{or}

				p.next()
			}
		case tokenOptionalRightMeta:
			node.Node = ast.Pattern{
				Nodes: wrappedNodes,
			}

			p.backup()

			break Loop
		default:
			errToReturn = fmt.Errorf("parse line count: unexpected token %s", p.currToken.val)
			break Loop
		}
	}

	return node, errToReturn
}

func (p *parser) parsePattern() (ast.Pattern, error) {
	node := ast.Pattern{}
	var err error

	lineCountFound := false

Loop:
	for {
		switch p.currToken.typ {
		case tokenEOF:
			break Loop
		case tokenLiteral:
			node.Nodes = append(node.Nodes, ast.Literal{
				Value: p.currToken.val,
			})

			p.next()
		case tokenOptionalLeftMeta:
			optional, err := p.parseOptional()
			if err != nil {
				return node, fmt.Errorf("parse pattern: %w", err)
			}

			node.Nodes = append(node.Nodes, optional)

			p.next()
		case tokenCharCount:
			charGroup := p.parseCharGroup()

			node.Nodes = append(node.Nodes, charGroup)

			p.next()
		case tokenLineCount:
			if !lineCountFound && len(node.Nodes) > 0 {
				node.Nodes = []ast.Node{
					ast.LineCountExpression{
						LineCount: 1,
						Node:      node,
					},
				}
			}

			lineCount, err := p.parseLineCount()
			if err != nil {
				return node, fmt.Errorf("parse pattern: %w", err)
			}

			node.Nodes = append(node.Nodes, lineCount)

			lineCountFound = true

			p.next()
		case tokenOrPatternMeta:
			or, err := p.parseOr(node)
			if err != nil {
				return node, fmt.Errorf("parse pattern: %w", err)
			}

			node.Nodes = []ast.Node{or}

			p.next()
		}
	}

	return node, err
}
