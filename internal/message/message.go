// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package message

import (
	"bufio"
	"context"
	"io"
)

type Config struct {
	StopOnError bool
}

type Message struct {
	Line        int
	Raw         string
	BasicHeader Block
	AppHeader   Block
	UsrHeader   Block
	Body        map[string][]string
	Trailers    Block
}

type Error struct {
	Line int
	Err  error
}

func (err Error) String() string {
	return err.Err.Error()
}

func (err Error) Error() string {
	return err.String()
}

func Parse(ctx context.Context, rd io.Reader, cfg Config) (chan Message, chan Error) {
	lexer := newLexer(ctx, bufio.NewReader(rd))
	parser := newParser(cfg, lexer)
	return parser.messages, parser.errors
}
