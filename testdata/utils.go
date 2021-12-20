// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package testdata

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/DennisVis/mt"
)

var ErrReadInvalid = fmt.Errorf("invalid")

type TestReaderInvalid struct{}

func (tr *TestReaderInvalid) Read(p []byte) (int, error) {
	for i := 0; i < len(p); i++ {
		p[i] = 'x'
	}

	return len(p), ErrReadInvalid
}

func MustParseTime(value string) mt.Time {
	var time mt.Time
	err := time.UnmarshalMT(value)
	if err != nil {
		panic(err)
	}

	return time
}

func MustParseDate(value string) mt.Date {
	var time mt.Date
	err := time.UnmarshalMT(value)
	if err != nil {
		panic(err)
	}

	return time
}

func MustParseDateTime(value string) mt.DateTime {
	var time mt.DateTime
	err := time.UnmarshalMT(value)
	if err != nil {
		panic(err)
	}

	return time
}

func MustParseDateOrDateTime(value string) mt.DateOrDateTime {
	var time mt.DateOrDateTime
	err := time.UnmarshalMT(value)
	if err != nil {
		panic(err)
	}

	return time
}

func MustOpenFile(path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return f
}

func ValidateError(t *testing.T, expected, actual error) {
	formatError := func(err error) string {
		if err == nil {
			return ""
		}

		return strings.Join(strings.Split(err.Error(), "\n"), "\n\t")
	}

	t.Run("Error", func(t *testing.T) {
		switch {
		case expected == nil && actual != nil:
			t.Errorf("\nexpected nil error, got:\n\t%s", formatError(actual))
		case expected != nil && actual == nil:
			t.Errorf("\nexpected error:\n\t%s, got nil", formatError(expected))
		case expected != nil && actual != nil && !strings.Contains(actual.Error(), expected.Error()):
			t.Errorf("\nexpected error:\n\t%s\ngot:\n\t%s", formatError(expected), formatError(actual))
		}
	})
}

func ValidateErrors(t *testing.T, e, a error) {
	expected, eok := e.(mt.Errors)
	if !eok {
		expectedErr, ok := e.(mt.Error)
		if ok {
			expected = mt.Errors{expectedErr}
			eok = true
		}
	}

	actual, aok := a.(mt.Errors)

	switch {
	case expected == nil && actual == nil:
		return
	case expected.Error() == "" && actual == nil:
		return
	case expected.Error() == "" && actual.Error() == "":
		return
	case expected == nil && actual != nil:
		t.Errorf("expected nil error, got: %s", actual)
	case expected != nil && actual == nil:
		t.Errorf("expected error: %s, got nil", expected)
	case !eok && aok:
		t.Errorf("expected not of type mt.Errors: %T", expected)
	case eok && !aok:
		t.Errorf("actual not of type mt.Errors: %T", actual)
	}

	if len(expected) < len(actual) {
		t.Errorf("expected %d parse errors, got %d: %s", len(expected), len(actual), actual)
		return
	}

	for i := 0; i < len(expected); i++ {
		if i+1 > len(actual) {
			t.Errorf("expected at least %d parse errors, got %d", i+1, len(actual))
			break
		}

		t.Run(fmt.Sprintf("ParseErrors[%d]", i), func(t *testing.T) {
			exp := expected[i]
			act := actual[i]

			switch {
			case exp.Cause() != nil && !strings.Contains(act.Cause().Error(), exp.Cause().Error()):
				t.Errorf("expected Error to be %q, got %q", exp.Error(), act.Error())
			case exp.Line() > 0 && act.Line() != exp.Line():
				t.Errorf("expected Line to be %d, got %d", exp.Line, act.Line)
			}
		})
	}
}

func ValidateStringMap(t *testing.T, name string, expected, actual map[string]string) {
	for k, v := range expected {
		t.Run(fmt.Sprintf(name+"[%s]", k), func(t *testing.T) {
			if actual[k] != v {
				t.Errorf("expected %v, got %v", v, actual[k])
			}
		})
	}
}

