// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

//go:generate go run ./cmd/generate

// A library for parsing of MT (message text) messages according to the SWIFT specification.
package mt

import (
	"context"
	"io"
	"sync"

	"github.com/DennisVis/mt/internal/message"
)

// ParseMTx takes as input a reader and will attempt to parse all MT messages in the input and publish them to the
// returned channel. Any messages that cannot be parsed are discarded. The errors encountered during parsing are
// published on the returned parse error channel.
//
// This function returns generic MT messages, meaning the body is not parsed but simply returned as a
// map[string][]string. Look to the specialized derivatives for messages with fully parsed bodies.
//
// Using channels here means that potentially very large inputs can be read without running out of memory. If input is
// expected to easily fit into memory it is advised to use ParseAllMTx for convenience instead.
//
// Example usage:
//
//	f, err := os.Open("/path/to/mt/file.txt")
//	if err != nil {
//		return nil, fmt.Errorf("could not open file: %w", err)
//	}
//	defer f.Close()
//
//	messages, errors := ParseMTx(f)
//
//	// handle the errors from the errors channel
//	// process the messages from the messages channel
func ParseMTx(ctx context.Context, rd io.Reader, options ...option) (chan MTx, chan Error) {
	cfg := optionsToConfig(options)

	msgs, errs := message.Parse(ctx, rd, message.Config{
		StopOnError: cfg.StopOnError,
	})

	wg := &sync.WaitGroup{}
	mtxCh := make(chan MTx)
	errCh := make(chan Error)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for msg := range msgs {
			mtx, errs := messageToMTx(msg)
			if errs != nil {
				for _, err := range errs {
					errCh <- err
				}

				continue
			}

			mtxCh <- mtx
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for err := range errs {
			errCh <- NewError(err.Err, err.Line)
		}
	}()

	go func() {
		wg.Wait()
		close(mtxCh)
		close(errCh)
	}()

	return mtxCh, errCh
}

// ParseAllMTx takes as input a reader and will attempt to parse all MT messages in the input and return them to the
// caller. It's a convenience function for reading an entire input. If the input is expected to be very large, too large
// to fit in memory, use ParseMTx instead.
//
// In case of any errors during parsing a custom error is returned that encapsulates the parse errors. The presence of
// this error does not discredit the messages in the messages slice. Those were successfully parsed.
//
// This function returns generic MT messages, meaning the body is not parsed but simply returned as a map[string]string.
// Look to the specialized derivatives for messages with fully parsed bodies.
//
// Example usage:
//
//	f, err := os.Open("/path/to/mt/file.txt")
//	if err != nil {
//		return nil, fmt.Errorf("could not open file: %w", err)
//	}
//	defer f.Close()
//
//	messages, results, err := ParseAllMTx(f)
//	if err != nil {
//		// handle parse errors
//	}
//
// 	return messages, nil
func ParseAllMTx(ctx context.Context, rd io.Reader, options ...option) ([]MTx, error) {
	genericMessagesCh, parseErrorsCh := ParseMTx(ctx, rd, options...)

	genericMessages := make([]MTx, 0)
	parseErrors := make(Errors, 0)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for genericMessage := range genericMessagesCh {
			genericMessages = append(genericMessages, genericMessage)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for err := range parseErrorsCh {
			parseErrors = append(parseErrors, err)
		}
	}()

	wg.Wait()

	if len(parseErrors) > 0 {
		return genericMessages, parseErrors
	}

	return genericMessages, nil
}
