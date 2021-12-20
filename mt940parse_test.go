// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mt_test

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/DennisVis/mt"
	mttest "github.com/DennisVis/mt/testdata"
)

var ctx = context.Background()

type TestMT940 struct {
	mt.MT940
	BasicHeader     mt.BasicHeader
	AppHeaderInput  mt.AppHeaderInput
	AppHeaderOutput mt.AppHeaderOutput
	UsrHeader       mt.UsrHeader
	Trailers        mt.Trailers
}

func (tmt TestMT940) toMT940() mt.MT940 {
	msg := tmt.MT940

	msg.BasicHeader = tmt.BasicHeader
	msg.AppHeaderInput = tmt.AppHeaderInput
	msg.AppHeaderOutput = tmt.AppHeaderOutput
	msg.UsrHeader = tmt.UsrHeader
	msg.Trailers = tmt.Trailers

	return msg
}

type TestMT940s []TestMT940

func (tmts TestMT940s) toMT940() []mt.MT940 {
	mts := make([]mt.MT940, len(tmts))
	for i, tmt := range tmts {
		mts[i] = tmt.toMT940()
	}
	return mts
}

func validateMT940s(t *testing.T, expectedMessages, messages []mt.MT940) {
	for i := 0; i < len(expectedMessages); i++ {
		if i+1 > len(messages) {
			t.Fatalf("expected at least %d parsed messages, got %d", i+1, len(messages))
			break
		}

		t.Run(fmt.Sprintf("MT940[%d]", i), func(t *testing.T) {
			expected := expectedMessages[i]
			actual := messages[i]

			mttest.ValidateBasicHeader(t, expected.BasicHeader, actual.BasicHeader)
			mttest.ValidateAppHeaderInput(t, expected.AppHeaderInput, actual.AppHeaderInput)
			mttest.ValidateAppHeaderOutput(t, expected.AppHeaderOutput, actual.AppHeaderOutput)
			mttest.ValidateUsrHeader(t, expected.UsrHeader, actual.UsrHeader)
			mttest.ValidateTrailers(t, expected.Trailers, actual.Trailers)
			mttest.ValidateBalance(t, "OpeningBalance", expected.OpeningBalance, actual.OpeningBalance)
			mttest.ValidateStatementLines(t, expected.StatementLines, actual.StatementLines)
			mttest.ValidateStringSlice(t, "AccountOwnerInformation", expected.AccountOwnerInformation, actual.AccountOwnerInformation)

			if expected.Reference != "" && expected.Reference != actual.Reference {
				t.Errorf("Reference expected %v, got %v", expected.Reference, actual.Reference)
			}
			if expected.AccountIdentification != "" && expected.AccountIdentification != actual.AccountIdentification {
				t.Errorf(
					"AccountIdentification expected %v, got %v",
					expected.AccountIdentification,
					actual.AccountIdentification,
				)
			}
			if expected.StatementNumberSequenceNumber != "" && expected.StatementNumberSequenceNumber != actual.StatementNumberSequenceNumber {
				t.Errorf(
					"StatementNumberSequenceNumber expected %v, got %v",
					expected.StatementNumberSequenceNumber,
					actual.StatementNumberSequenceNumber,
				)
			}
		})
	}
}

func TestParseMT940(t *testing.T) {
	for _, test := range []struct {
		name                string
		input               io.Reader
		expectedParseErrors mt.Errors
		expectedMT940s      TestMT940s
	}{
		{
			name:                "InvalidInput",
			input:               &mttest.TestReaderInvalid{},
			expectedParseErrors: []mt.Error{mt.NewError(mttest.ErrReadInvalid, 1)},
		},
		{
			name:  "SampleFile",
			input: mttest.MustOpenFile("testdata/sample-file-mt940.txt"),
			expectedMT940s: TestMT940s{
				{
					AppHeaderInput: mt.AppHeaderInput{
						MessagePriority: mt.PriorityNormal,
					},
					MT940: mt.MT940{
						Reference: "TELEWIZORY S.A.",
						OpeningBalance: mt.Balance{
							CreditDebit: mt.Credit,
							Date: mt.Date{
								Raw: "031002",
							},
							Currency: "PLN",
							Amount:   40000.00,
						},
					},
				},
			},
		},
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			msgs, err := mt.ParseAllMT940(ctx, test.input)
			mttest.ValidateErrors(t, test.expectedParseErrors, err)
			validateMT940s(t, test.expectedMT940s.toMT940(), msgs)
		})
	}
}
