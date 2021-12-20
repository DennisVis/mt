// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/ MIT

package mt_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/DennisVis/mt"
	mttest "github.com/DennisVis/mt/testdata"
)

const messageInput = `{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMN}{4:
:20:TELEWIZORY S.A.
:25:BPHKPLPK/320000546101
:28C:00084/001
:60F:C031002PLN40000,00
:61:0310201020C20000,00FMSCNONREF//8327000090031789
Card transaction
:86: 020?00Wyplata-(dysp/przel)?2008106000760000777777777777?2115617?
22INFO INFO INFO INFO INFO INFO 1 END?23INFO INFO INFO INFO INFO
INFO 2 END?24ZAPLATA ZA FABRYKATY DO TUB?25 - 200 S ZTUK, TRANZY
STORY-?26300 SZT GR544 I OPORNIKI-5?2700 SZT GTX847 FAKTURA 333/
2?28003.?3010600076?310000777777777777?32HUTA SZKLA TOPIC UL
PRZEMY?33SLOWA 67 32-669 WROCLAW?38PL081060007600007777777
77777
:61:0310201020D10000,00FTRFREF 25611247//8327000090031790
Transfer
:86: 020?00Wyplata-(dysp/przel)?2008106000760000777777777777?2115617?
22INFO INFO INFO INFO INFO INFO 1 END?23INFO INFO INFO INFO INFO
INFO 2 END?24ZAPLATA ZA FABRYKATY DO TUB?25 - 200 S ZTUK, TRANZY
STORY-?26300 SZT GR544 I OPORNIKI-5?2700 SZT GTX847 FAKTURA 333/
2?28003.?3010600076?310000777777777777?38PL081060007600007777777
77777
:61:0310201020C40,00FTRFNONREF//8327000090031791
Interest credit 
:86: 844?00Uznanie kwotą odsetek?20Odsetki od lokaty nr 101000?21022086
:62F:C020325PLN50040,00
-}`

func validateBody(t *testing.T, expected, actual map[string][]string) {
	t.Run("Body", func(t *testing.T) {
		for k, exp := range expected {
			act, ok := actual[k]
			if !ok {
				t.Errorf("expected body to contain %s, did not", k)
				return
			}

			mttest.ValidateStringSlice(t, k, exp, act)
		}
	})
}

func validateMTx(t *testing.T, expected, actual mt.MTx) {
	t.Run("MTx", func(t *testing.T) {
		if expected.BasicHeader.Raw != "" {
			mttest.ValidateBasicHeader(t, expected.BasicHeader, actual.BasicHeader)
		}
		if expected.AppHeaderInput.Set {
			mttest.ValidateAppHeaderInput(t, expected.AppHeaderInput, actual.AppHeaderInput)
		}
		if expected.AppHeaderOutput.Set {
			mttest.ValidateAppHeaderOutput(t, expected.AppHeaderOutput, actual.AppHeaderOutput)
		}
		if expected.UsrHeader.Set {
			mttest.ValidateUsrHeader(t, expected.UsrHeader, actual.UsrHeader)
		}
		if expected.Body != nil {
			validateBody(t, expected.Body, actual.Body)
		}
		if expected.Trailers.Set {
			mttest.ValidateTrailers(t, expected.Trailers, actual.Trailers)
		}
	})
}

func validateMTxs(t *testing.T, expected, actual []mt.MTx) {
	for i, exp := range expected {
		t.Run(fmt.Sprintf("MTxMessages[%d]", i), func(t *testing.T) {
			if i > len(actual)-1 {
				t.Errorf("expected at least %d messages, got only %d", i, len(actual))
				return
			}

			validateMTx(t, exp, actual[i])
		})
	}
}

