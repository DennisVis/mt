// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package message

import (
	"fmt"
	"strings"
)

const (
	blockLabelBasicHeader = "1"
	blockLabelAppHeader   = "2"
	blockLabelUsrHeader   = "3"
	blockLabelBody        = "4"
	blockLabelTrailers    = "5"
)

type SubBlock struct {
	Label   string
	Content string
}

func newMessageSubBlock() SubBlock {
	return SubBlock{}
}

type Block struct {
	Label   string
	Content string
	Fields  map[string][]string
	Blocks  []SubBlock
}

func newBlock() Block {
	return Block{
		Fields: make(map[string][]string),
		Blocks: make([]SubBlock, 0),
	}
}

type parser struct {
	cfg        Config
	lexerItems chan item
	messages   chan Message
	errors     chan Error
}

func newParser(cfg Config, lexer *lexer) *parser {
	p := &parser{
		cfg:        cfg,
		lexerItems: lexer.items,
		messages:   make(chan Message),
		errors:     make(chan Error),
	}

	go p.run()

	return p
}

// blocksToMessage takes a slice of blocks, that should form a complete message, and parses them into a message struct.
// It delegates parsing of each type of blog to its respective function.
func (p *parser) blocksToMessage(blocks []Block, line int) Message {
	m := Message{
		Line: line,
	}

	var rawHeader string
	var rawAppHeader string
	var rawUsrHeader string
	var rawBody string
	var rawTrailers string

	for _, block := range blocks {
		switch block.Label {
		case blockLabelBasicHeader:
			m.BasicHeader = block
			rawHeader = fmt.Sprintf("{%s:%s}", blockLabelBasicHeader, block.Content)
		case blockLabelAppHeader:
			m.AppHeader = block
			rawAppHeader = fmt.Sprintf("{%s:%s}", blockLabelAppHeader, block.Content)
		case blockLabelUsrHeader:
			m.UsrHeader = block
			rawUsrHeader = fmt.Sprintf("{%s:%s}", blockLabelUsrHeader, block.Content)
		case blockLabelBody:
			m.Body = block.Fields
			rawBody = fmt.Sprintf("{%s:%s}", blockLabelBody, block.Content)
		case blockLabelTrailers:
			m.Trailers = block
			rawTrailers = fmt.Sprintf("{%s:%s}", blockLabelTrailers, block.Content)
		}
	}

	m.Raw = rawHeader + rawAppHeader + rawUsrHeader + rawBody + rawTrailers

	return m
}

// run runs the parser. This means it will read the items it receives from the lexer and parses them into complete
// messages.
func (p *parser) run() {
	blocks := make([]Block, 0)

	currLine := 1
	currBlock := newBlock()

	var currSubBlock SubBlock
	var currTag string

	sendMessage := func() {
		if len(blocks) > 0 {
			p.messages <- p.blocksToMessage(blocks, currLine)
		}
	}

Loop:
	for item := range p.lexerItems {
		switch item.typ {
		case itemBlockLabel:
			// if we receive a new basic header block it means a new message
			if item.val == blockLabelBasicHeader {
				// if we had blocks before this new message we process them before starting on the new message
				sendMessage()

				currLine = item.line
				blocks = make([]Block, 0)
			}

			currBlock = newBlock()
			currBlock.Label = item.val
		case itemBlockContent:
			currBlock.Content = item.val
		case itemSubBlockLeftMeta:
			currSubBlock = newMessageSubBlock()
		case itemSubBlockLabel:
			currSubBlock.Label = item.val
		case itemSubBlockContent:
			currSubBlock.Content = item.val
		case itemSubBlockRightMeta:
			currBlock.Blocks = append(currBlock.Blocks, currSubBlock)
		case itemTagContent:
			currTag = item.val
		case itemFieldContent:
			_, ok := currBlock.Fields[currTag]
			if !ok {
				currBlock.Fields[currTag] = make([]string, 0)
			}

			currBlock.Fields[currTag] = append(currBlock.Fields[currTag], strings.TrimSpace(item.val))
			currTag = ""
		case itemBlockRightMeta:
			blocks = append(blocks, currBlock)
		case itemError:
			p.errors <- Error{
				Err:  fmt.Errorf(item.val),
				Line: currLine,
			}
			if p.cfg.StopOnError {
				break Loop
			}
		case itemEOF:
			// If we've reached the end of the file and still have unprocessed blocks left these are processed as the
			// last message
			sendMessage()

			break Loop
		}
	}

	close(p.messages)
	close(p.errors)
}