func ValidateStringSlice(t *testing.T, name string, expected, actual []string) {
	for i, v := range expected {
		t.Run(fmt.Sprintf(name+"[%d]", i), func(t *testing.T) {
			if actual[i] != v {
				t.Errorf("expected %v, got %v", v, actual[i])
			}
		})
	}
}

func ValidateStringSliceMap(t *testing.T, name string, expected, actual map[string][]string) {
	for k, v := range expected {
		ValidateStringSlice(t, name+"[%s]", v, actual[k])
	}
}

func ValidateRaw(t *testing.T, expected, actual string) {
	t.Run("Raw", func(t *testing.T) {
		if expected != "" && expected != actual {
			t.Errorf("expected Raw %s, got %s", expected, actual)
		}
	})
}

func ValidateTime(t *testing.T, expected, actual mt.Time) {
	t.Run("Time", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if !expected.Time.IsZero() && expected.Time != actual.Time {
			t.Errorf("expected Time %s, got %s", expected.Time, actual.Time)
		}
	})
}

func ValidateMonth(t *testing.T, name string, expected, actual mt.Month) {
	t.Run(name, func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if !expected.Time.IsZero() && expected.Time != actual.Time {
			t.Errorf("expected Time %s, got %s", expected.Time, actual.Time)
		}
	})
}

func ValidateDate(t *testing.T, expected, actual mt.Date) {
	t.Run("Date", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if !expected.Time.IsZero() && expected.Time != actual.Time {
			t.Errorf("expected Time %s, got %s", expected.Time, actual.Time)
		}
	})
}

func ValidateDateTime(t *testing.T, expected, actual mt.DateTime) {
	t.Run("DateTime", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if !expected.Time.IsZero() && expected.Time != actual.Time {
			t.Errorf("expected Time %s, got %s", expected.Time, actual.Time)
		}
	})
}

func ValidateDateOrDateTime(t *testing.T, expected, actual mt.DateOrDateTime) {
	t.Run("DateOrDateTime", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if !expected.Time.IsZero() && expected.Time != actual.Time {
			t.Errorf("expected Time %s, got %s", expected.Time, actual.Time)
		}
	})
}

func ValidateDateTimeSecOptCent(t *testing.T, expected, actual mt.DateTimeSecOptCent) {
	t.Run("DateTimeSecOptCent", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if !expected.Time.IsZero() && expected.Time != actual.Time {
			t.Errorf("expected Time %s, got %s", expected.Time, actual.Time)
		}
	})
}

func ValidateOutputReference(t *testing.T, expected, actual mt.OutputReference) {
	t.Run("OutputReference", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if expected.LogicalTerminalAddress != "" && expected.LogicalTerminalAddress != actual.LogicalTerminalAddress {
			t.Errorf(
				"LogicalTerminalAddress expected %v, got %v",
				expected.LogicalTerminalAddress,
				actual.LogicalTerminalAddress,
			)
		}
		if expected.SessionNumber != "" && expected.SessionNumber != actual.SessionNumber {
			t.Errorf("SessionNumber expected %v, got %v", expected.SessionNumber, actual.SessionNumber)
		}
		if expected.SequenceNumber != "" && expected.SequenceNumber != actual.SequenceNumber {
			t.Errorf("SequenceNumber expected %v, got %v", expected.SequenceNumber, actual.SequenceNumber)
		}
		ValidateDateOrDateTime(t, expected.DateOrDateTime, actual.DateOrDateTime)
	})
}

func ValidateInputReference(t *testing.T, expected, actual mt.InputReference) {
	t.Run("InputReference", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if expected.LogicalTerminalAddress != "" && expected.LogicalTerminalAddress != actual.LogicalTerminalAddress {
			t.Errorf(
				"LogicalTerminalAddress expected %v, got %v",
				expected.LogicalTerminalAddress,
				actual.LogicalTerminalAddress,
			)
		}
		if expected.SessionNumber != "" && expected.SessionNumber != actual.SessionNumber {
			t.Errorf("SessionNumber expected %v, got %v", expected.SessionNumber, actual.SessionNumber)
		}
		if expected.SequenceNumber != "" && expected.SequenceNumber != actual.SequenceNumber {
			t.Errorf("SequenceNumber expected %v, got %v", expected.SequenceNumber, actual.SequenceNumber)
		}
		ValidateDateOrDateTime(t, expected.DateOrDateTime, actual.DateOrDateTime)
	})
}