func TestParseBasicHeader(t *testing.T) {
	for _, test := range []struct {
		name                string
		input               io.Reader
		expectedError       error
		expectedBasicHeader mt.BasicHeader
	}{
		{
			name:          "BasicHeaderTooShort",
			input:         strings.NewReader(`{1:122}{2:I940BOFAUS6BXBAMN}`),
			expectedError: errors.New("invalid basic header block content length"),
		},
		{
			name:  "BasicHeaderAppIDFinancial",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:I940BOFAUS6BXBAMN}`),
			expectedBasicHeader: mt.BasicHeader{
				Raw:   "{1:F01SCBLZAJJXXXX5712100002}",
				AppID: mt.ApplicationIDFinancial,
			},
		},
		{
			name:  "BasicHeaderAppIDGeneral",
			input: strings.NewReader(`{1:A01SCBLZAJJXXXX5712100002}{2:I940BOFAUS6BXBAMN}`),
			expectedBasicHeader: mt.BasicHeader{
				Raw:   "{1:A01SCBLZAJJXXXX5712100002}",
				AppID: mt.ApplicationIDGeneral,
			},
		},
		{
			name:  "BasicHeaderAppIDLogin",
			input: strings.NewReader(`{1:L01SCBLZAJJXXXX5712100002}{2:I940BOFAUS6BXBAMN}`),
			expectedBasicHeader: mt.BasicHeader{
				Raw:   "{1:L01SCBLZAJJXXXX5712100002}",
				AppID: mt.ApplicationIDLogin,
			},
		},
		{
			name:          "BasicHeaderAppIDUnknown",
			input:         strings.NewReader(`{1:X01SCBLZAJJXXXX5712100002}{2:I940BOFAUS6BXBAMN}`),
			expectedError: errors.New("invalid basic header: unknown application id in basic header block content: X"),
		},
		{
			name:  "BasicHeaderServiceIDFINGPA",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:I940BOFAUS6BXBAMN}`),
			expectedBasicHeader: mt.BasicHeader{
				Raw:       "{1:F01SCBLZAJJXXXX5712100002}",
				ServiceID: mt.ServiceIDFINGPA,
			},
		},
		{
			name:  "BasicHeaderServiceIDACKNACK",
			input: strings.NewReader(`{1:F21SCBLZAJJXXXX5712100002}{2:I940BOFAUS6BXBAMN}`),
			expectedBasicHeader: mt.BasicHeader{
				Raw:       "{1:F21SCBLZAJJXXXX5712100002}",
				ServiceID: mt.ServiceIDACKNACK,
			},
		},
		{
			name:          "BasicHeaderServiceIDUnknown",
			input:         strings.NewReader(`{1:FXXSCBLZAJJXXXX5712100002}{2:I940BOFAUS6BXBAMN}`),
			expectedError: errors.New("invalid basic header: unknown service id in basic header block content: XX"),
		},
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			msgs, err := mt.ParseAllMTx(ctx, test.input)
			mttest.ValidateError(t, test.expectedError, err)
			if len(msgs) != 0 {
				mttest.ValidateBasicHeader(t, test.expectedBasicHeader, msgs[0].BasicHeader)
			}
		})
	}
}

func TestParseAppHeader(t *testing.T) {
	for _, test := range []struct {
		name          string
		input         io.Reader
		expectedError error
	}{
		{
			name:          "TooShort",
			input:         strings.NewReader("{1:F01BPHKPLPKXXXX0000000000}{2:I94}"),
			expectedError: fmt.Errorf("invalid app header block content length"),
		},
		{
			name:          "WrongType",
			input:         strings.NewReader("{1:F01BPHKPLPKXXXX0000000000}{2:X940}"),
			expectedError: fmt.Errorf("invalid app header message type"),
		},
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := mt.ParseAllMTx(ctx, test.input)
			mttest.ValidateError(t, test.expectedError, err)
		})
	}
}

