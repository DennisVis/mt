// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package message

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

// itemType identifies the type of items the message lexer can produce.
type itemType int

const (
	itemError itemType = iota // error occurred; value is text of error
	itemEOF
	itemIgnore
	itemBlockLeftMeta
	itemBlockLabelMeta
	itemBlockLabel
	itemBlockContent
	itemBlockRightMeta
	itemSubBlockLeftMeta
	itemSubBlockLabelMeta
	itemSubBlockLabel
	itemSubBlockContent
	itemSubBlockRightMeta
	itemTagLeftMeta
	itemTagContent
	itemTagRightMeta
	itemFieldContent
)

var (
	blockLeftMeta     = "{"
	blockLabelMeta    = ":"
	blockRightMeta    = "}"
	subBlockLeftMeta  = "{"
	subBlockLabelMeta = ":"
	subBlockRightMeta = "}"
	tagLeftMeta       = ":"
	tagRightMeta      = ":"
	fieldsRightMeta   = "-}"
)

const eof = -1

type item struct {
	typ  itemType // The type of this item.
	val  string   // The value of this item.
	line int      // The line number at the start of this item.
}

// lexer holds the state of the scanner.
type lexer struct {
	ctx   context.Context
	input *bufio.Reader // the bytes being scanned
	buff  string        // the buffer used for storing read bytes from input
	items chan item     // channel of scanned items
	line  int           // start line of the current item
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func() stateFn

func newLexer(ctx context.Context, input *bufio.Reader) *lexer {
	l := &lexer{
		ctx:   ctx,
		input: input,
		items: make(chan item),
		line:  1,
	}

	go l.run()

	return l
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	i := item{
		typ:  t,
		val:  l.buff,
		line: l.line,
	}

	l.items <- i

	l.buff = ""
}

// errorf returns an error token and terminates the scan by passing back a nil pointer that will be the next state,
// terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		typ:  itemError,
		val:  fmt.Sprintf(format, args...),
		line: l.line,
	}
	return nil
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	r, _, err := l.input.ReadRune()
	if errors.Is(err, io.EOF) {
		return eof
	}
	if err != nil {
		l.errorf("could not read next rune from reader: %v", err)
		return eof
	}

	l.buff += string(r)

	if r == '\n' {
		l.line++
	}

	return r
}

func (l *lexer) lexText(typ itemType, next map[string]stateFn) stateFn {
	for {
		for suffix, nextStateFn := range next {
			if strings.HasSuffix(l.buff, suffix) {
				l.buff = l.buff[:len(l.buff)-len(suffix)]
				l.emit(typ)
				l.buff = suffix
				return nextStateFn
			}
		}

		if l.next() == eof {
			break
		}
	}

	// Correctly reached EOF.
	l.emit(typ)

	l.emit(itemEOF) // Useful to make EOF a token.

	return nil // Stop the run loop.
}

func (l *lexer) lexMeta(
	typ itemType,
	metaChars string,
	next stateFn,
) stateFn {
	l.emit(typ)
	return next
}

func (l *lexer) lexFieldContent() stateFn {
	return l.lexText(itemFieldContent, map[string]stateFn{
		// stop when we find a new tag and start parsing that
		tagLeftMeta: l.lexTagLeftMeta,
		// also stop when we find the end of the fields, we'll finish parsing of the block in that case
		fieldsRightMeta: l.lexBlockContent,
	})
}

func (l *lexer) lexTagRightMeta() stateFn {
	return l.lexMeta(
		itemTagRightMeta,
		tagRightMeta,
		l.lexFieldContent, // Now outside tag.
	)
}

func (l *lexer) lexTagContent() stateFn {
	return l.lexText(itemTagContent, map[string]stateFn{
		tagRightMeta: l.lexTagRightMeta,
	})
}

func (l *lexer) lexTagLeftMeta() stateFn {
	return l.lexMeta(
		itemTagLeftMeta,
		tagLeftMeta,
		l.lexTagContent, // Now inside tag.
	)
}

func (l *lexer) lexSubBlockRightMeta() stateFn {
	return l.lexMeta(
		itemSubBlockRightMeta,
		subBlockRightMeta,
		l.lexBlockContent, // we've reached the end of the sub block, we can now return to lexing the block
	)
}

func (l *lexer) lexBlockRightMeta() stateFn {
	return l.lexMeta(
		itemBlockRightMeta,
		blockLeftMeta,
		l.lexToBlock, // Now outside block, need to find new block.
	)
}

func (l *lexer) lexSubBlockContent() stateFn {
	return l.lexText(itemSubBlockContent, map[string]stateFn{
		// we've reached the end of the sub block, we can now return to lexing the block
		subBlockRightMeta: l.lexSubBlockRightMeta,
	})
}

func (l *lexer) lexSubBlockLabelMeta() stateFn {
	return l.lexMeta(
		itemSubBlockLabelMeta,
		subBlockLabelMeta,
		l.lexSubBlockContent,
	)
}

func (l *lexer) lexSubBlockLabel() stateFn {
	return l.lexText(itemSubBlockLabel, map[string]stateFn{
		subBlockLabelMeta: l.lexSubBlockLabelMeta,
	})
}

func (l *lexer) lexSubBlockLeftMeta() stateFn {
	return l.lexMeta(
		itemSubBlockLeftMeta,
		subBlockLeftMeta,
		l.lexSubBlockLabel, // Now inside subBlock.
	)
}

func (l *lexer) lexBlockContent() stateFn {
	return l.lexText(itemBlockContent, map[string]stateFn{
		blockRightMeta: l.lexBlockRightMeta,
		// a block can contain a sub-block, if it does we start parsing it
		subBlockLeftMeta: l.lexSubBlockLeftMeta,
		// a block can contain a tag, if it does we start parsing it
		tagLeftMeta: l.lexTagLeftMeta,
	})
}

func (l *lexer) lexBlockLabelMeta() stateFn {
	return l.lexMeta(
		itemBlockLabelMeta,
		blockLabelMeta,
		l.lexBlockContent,
	)
}

func (l *lexer) lexBlockLabel() stateFn {
	return l.lexText(itemBlockLabel, map[string]stateFn{
		blockLabelMeta: l.lexBlockLabelMeta,
	})
}

func (l *lexer) lexBlockLeftMeta() stateFn {
	return l.lexMeta(
		itemBlockLeftMeta,
		blockLeftMeta,
		l.lexBlockLabel, // Now inside block.
	)
}

func (l *lexer) lexToBlock() stateFn {
	return l.lexText(itemIgnore, map[string]stateFn{
		blockLeftMeta: l.lexBlockLeftMeta,
	})
}

// run lexes the input by executing state functions until the state is nil.
func (l *lexer) run() {
	state := l.lexToBlock

Loop:
	for {
		select {
		case <-l.ctx.Done():
		default:
			state = state()
			if state == nil {
				break Loop
			}
		}
	}

	close(l.items) // No more tokens will be delivered.
}
