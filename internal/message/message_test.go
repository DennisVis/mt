// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package message_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/DennisVis/mt/internal/message"
	mttest "github.com/DennisVis/mt/testdata"
)

var ctx = context.Background()

func collectAllMessagesAndErrors(
	genericMessagesCh chan message.Message,
	errorsCh chan message.Error,
) ([]message.Message, []message.Error) {
	genericMessages := make([]message.Message, 0)
	errors := make([]message.Error, 0)

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

		for err := range errorsCh {
			errors = append(errors, err)
		}
	}()

	wg.Wait()

	return genericMessages, errors
}

func validateErrors(t *testing.T, expectedErrors, errors []message.Error) {
	if len(expectedErrors) < len(errors) {
		t.Errorf("expected %d parse errors, got %d: %s", len(expectedErrors), len(errors), errors)
		return
	}

	for i := 0; i < len(expectedErrors); i++ {
		if i+1 > len(errors) {
			t.Errorf("expected at least %d parse errors, got %d", i+1, len(errors))
			break
		}

		t.Run(fmt.Sprintf("ParseErrors[%d]", i+1), func(t *testing.T) {
			err := errors[i]
			expected := expectedErrors[i]

			switch {
			case expected.Err != nil && !strings.Contains(err.Error(), expected.Error()):
				t.Errorf("expected Error to be %q, got %q", expected.Err, err.Err)
			case expected.Line > 0 && err.Line != expected.Line:
				t.Errorf("expected Line to be %d, got %d", expected.Line, err.Line)
			}
		})
	}
}

func validateSubBlock(t *testing.T, name string, expected message.SubBlock, actual message.SubBlock) {
	t.Run(name, func(t *testing.T) {
		if expected.Label != "" && expected.Label != actual.Label {
			t.Errorf("expected label %q, got %q", expected.Label, actual.Label)
		}
		if expected.Content != "" && expected.Content != actual.Content {
			t.Errorf("expected content %q, got %q", expected.Content, actual.Content)
		}
	})
}

func validateSubBlocks(t *testing.T, name string, expected []message.SubBlock, actual []message.SubBlock) {
	if len(expected) != len(actual) {
		t.Errorf("expected %d sub-blocks, got %d", len(expected), len(actual))
	}
	for i, expectedSubBlock := range expected {
		actualSubBlock := actual[i]
		validateSubBlock(t, expectedSubBlock.Label, expectedSubBlock, actualSubBlock)
	}
}

func validateBlock(t *testing.T, name string, expected, actual message.Block) {
	t.Run(name, func(t *testing.T) {
		if expected.Label != "" && expected.Label != actual.Label {
			t.Errorf("expected label %q, got %q", expected.Label, actual.Label)
		}
		if expected.Content != "" && expected.Content != actual.Content {
			t.Errorf("expected content %q, got %q", expected.Content, actual.Content)
		}
		mttest.ValidateStringSliceMap(t, "Fields", expected.Fields, actual.Fields)
		validateSubBlocks(t, "Blocks", expected.Blocks, actual.Blocks)
	})
}

func validateBody(t *testing.T, expected, actual map[string][]string) {
	for k, vs := range expected {
		ovs, ok := actual[k]
		if !ok {
			t.Errorf("expected key %s in body, not found", k)
			continue
		}

		mttest.ValidateStringSlice(t, k, vs, ovs)
	}
}