func TestParseAppHeaderInput(t *testing.T) {
	for _, test := range []struct {
		name                   string
		input                  io.Reader
		expectedError          mt.Error
		expectedAppHeaderInput mt.AppHeaderInput
	}{
		{
			name:          "AppHeaderInputTooShort",
			input:         strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBA}`),
			expectedError: mt.NewError(errors.New("invalid app header input block content length"), 1),
		},
		{
			name:          "AppHeaderInputTooLong",
			input:         strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMN2020X}`),
			expectedError: mt.NewError(errors.New("invalid app header input block content length"), 1),
		},
		{
			name:          "AppHeaderInputPriorityUnknown17",
			input:         strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMX}`),
			expectedError: mt.NewError(errors.New("invalid priority or delivery monitor in app header input"), 1),
		},
		{
			name:          "AppHeaderInputPriorityUnknown18",
			input:         strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMX2}`),
			expectedError: mt.NewError(errors.New("unknown message priority in app header input"), 1),
		},
		{
			name:  "AppHeaderInputPriorityS",
			input: strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMS}`),
			expectedAppHeaderInput: mt.AppHeaderInput{
				Raw:             "{2:I940BOFAUS6BXBAMS}",
				MessagePriority: mt.PrioritySystem,
			},
		},
		{
			name:  "AppHeaderInputPriorityU",
			input: strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMU}`),
			expectedAppHeaderInput: mt.AppHeaderInput{
				Raw:             "{2:I940BOFAUS6BXBAMU}",
				MessagePriority: mt.PriorityUrgent,
			},
		},
		{
			name:          "AppHeaderInputDeliveryMonitorUnknown17",
			input:         strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMX}`),
			expectedError: mt.NewError(errors.New("invalid priority or delivery monitor in app header input"), 1),
		},
		{
			name:          "AppHeaderInputDeliveryMonitorUnknown18",
			input:         strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMNX}`),
			expectedError: mt.NewError(errors.New("invalid delivery monitor in app header input"), 1),
		},
		{
			name:  "AppHeaderInputDeliveryMonitorUnknown1",
			input: strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMN1}`),
			expectedAppHeaderInput: mt.AppHeaderInput{
				Raw:             "{2:I940BOFAUS6BXBAMN1}",
				DeliveryMonitor: mt.DeliveryMonitorNonDelivery,
			},
		},
		{
			name:  "AppHeaderInputDeliveryMonitorUnknown2",
			input: strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMN2}`),
			expectedAppHeaderInput: mt.AppHeaderInput{
				Raw:             "{2:I940BOFAUS6BXBAMN2}",
				DeliveryMonitor: mt.DeliveryMonitorDelivery,
			},
		},
		{
			name:  "AppHeaderInputDeliveryMonitorUnknown3",
			input: strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAM3}`),
			expectedAppHeaderInput: mt.AppHeaderInput{
				Raw:             "{2:I940BOFAUS6BXBAM3}",
				DeliveryMonitor: mt.DeliveryMonitorBoth,
			},
		},
		{
			name:  "AppHeaderInputNoOptionalFields",
			input: strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAM}`),
		},
		{
			name:          "AppHeaderInputObsolescenceInvalidNumber",
			input:         strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAM02X}`),
			expectedError: mt.NewError(errors.New("invalid obsolescence period in app header input"), 1),
		},
		{
			name:  "AppHeaderInputObsolescenceValid",
			input: strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAM020}`),
			expectedAppHeaderInput: mt.AppHeaderInput{
				Raw:                         "{2:I940BOFAUS6BXBAM020}",
				ObsolescencePeriodInMinutes: 100,
			},
		},
		{
			name:          "AppHeaderInputObsolescenceValidAndPriorityInvalid",
			input:         strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMX020}`),
			expectedError: mt.NewError(errors.New("invalid priority or delivery monitor in app header input"), 1),
		},
		{
			name:          "AppHeaderInputObsolescenceInvalidAndPriorityValid",
			input:         strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMN02X}`),
			expectedError: mt.NewError(errors.New("invalid obsolescence period in app header input"), 1),
		},
		{
			name:  "AppHeaderInputObsolescenceAndPriorityValid",
			input: strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMU020}`),
			expectedAppHeaderInput: mt.AppHeaderInput{
				Raw:                         "{2:I940BOFAUS6BXBAMU020}",
				MessagePriority:             mt.PriorityUrgent,
				ObsolescencePeriodInMinutes: 100,
			},
		},
		{
			name:          "AppHeaderInputPriorityInvalidDeliveryMonitorValidObsolescenceValid",
			input:         strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMX2020}`),
			expectedError: mt.NewError(errors.New("unknown message priority in app header input"), 1),
		},
		{
			name:          "AppHeaderInputPriorityValidDeliveryMonitorInvalidObsolescenceValid",
			input:         strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMUX020}`),
			expectedError: mt.NewError(errors.New("invalid delivery monitor in app header input"), 1),
		},
		{
			name:          "AppHeaderInputPriorityValidDeliveryMonitorValidObsolescenceInvalid",
			input:         strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMU202X}`),
			expectedError: mt.NewError(errors.New("invalid obsolescence period in app header input"), 1),
		},
		{
			name:  "AppHeaderInputPriorityValidDeliveryMonitorValidObsolescenceValid",
			input: strings.NewReader(`{1:F01BPHKPLPKXXXX0000000000}{2:I940BOFAUS6BXBAMU2020}`),
			expectedAppHeaderInput: mt.AppHeaderInput{
				Raw:                         "{2:I940BOFAUS6BXBAMU2020}",
				MessagePriority:             mt.PriorityUrgent,
				DeliveryMonitor:             mt.DeliveryMonitorDelivery,
				ObsolescencePeriodInMinutes: 100,
			},
		},
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			msgs, err := mt.ParseAllMTx(ctx, test.input)
			if test.expectedError.Cause() != nil {
				mttest.ValidateErrors(t, test.expectedError, err)
			}
			if test.expectedError.Cause() == nil {
				mttest.ValidateAppHeaderInput(t, test.expectedAppHeaderInput, msgs[0].AppHeaderInput)
			}
		})
	}
}