func ValidateReference(t *testing.T, expected, actual mt.Reference) {
	t.Run("Reference", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		ValidateDateTime(t, expected.DateTime, actual.DateTime)
		ValidateInputReference(t, expected.MessageInputReference, actual.MessageInputReference)
	})
}

func ValidatePossibleDuplicateEmission(t *testing.T, expected, actual mt.PossibleDuplicateEmission) {
	t.Run("PossibleDuplicateEmission", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		ValidateTime(t, expected.Time, actual.Time)
		ValidateInputReference(t, expected.MessageInputReference, actual.MessageInputReference)
	})
}

func ValidatePossibleDuplicateMessage(t *testing.T, expected, actual mt.PossibleDuplicateMessage) {
	t.Run("PossibleDuplicateMessage", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		ValidateTime(t, expected.Time, actual.Time)
		ValidateOutputReference(t, expected.MessageOutputReference, actual.MessageOutputReference)
	})
}

func ValidateSystemOriginatedMessage(t *testing.T, expected, actual mt.SystemOriginatedMessage) {
	t.Run("SystemOriginatedMessage", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		ValidateTime(t, expected.Time, actual.Time)
		ValidateInputReference(t, expected.MessageInputReference, actual.MessageInputReference)
	})
}

func ValidateBasicHeader(t *testing.T, expected, actual mt.BasicHeader) {
	t.Run("BasicHeader", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if expected.AppID != actual.AppID {
			t.Errorf("AppID expected %v, got %v", expected.AppID, actual.AppID)
		}
		if expected.ServiceID != actual.ServiceID {
			t.Errorf("ServiceID expected %v, got %v", expected.ServiceID, actual.ServiceID)
		}
		if expected.SessionNumber != "" && expected.SessionNumber != actual.SessionNumber {
			t.Errorf("SessionNumber expected %v, got %v", expected.SessionNumber, actual.SessionNumber)
		}
		if expected.SequenceNumber != "" && expected.SequenceNumber != actual.SequenceNumber {
			t.Errorf("SequenceNumber expected %v, got %v", expected.SequenceNumber, actual.SequenceNumber)
		}
		if expected.LogicalTerminalAddress != "" && expected.LogicalTerminalAddress != actual.LogicalTerminalAddress {
			t.Errorf(
				"LogicalTerminalAddress expected %v, got %v",
				expected.LogicalTerminalAddress,
				actual.LogicalTerminalAddress,
			)
		}
	})
}

func ValidateAppHeaderInput(t *testing.T, expected, actual mt.AppHeaderInput) {
	t.Run("AppHeaderInput", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if expected.ObsolescencePeriodInMinutes != 0 && expected.ObsolescencePeriodInMinutes != actual.ObsolescencePeriodInMinutes {
			t.Errorf(
				"ObsolescencePeriodInMinutes expected %v, got %v",
				expected.ObsolescencePeriodInMinutes,
				actual.ObsolescencePeriodInMinutes,
			)
		}
		if expected.MessageType != "" && expected.MessageType != actual.MessageType {
			t.Errorf("MessageType expected %v, got %v", expected.MessageType, actual.MessageType)
		}
		if expected.ReceiverAddress != "" && expected.ReceiverAddress != actual.ReceiverAddress {
			t.Errorf("ReceiverAddress expected %v, got %v", expected.ReceiverAddress, actual.ReceiverAddress)
		}
		if expected.MessagePriority != actual.MessagePriority {
			t.Errorf("MessagePriority expected %v, got %v", expected.MessagePriority, actual.MessagePriority)
		}
		if expected.DeliveryMonitor != actual.DeliveryMonitor {
			t.Errorf("DeliveryMonitor expected %v, got %v", expected.DeliveryMonitor, actual.DeliveryMonitor)
		}
	})
}