func TestParse(t *testing.T) {
	for _, test := range []struct {
		name                string
		cfg                 message.Config
		input               io.Reader
		expectMessage       bool
		expectedErrors      []message.Error
		expectedBasicHeader *message.Block
		expectedAppHeader   *message.Block
		expectedUsrHeader   *message.Block
		expectedBody        *map[string][]string
		expectedTrailers    *message.Block
	}{
		{
			name:  "InvalidInput",
			input: &mttest.TestReaderInvalid{},
			expectedErrors: []message.Error{
				{
					Err: fmt.Errorf("could not read next rune from reader: invalid"),
				},
			},
		},
		{
			name:  "InvalidInputStopOnError",
			cfg:   message.Config{StopOnError: true},
			input: &mttest.TestReaderInvalid{},
			expectedErrors: []message.Error{
				{
					Err: fmt.Errorf("could not read next rune from reader: invalid"),
				},
			},
		},
		{
			name:          "BasicHeader",
			input:         strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}`),
			expectMessage: true,
			expectedBasicHeader: &message.Block{
				Label:   "1",
				Content: `F01SCBLZAJJXXXX5712100002`,
			},
		},
		{
			name:          "AppHeader",
			input:         strings.NewReader(`{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}`),
			expectMessage: true,
			expectedAppHeader: &message.Block{
				Label:   "2",
				Content: `O9401157091028SCBLZAJJXXXX57121000020910281157N`,
			},
		},
		{
			name:          "UsrHeader",
			input:         strings.NewReader(`{3:O9401157091028SCBLZAJJXXXX57121000020910281157N}`),
			expectMessage: true,
			expectedUsrHeader: &message.Block{
				Label:   "3",
				Content: `O9401157091028SCBLZAJJXXXX57121000020910281157N`,
			},
		},
		{
			name:          "TrailersAllCorrect",
			input:         strings.NewReader(`{5:{CHK:my checksum}{TNG:}{PDE:1348120811BANKFRPPAXXX2222123456}{DLM:}{MRF:1806271539180626BANKFRPPAXXX2222123456}{PDM:1213120811BANKFRPPAXXX2222123456}{SYS:1454120811BANKFRPPAXXX2222123456}}`),
			expectMessage: true,
			expectedTrailers: &message.Block{
				Label:   "5",
				Content: "",
				Blocks: []message.SubBlock{
					{
						Label:   "CHK",
						Content: "my checksum",
					},
					{
						Label:   "TNG",
						Content: "",
					},
					{
						Label:   "PDE",
						Content: "1348120811BANKFRPPAXXX2222123456",
					},
					{
						Label:   "DLM",
						Content: "",
					},
					{
						Label:   "MRF",
						Content: "1806271539180626BANKFRPPAXXX2222123456",
					},
					{
						Label:   "PDM",
						Content: "1213120811BANKFRPPAXXX2222123456",
					},
					{
						Label:   "SYS",
						Content: "1454120811BANKFRPPAXXX2222123456",
					},
				},
			},
		},
		{
			name: "BodyAllCorrect",
			input: strings.NewReader(`{4:
:20:Test1
:20a:Test2
:21:Test3
:21:Test4
-}`),
			expectMessage: true,
			expectedBody: &map[string][]string{
				"20":  {"Test1"},
				"20a": {"Test2"},
				"21":  {"Test3", "Test4"},
			},
		},
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			msgch, errch := message.Parse(ctx, test.input, test.cfg)
			msgs, errs := collectAllMessagesAndErrors(msgch, errch)
			validateErrors(t, test.expectedErrors, errs)

			if test.expectMessage && len(msgs) < 1 {
				t.Fatalf("expected at least 1 message, got 0")
			}

			if len(msgs) > 0 {
				if test.expectedBasicHeader != nil {
					validateBlock(t, "BasicHeader", *test.expectedBasicHeader, msgs[0].BasicHeader)
				}
				if test.expectedAppHeader != nil {
					validateBlock(t, "AppHeader", *test.expectedAppHeader, msgs[0].AppHeader)
				}
				if test.expectedUsrHeader != nil {
					validateBlock(t, "UsrHeader", *test.expectedUsrHeader, msgs[0].UsrHeader)
				}
				if test.expectedBody != nil {
					validateBody(t, *test.expectedBody, msgs[0].Body)
				}
				if test.expectedTrailers != nil {
					validateBlock(t, "Trailers", *test.expectedTrailers, msgs[0].Trailers)
				}
			}
		})
	}
}