func TestParseAppHeaderOutput(t *testing.T) {
	for _, test := range []struct {
		name                    string
		input                   io.Reader
		expectedError           mt.Error
		expectedAppHeaderOutput mt.AppHeaderOutput
	}{
		{
			name:          "TooShort",
			input:         strings.NewReader("{1:F01BPHKPLPKXXXX0000000000}{2:O9401157091028SCBLZAJJXXXX5712100002091028115}"),
			expectedError: mt.NewError(fmt.Errorf("invalid app header output block content length"), 1),
		},
		{
			name:          "InvalidInputTime",
			input:         strings.NewReader("{1:F01BPHKPLPKXXXX0000000000}{2:O9401X57091028SCBLZAJJXXXX57121000020910281157N}"),
			expectedError: mt.NewError(fmt.Errorf("invalid input time in app header output"), 1),
		},
		{
			name:          "InvalidOutputDate",
			input:         strings.NewReader("{1:F01BPHKPLPKXXXX0000000000}{2:O9401157091028SCBLZAJJXXXX571210000209X0281157N}"),
			expectedError: mt.NewError(fmt.Errorf("invalid output date in app header output"), 1),
		},
		{
			name:          "InvalidOutputTime",
			input:         strings.NewReader("{1:F01BPHKPLPKXXXX0000000000}{2:O9401157091028SCBLZAJJXXXX57121000020910281X57N}"),
			expectedError: mt.NewError(fmt.Errorf("invalid output time in app header output"), 1),
		},
		{
			name:          "InvalidMessageInputReference",
			input:         strings.NewReader("{1:F01BPHKPLPKXXXX0000000000}{2:O940115709X028SCBLZAJJXXXX57121000020910281157N}"),
			expectedError: mt.NewError(fmt.Errorf("could not parse message input reference with date"), 1),
		},
		{
			name:  "MessagePriorityN",
			input: strings.NewReader("{1:F01BPHKPLPKXXXX0000000000}{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}"),
			expectedAppHeaderOutput: mt.AppHeaderOutput{
				MessagePriority: mt.PriorityNormal,
			},
		},
		{
			name:  "MessagePriorityS",
			input: strings.NewReader("{1:F01BPHKPLPKXXXX0000000000}{2:O9401157091028SCBLZAJJXXXX57121000020910281157S}"),
			expectedAppHeaderOutput: mt.AppHeaderOutput{
				MessagePriority: mt.PrioritySystem,
			},
		},
		{
			name:  "MessagePriorityU",
			input: strings.NewReader("{1:F01BPHKPLPKXXXX0000000000}{2:O9401157091028SCBLZAJJXXXX57121000020910281157U}"),
			expectedAppHeaderOutput: mt.AppHeaderOutput{
				MessagePriority: mt.PriorityUrgent,
			},
		},
		{
			name:          "MessagePriorityUnknown",
			input:         strings.NewReader("{1:F01BPHKPLPKXXXX0000000000}{2:O9401157091028SCBLZAJJXXXX57121000020910281157X}"),
			expectedError: mt.NewError(fmt.Errorf("invalid message priority"), 1),
		},
	} {
		// rebind for parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			msgs, err := mt.ParseAllMTx(ctx, test.input)
			if test.expectedError.Cause() != nil {
				mttest.ValidateErrors(t, test.expectedError, err)
			}
			if len(msgs) > 0 {
				mttest.ValidateAppHeaderOutput(t, test.expectedAppHeaderOutput, msgs[0].AppHeaderOutput)
			}
		})
	}
}