func ValidateAppHeaderOutput(t *testing.T, expected, actual mt.AppHeaderOutput) {
	t.Run("AppHeaderOutput", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if expected.MessageType != "" && expected.MessageType != actual.MessageType {
			t.Errorf("MessageType expected %v, got %v", expected.MessageType, actual.MessageType)
		}
		if expected.MessagePriority != actual.MessagePriority {
			t.Errorf("MessagePriority expected %v, got %v", expected.MessagePriority, actual.MessagePriority)
		}
		ValidateTime(t, expected.InputTime, actual.InputTime)
		ValidateDate(t, expected.OutputDate, actual.OutputDate)
		ValidateTime(t, expected.OutputTime, actual.OutputTime)
		ValidateInputReference(t, expected.MessageInputReference, actual.MessageInputReference)
	})
}

func ValidateUsrHeader(t *testing.T, expected, actual mt.UsrHeader) {
	t.Run("UsrHeader", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if expected.ServiceID != "" && expected.ServiceID != actual.ServiceID {
			t.Errorf("ServiceID expected %v, got %v", expected.ServiceID, actual.ServiceID)
		}
		if expected.AddresseeInformation != "" && expected.AddresseeInformation != actual.AddresseeInformation {
			t.Errorf(
				"AddresseeInformation expected %v, got %v",
				expected.AddresseeInformation,
				actual.AddresseeInformation,
			)
		}
		if expected.BankingPriority != "" && expected.BankingPriority != actual.BankingPriority {
			t.Errorf(
				"BankingPriority expected %v, got %v",
				expected.BankingPriority,
				actual.BankingPriority,
			)
		}
		if expected.MessageUserReference != "" && expected.MessageUserReference != actual.MessageUserReference {
			t.Errorf(
				"MessageUserReference expected %v, got %v",
				expected.MessageUserReference,
				actual.MessageUserReference,
			)
		}
		if expected.ValidationFlag != "" && expected.ValidationFlag != actual.ValidationFlag {
			t.Errorf(
				"ValidationFlag expected %v, got %v",
				expected.ValidationFlag,
				actual.ValidationFlag,
			)
		}
		if expected.RelatedReference != "" && expected.RelatedReference != actual.RelatedReference {
			t.Errorf(
				"RelatedReference expected %v, got %v",
				expected.RelatedReference,
				actual.RelatedReference,
			)
		}
		if expected.ServiceTypeID != "" && expected.ServiceTypeID != actual.ServiceTypeID {
			t.Errorf(
				"ServiceTypeID expected %v, got %v",
				expected.ServiceTypeID,
				actual.ServiceTypeID,
			)
		}
		if expected.UniqueEndToEndTransactionReference != "" && expected.UniqueEndToEndTransactionReference != actual.UniqueEndToEndTransactionReference {
			t.Errorf(
				"UniqueEndToEndTransactionReference expected %v, got %v",
				expected.UniqueEndToEndTransactionReference,
				actual.UniqueEndToEndTransactionReference,
			)
		}
		if expected.PaymentReleaseInformation != "" && expected.PaymentReleaseInformation != actual.PaymentReleaseInformation {
			t.Errorf(
				"PaymentReleaseInformation expected %v, got %v",
				expected.PaymentReleaseInformation,
				actual.PaymentReleaseInformation,
			)
		}
		if expected.SanctionsScreeningInformation != "" && expected.SanctionsScreeningInformation != actual.SanctionsScreeningInformation {
			t.Errorf(
				"SanctionsScreeningInformation expected %v, got %v",
				expected.SanctionsScreeningInformation,
				actual.SanctionsScreeningInformation,
			)
		}
		if expected.PaymentControlsInformation != "" && expected.PaymentControlsInformation != actual.PaymentControlsInformation {
			t.Errorf(
				"PaymentControlsInformation expected %v, got %v",
				expected.PaymentControlsInformation,
				actual.PaymentControlsInformation,
			)
		}
		ValidateDateTimeSecOptCent(t, expected.BalanceCheckpointDateTime, actual.BalanceCheckpointDateTime)
		ValidateInputReference(t, expected.MessageInputReference, actual.MessageInputReference)
	})
}

