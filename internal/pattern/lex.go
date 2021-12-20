// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pattern

import (
	"unicode"
	"unicode/utf8"
)

// tokenType identifies the type of items the pattern lexer can produce.
type tokenType int

const (
	tokenEOF tokenType = iota
	tokenLiteral
	tokenOptionalLeftMeta
	tokenOptionalRightMeta
	tokenLineCountMeta
	tokenLineCount
	tokenCharCount
	tokenCharCountStrictMeta
	tokenCharSet
	tokenOrPatternMeta
)

var (
	optionalLeftMetaRune    = '('
	optionalRightMetaRune   = ')'
	lineCountMetaRune       = '*'
	charCountStrictMetaRune = '!'
	orPatternMetaRune       = '|'
)

const eof = -1

type token struct {
	typ tokenType // The type of this item.
	val string    // The value of this item.
}

type lexer struct {
	input  string     // the string being scanned
	pos    int        // current position in the input
	start  int        // start position of this item
	width  int        // width of last rune read from input
	tokens chan token // channel of scanned tokenS
}

func lex(input string) chan token {
	l := &lexer{
		input:  input,
		tokens: make(chan token),
	}

	go l.run()

	return l.tokens
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func() stateFn

func isCharSetSpecifier(r rune) bool {
	return charSetsKeys.contains(r)
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width

	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t tokenType) {
	l.tokens <- token{
		typ: t,
		val: l.input[l.start:l.pos],
	}

	l.start = l.pos
}

func (l *lexer) lexMeta(typ tokenType, next stateFn) stateFn {
	l.next()
	l.emit(typ)
	return next
}

func (l *lexer) lexCharSet() stateFn {
	l.next()
	l.emit(tokenCharSet)
	return l.lexToPattern
}

func (l *lexer) lexItemCharCountStrictMeta() stateFn {
	return l.lexMeta(tokenCharCountStrictMeta, l.lexCharSet)
}

func (l *lexer) lexItemLineCountMeta() stateFn {
	return l.lexMeta(tokenLineCountMeta, l.lexToPattern)
}

func (l *lexer) lexNumber() stateFn {
	for {
		switch r := l.next(); {
		case unicode.IsDigit(r):
			// consume
		// we've reached a line number meta, meaning the number we lexed is the line count
		// and we'll need to lex the line count meta
		case r == lineCountMetaRune:
			l.backup()
			l.emit(tokenLineCount)
			return l.lexItemLineCountMeta
		// we've reached a char count strict meta, meaning the number we lexed is the char count
		// and we'll need to lex the char count strict meta
		case r == charCountStrictMetaRune:
			l.backup()
			l.emit(tokenCharCount)
			return l.lexItemCharCountStrictMeta
		// we've reached a char set specifier, meaning the number we lexed is the char count
		// and we'll need to lex the char set specifier
		case isCharSetSpecifier(r):
			l.backup()
			l.emit(tokenCharCount)
			return l.lexCharSet
		// we've reached the end or an unknown character, meaning it wasn't known type we were parsing
		// it must be a literal instead, so we backup and parse it as a literal
		default:
			l.backup()
			return l.lexLiteral
		}
	}
}

func (l *lexer) lexLiteral() stateFn {
	digitsFound := 0

	for {
		switch r := l.next(); {
		case r == eof:
			l.emit(tokenLiteral)
			return l.lexToPattern
		case r == optionalLeftMetaRune:
			l.backup()
			l.emit(tokenLiteral)
			return l.lexOptionalLeftMeta
		case r == optionalRightMetaRune:
			l.backup()
			l.emit(tokenLiteral)
			return l.lexOptionalRightMeta
		case r == orPatternMetaRune:
			l.backup()
			l.emit(tokenLiteral)
			return l.lexOrPatternMeta
		// we've reached a reserved character meaning the number we were parsing is not a literal
		// we backup to the start of the number, emit the literal before it, and start to parse the number
		case digitsFound > 0 && (r == charCountStrictMetaRune || isCharSetSpecifier(r) || r == lineCountMetaRune):
			l.pos -= (digitsFound + 1)
			l.emit(tokenLiteral)
			return l.lexNumber
		case unicode.IsDigit(r):
			digitsFound++
		default:
			digitsFound = 0
		}
	}
}

func (l *lexer) lexOptionalRightMeta() stateFn {
	return l.lexMeta(tokenOptionalRightMeta, l.lexToPattern)
}

func (l *lexer) lexOptionalLeftMeta() stateFn {
	l.next()
	l.emit(tokenOptionalLeftMeta)
	return l.lexToPattern
}

func (l *lexer) lexOrPatternMeta() stateFn {
	return l.lexMeta(tokenOrPatternMeta, l.lexToPattern)
}

func (l *lexer) lexToPattern() stateFn {
Loop:
	for {
		switch r := l.next(); {
		case r == eof:
			break Loop
		case unicode.IsDigit(r):
			l.backup()
			return l.lexNumber
		case r == optionalLeftMetaRune:
			l.backup()
			return l.lexOptionalLeftMeta
		case r == optionalRightMetaRune:
			l.backup()
			return l.lexOptionalRightMeta
		case r == orPatternMetaRune:
			l.backup()
			return l.lexOrPatternMeta
		default:
			l.backup()
			return l.lexLiteral
		}
	}

	l.emit(tokenEOF) // Useful to make EOF a token.

	return nil // Stop the run loop.
}

func (l *lexer) run() {
	for state := l.lexToPattern; state != nil; {
		state = state()
	}

	close(l.tokens)
}