// {1:F01SCBLZAJJXXXX5712100002}{2:I940BOFAUS6BXBAMN1}
func TestParseUsrHeader(t *testing.T) {
	for _, test := range []struct {
		name              string
		input             io.Reader
		expectedError     mt.Error
		expectedUsrHeader mt.UsrHeader
	}{
		{
			name:          "InvalidLabel",
			input:         strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:I940BOFAUS6BXBAMN1}{3:{555:123}}`),
			expectedError: mt.NewError(fmt.Errorf("invalid usr header block sub block label"), 1),
		},
		{
			name:          "InvalidMessageInputReference",
			input:         strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:I940BOFAUS6BXBAMN1}{3:{106:091X28SCBLZAJJXXXX57121000020}}`),
			expectedError: mt.NewError(fmt.Errorf("invalid message input reference"), 1),
		},
		{
			name:          "InvalidBalanceCheckpointDateTime",
			input:         strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:I940BOFAUS6BXBAMN1}{3:{423:123}}`),
			expectedError: mt.NewError(fmt.Errorf("invalid balance checkpoint time in usr header"), 1),
		},
		{
			name: "ValidAndComplete",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}
{2:I940BOFAUS6BXBAMN1}
{3:
{103:MyServiceID}
{106:120811BANKFRPPAXXX2222123456}
{108:MyUserReference}
{111:MyServiceTypeID}
{113:MyBankingPriority}
{115:MyAddressInformation}
{119:MyValidationFlag}
{121:MyUE2ETRef}
{165:MyPaymentReleaseInformation}
{423:060102150405000}
{424:MyRelatedReference}
{433:MySanctionsScreeningInformation}
{434:MyPaymentControlsInformation}
}`),
			expectedUsrHeader: mt.UsrHeader{
				Set:       true,
				ServiceID: "MyServiceID",
				MessageInputReference: mt.InputReference{
					Raw: "120811BANKFRPPAXXX2222123456",
				},
				MessageUserReference:               "MyUserReference",
				ServiceTypeID:                      "MyServiceTypeID",
				BankingPriority:                    "MyBankingPriority",
				AddresseeInformation:               "MyAddressInformation",
				ValidationFlag:                     "MyValidationFlag",
				UniqueEndToEndTransactionReference: "MyUE2ETRef",
				PaymentReleaseInformation:          "MyPaymentReleaseInformation",
				BalanceCheckpointDateTime: mt.DateTimeSecOptCent{
					Set: true,
					Raw: "060102150405000",
				},
				RelatedReference:              "MyRelatedReference",
				SanctionsScreeningInformation: "MySanctionsScreeningInformation",
				PaymentControlsInformation:    "MyPaymentControlsInformation",
			},
		},
	} {
		// rebing for parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			msgs, err := mt.ParseAllMTx(ctx, test.input)
			if test.expectedError.Cause() != nil {
				mttest.ValidateErrors(t, test.expectedError, err)
			} else if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if len(msgs) > 0 {
				mttest.ValidateUsrHeader(t, test.expectedUsrHeader, msgs[0].UsrHeader)
			}
		})
	}
}