func ValidateTrailers(t *testing.T, expected, actual mt.Trailers) {
	t.Run("Trailers", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if expected.DelayedMessage != actual.DelayedMessage {
			t.Errorf("DelayedMessage expected %v, got %v", expected.DelayedMessage, actual.DelayedMessage)
		}
		if expected.TestAndTrainingMessage != actual.TestAndTrainingMessage {
			t.Errorf(
				"TestAndTrainingMessage expected %v, got %v",
				expected.TestAndTrainingMessage,
				actual.TestAndTrainingMessage,
			)
		}
		if expected.Checksum != "" && expected.Checksum != actual.Checksum {
			t.Errorf("Checksum expected %v, got %v", expected.Checksum, actual.Checksum)
		}
		ValidateReference(t, expected.MessageReference, actual.MessageReference)
		ValidatePossibleDuplicateEmission(t, expected.PossibleDuplicateEmission, actual.PossibleDuplicateEmission)
		ValidatePossibleDuplicateMessage(t, expected.PossibleDuplicateMessage, actual.PossibleDuplicateMessage)
		ValidateSystemOriginatedMessage(t, expected.SystemOriginatedMessage, actual.SystemOriginatedMessage)
		ValidateStringMap(t, "AdditionalTrailers", expected.AdditionalTrailers, actual.AdditionalTrailers)
	})
}

func ValidateBalance(t *testing.T, name string, expected, actual mt.Balance) {
	t.Run(name, func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if expected.CreditDebit != actual.CreditDebit {
			t.Errorf("expected credit/debit %v, got %v", expected.CreditDebit, actual.CreditDebit)
		}
		if expected.Currency != "" && expected.Currency != actual.Currency {
			t.Errorf("expected currency %s, got %s", expected.Currency, actual.Currency)
		}
		if expected.Amount != actual.Amount {
			t.Errorf("expected amount %f, got %f", expected.Amount, actual.Amount)
		}
		ValidateDate(t, expected.Date, actual.Date)
	})
}

func ValidateStatementLine(t *testing.T, expected, actual mt.StatementLine) {
	t.Run("StatementLine", func(t *testing.T) {
		ValidateRaw(t, expected.Raw, actual.Raw)
		if expected.FundsCode != actual.FundsCode {
			t.Errorf("expected funds code %s, got %s", expected.FundsCode, actual.FundsCode)
		}
		if expected.Amount != actual.Amount {
			t.Errorf("expected amount %f, got %f", expected.Amount, actual.Amount)
		}
		if expected.SwiftCode != "" && expected.SwiftCode != actual.SwiftCode {
			t.Errorf("expected SWIFT code %s, got %s", expected.SwiftCode, actual.SwiftCode)
		}
		if expected.AccountOwnerReference != "" && expected.AccountOwnerReference != actual.AccountOwnerReference {
			t.Errorf("expected account owner reference %s, got %s", expected.AccountOwnerReference, actual.AccountOwnerReference)
		}
		if expected.BankReference != "" && expected.BankReference != actual.BankReference {
			t.Errorf("expected bank reference %s, got %s", expected.BankReference, actual.BankReference)
		}
		if expected.Description != "" && expected.Description != actual.Description {
			t.Errorf("expected description %s, got %s", expected.Description, actual.Description)
		}
		ValidateDate(t, expected.Date, actual.Date)
		ValidateMonth(t, "EntryDate", expected.EntryDate, actual.EntryDate)
	})
}

func ValidateStatementLines(t *testing.T, expected, actual []mt.StatementLine) {
	for i, sl := range expected {
		t.Run(fmt.Sprintf("StatementLine[%d]", i), func(t *testing.T) {
			ValidateStatementLine(t, sl, actual[i])
		})
	}
}