func TestParseTrailers(t *testing.T) {
	for _, test := range []struct {
		name             string
		input            io.Reader
		expectedErrors   mt.Errors
		expectedTrailers mt.Trailers
	}{
		{
			name: "TrailersTooShort",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}{4:-}
		{5:{CHK:my checksum}{TNG:}{PDE:1348120811BANKFRPPAXXX222212345}{DLM:}{MRF:1806271539180626BANKFRPPAXXX222212345}{PDM:1213120811BANKFRPPAXXX222212345}{SYS:1454120811BANKFRPPAXXX222212345}}`),
			expectedErrors: mt.Errors{
				mt.NewError(errors.New("invalid possible duplicate emission"), 1),
				mt.NewError(errors.New("invalid message reference"), 1),
				mt.NewError(errors.New("invalid possible duplicate message"), 1),
				mt.NewError(errors.New("invalid system originated message"), 1),
			},
		},
		{
			name: "TrailersPDEInvalidTime",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}{4:-}
			{5:{PDE:13X8120811BANKFRPPAXXX2222123456}}`),
			expectedErrors: mt.Errors{
				mt.NewError(errors.New("invalid possible duplicate emission time"), 1),
			},
		},
		{
			name: "TrailersPDEInvalidInputReferenceDate",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}{4:-}
			{5:{PDE:131812X811BANKFRPPAXXX2222123456}}`),
			expectedErrors: mt.Errors{
				mt.NewError(errors.New("invalid message input reference with date date string"), 1),
			},
		},
		{
			name: "TrailersMRFInvaliDate",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}{4:-}
			{5:{MRF:18X6271539180626BANKFRPPAXXX2222123456}}`),
			expectedErrors: mt.Errors{
				mt.NewError(errors.New("invalid message reference date/time string"), 1),
			},
		},
		{
			name: "TrailersPDMInvalidTime",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}{4:-}
			{5:{PDM:12X3120811BANKFRPPAXXX2222123456}}`),
			expectedErrors: mt.Errors{
				mt.NewError(errors.New("invalid possible duplicate message time"), 1),
			},
		},
		{
			name: "TrailersPDMInvalidOutputReferenceDate",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}{4:-}
			{5:{PDM:12131208X1BANKFRPPAXXX2222123456}}`),
			expectedErrors: mt.Errors{
				mt.NewError(errors.New("invalid message output reference date/time string"), 1),
			},
		},
		{
			name: "TrailersPDMWithTimeInvalidTime",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}{4:-}
			{5:{PDM:134812X3120811BANKFRPPAXXX2222123456}}`),
			expectedErrors: mt.Errors{
				mt.NewError(errors.New("invalid message output reference date/time string"), 1),
			},
		},
		{
			name: "TrailersPDMWithTimeInvalidOutputReferenceDate",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}{4:-}
			{5:{PDM:134812131208X1BANKFRPPAXXX2222123456}}`),
			expectedErrors: mt.Errors{
				mt.NewError(errors.New("invalid message output reference date/time string"), 1),
			},
		},
		{
			name: "TrailersValidPDMWithTime",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}{4:-}
			{5:{PDM:12131213120811BANKFRPPAXXX2222123456}}`),
		},
		{
			name: "TrailersSYSInvalidTime",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}{4:-}
			{5:{SYS:14X4120811BANKFRPPAXXX2222123456}}`),
			expectedErrors: mt.Errors{
				mt.NewError(errors.New("invalid system originated message time"), 1),
			},
		},
		{
			name: "TrailersSYSInvalidMessageInputReference",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}{4:-}
			{5:{SYS:140412X811BANKFRPPAXXX2222123456}}`),
			expectedErrors: mt.Errors{
				mt.NewError(errors.New("invalid message input reference"), 1),
			},
		},
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			msgs, err := mt.ParseAllMTx(ctx, test.input)
			mttest.ValidateErrors(t, test.expectedErrors, err)
			if len(msgs) != 0 {
				mttest.ValidateTrailers(t, test.expectedTrailers, msgs[0].Trailers)
			}
		})
	}
}

func TestParseMTx(t *testing.T) {
	for _, test := range []struct {
		name           string
		input          io.Reader
		expectedErrors mt.Errors
		expectedMTxs   []mt.MTx
	}{
		{
			name: "ValidOutputMessage",
			input: strings.NewReader(`{1:F01SCBLZAJJXXXX5712100002}{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}{4:-}
		{5:{CHK:my checksum}{TNG:}{PDE:1348120811BANKFRPPAXXX2222123456}{DLM:}{MRF:1806271539180626BANKFRPPAXXX2222123456}{PDM:1213120811BANKFRPPAXXX2222123456}{SYS:1454120811BANKFRPPAXXX2222123456}}`),
			expectedMTxs: []mt.MTx{
				{
					Base: mt.Base{
						BasicHeader: mt.BasicHeader{
							Raw:                    "{1:F01SCBLZAJJXXXX5712100002}",
							AppID:                  mt.ApplicationIDFinancial,
							ServiceID:              mt.ServiceIDFINGPA,
							LogicalTerminalAddress: "SCBLZAJJXXXX",
							SessionNumber:          "5712",
							SequenceNumber:         "100002",
						},
						AppHeaderOutput: mt.AppHeaderOutput{
							Set:         true,
							Raw:         "{2:O9401157091028SCBLZAJJXXXX57121000020910281157N}",
							MessageType: "940",
							MessageInputReference: mt.InputReference{
								Set:                    true,
								Raw:                    "091028SCBLZAJJXXXX5712100002",
								LogicalTerminalAddress: "SCBLZAJJXXXX",
								SessionNumber:          "5712",
								SequenceNumber:         "100002",
								DateOrDateTime: mt.DateOrDateTime{
									Set: true,
									Raw: "091028",
								},
							},
							MessagePriority: mt.PriorityNormal,
						},
						Trailers: mt.Trailers{
							Set:                    true,
							Raw:                    "{5:{CHK:my checksum}{TNG:}{PDE:1348120811BANKFRPPAXXX2222123456}{DLM:}{MRF:1806271539180626BANKFRPPAXXX2222123456}{PDM:1213120811BANKFRPPAXXX2222123456}{SYS:1454120811BANKFRPPAXXX2222123456}}",
							DelayedMessage:         true,
							TestAndTrainingMessage: true,
							Checksum:               "my checksum",
							MessageReference: mt.Reference{
								DateTime: mttest.MustParseDateTime("1806271539"),
								MessageInputReference: mt.InputReference{
									DateOrDateTime:         mttest.MustParseDateOrDateTime("180626"),
									LogicalTerminalAddress: "BANKFRPPAXXX",
									SessionNumber:          "2222",
									SequenceNumber:         "123456",
								},
							},
							PossibleDuplicateEmission: mt.PossibleDuplicateEmission{
								Time: mttest.MustParseTime("1348"),
								MessageInputReference: mt.InputReference{
									DateOrDateTime:         mttest.MustParseDateOrDateTime("120811"),
									LogicalTerminalAddress: "BANKFRPPAXXX",
									SessionNumber:          "2222",
									SequenceNumber:         "123456",
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "ValidInputMessage",
			input: strings.NewReader(messageInput),
			expectedMTxs: []mt.MTx{
				{
					Base: mt.Base{
						BasicHeader: mt.BasicHeader{
							Raw:                    "{1:F01BPHKPLPKXXXX0000000000}",
							AppID:                  mt.ApplicationIDFinancial,
							ServiceID:              mt.ServiceIDFINGPA,
							LogicalTerminalAddress: "BPHKPLPKXXXX",
							SessionNumber:          "0000",
							SequenceNumber:         "000000",
						},
						AppHeaderInput: mt.AppHeaderInput{
							Set:             true,
							Raw:             "{2:I940BOFAUS6BXBAMN}",
							MessageType:     "940",
							ReceiverAddress: "BOFAUS6BXBAM",
							MessagePriority: mt.PriorityNormal,
						},
					},
					Body: map[string][]string{
						"20":  {"TELEWIZORY S.A."},
						"25":  {"BPHKPLPK/320000546101"},
						"28C": {"00084/001"},
						"60F": {"C031002PLN40000,00"},
						"61": {
							`0310201020C20000,00FMSCNONREF//8327000090031789
Card transaction`,
							`0310201020D10000,00FTRFREF 25611247//8327000090031790
Transfer`,
							`0310201020C40,00FTRFNONREF//8327000090031791
Interest credit`,
						},
						"86": {
							`020?00Wyplata-(dysp/przel)?2008106000760000777777777777?2115617?
22INFO INFO INFO INFO INFO INFO 1 END?23INFO INFO INFO INFO INFO
INFO 2 END?24ZAPLATA ZA FABRYKATY DO TUB?25 - 200 S ZTUK, TRANZY
STORY-?26300 SZT GR544 I OPORNIKI-5?2700 SZT GTX847 FAKTURA 333/
2?28003.?3010600076?310000777777777777?32HUTA SZKLA TOPIC UL
PRZEMY?33SLOWA 67 32-669 WROCLAW?38PL081060007600007777777
77777`,
							`020?00Wyplata-(dysp/przel)?2008106000760000777777777777?2115617?
22INFO INFO INFO INFO INFO INFO 1 END?23INFO INFO INFO INFO INFO
INFO 2 END?24ZAPLATA ZA FABRYKATY DO TUB?25 - 200 S ZTUK, TRANZY
STORY-?26300 SZT GR544 I OPORNIKI-5?2700 SZT GTX847 FAKTURA 333/
2?28003.?3010600076?310000777777777777?38PL081060007600007777777
77777`,
							`844?00Uznanie kwotą odsetek?20Odsetki od lokaty nr 101000?21022086`,
						},
						"62F": {"C020325PLN50040,00"},
					},
				},
			},
		},
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			msgs, err := mt.ParseAllMTx(ctx, test.input)
			mttest.ValidateErrors(t, test.expectedErrors, err)
			validateMTxs(t, test.expectedMTxs, msgs)
		})
	}
}

func BenchmarkParseMTxParallel(b *testing.B) {
	for _, msgCount := range []int{
		1,
		10,
		100,
		1000,
		10000,
	} {
		b.Run(fmt.Sprintf("MessageCount_%d", msgCount), func(b *testing.B) {
			messages := strings.Repeat(messageInput, msgCount)

			b.ResetTimer()

			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					mt.ParseMTx(ctx, strings.NewReader(messages))
				}
			})
		})
	}
}

func BenchmarkParseMTx(b *testing.B) {
	for _, msgCount := range []int{
		1,
		10,
		100,
		1000,
		10000,
	} {
		b.Run(fmt.Sprintf("MessageCount_%d", msgCount), func(b *testing.B) {
			messages := strings.Repeat(messageInput, msgCount)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				mt.ParseMTx(ctx, strings.NewReader(messages))
			}
		})
	}
}

func BenchmarkParseAllMTxParallel(b *testing.B) {
	for _, msgCount := range []int{
		1,
		10,
		100,
		1000,
		10000,
	} {
		b.Run(fmt.Sprintf("MessageCount_%d", msgCount), func(b *testing.B) {
			messages := strings.Repeat(messageInput, msgCount)

			b.ResetTimer()

			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					mt.ParseAllMTx(ctx, strings.NewReader(messages))
				}
			})
		})
	}
}

func BenchmarkParseAllMTx(b *testing.B) {
	for _, msgCount := range []int{
		1,
		10,
		100,
		1000,
		10000,
	} {
		b.Run(fmt.Sprintf("MessageCount_%d", msgCount), func(b *testing.B) {
			messages := strings.Repeat(messageInput, msgCount)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				mt.ParseAllMTx(ctx, strings.NewReader(messages))
			}
		})
	}
}
