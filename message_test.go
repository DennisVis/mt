// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mt_test

import (
	"fmt"
	"testing"

	"github.com/DennisVis/mt"
	mttest "github.com/DennisVis/mt/testdata"
)

func TestApplicationID(t *testing.T) {
	t.Parallel()

	if mt.ApplicationIDFinancial.RawString() != "F" {
		t.Error("ApplicationIDFinancial raw string is not F")
	}
	if mt.ApplicationIDGeneral.RawString() != "A" {
		t.Error("ApplicationIDGeneral raw string is not A")
	}
	if mt.ApplicationIDLogin.RawString() != "L" {
		t.Error("ApplicationIDLogin raw string is not L")
	}
}

func TestServiceID(t *testing.T) {
	t.Parallel()

	if mt.ServiceIDFINGPA.RawString() != "01" {
		t.Error("ServiceIDFINGPA raw string is not 01")
	}
	if mt.ServiceIDACKNACK.RawString() != "21" {
		t.Error("ServiceIDACKNACK raw string is not 21")
	}
}

func TestPriority(t *testing.T) {
	t.Parallel()

	if mt.PrioritySystem.RawString() != "S" {
		t.Error("PrioritySystem raw string is not S")
	}
	if mt.PriorityNormal.RawString() != "N" {
		t.Error("PriorityNormal raw string is not N")
	}
	if mt.PriorityUrgent.RawString() != "U" {
		t.Error("PriorityUrgent raw string is not U")
	}
}

func TestDeliveryMonitor(t *testing.T) {
	t.Parallel()

	if mt.DeliveryMonitorNonDelivery.RawString() != "1" {
		t.Error("DeliveryMonitorNonDelivery raw string is not 1")
	}
	if mt.DeliveryMonitorDelivery.RawString() != "2" {
		t.Error("DeliveryMonitorDelivery raw string is not 2")
	}
	if mt.DeliveryMonitorBoth.RawString() != "3" {
		t.Error("DeliveryMonitorBoth raw string is not 3")
	}
}

func TestCreditDebit(t *testing.T) {
	t.Parallel()

	if mt.Credit.RawString() != "C" {
		t.Error("Credit raw string is not C")
	}
	if mt.Debit.RawString() != "D" {
		t.Error("Debit raw string is not D")
	}
}

func TestBalance(t *testing.T) {
	if (mt.Balance{Raw: "123"}).RawString() != "123" {
		t.Error("Balance raw string is not 123")
	}

	for _, test := range []struct {
		name            string
		input           string
		expectedErr     error
		expectedBalance mt.Balance
	}{
		{
			name:        "InvalidInputLength",
			input:       "C0310",
			expectedErr: fmt.Errorf("balance: invalid input length: 5"),
		},
		{
			name:        "InvalidCreditDebit",
			input:       "E031002PLN40000,00",
			expectedErr: fmt.Errorf("balance: credit/debit: invalid indicator: E"),
		},
		{
			name:        "InvalidDate",
			input:       "CX31002PLN400X0,00",
			expectedErr: fmt.Errorf("balance: invalid date"),
		},
		{
			name:        "InvalidAmount",
			input:       "C031002PLN400X0,00",
			expectedErr: fmt.Errorf("balance: invalid amount"),
		},
		{
			name:  "ValidCredit",
			input: "C031002PLN40000,00",
			expectedBalance: mt.Balance{
				Set:         true,
				Raw:         "C031002PLN40000,00",
				CreditDebit: mt.Credit,
				Date: mt.Date{
					Set: true,
					Raw: "031002",
				},
				Currency: "PLN",
				Amount:   40000.00,
			},
		},
		{
			name:  "ValidDebit",
			input: "D031002PLN40000,00",
			expectedBalance: mt.Balance{
				Set:         true,
				Raw:         "D031002PLN40000,00",
				CreditDebit: mt.Debit,
				Date: mt.Date{
					Set: true,
					Raw: "031002",
				},
				Currency: "PLN",
				Amount:   40000.00,
			},
		},
	} {
		test := test

		t.Run("UnmarshalMT/"+test.name, func(t *testing.T) {
			t.Parallel()

			var balance mt.Balance
			err := balance.UnmarshalMT(test.input)
			mttest.ValidateError(t, test.expectedErr, err)
			mttest.ValidateBalance(t, "Result", test.expectedBalance, balance)
		})
	}
}

func TestFundsCode(t *testing.T) {
	t.Parallel()

	if mt.FundsCodeCredit.RawString() != "C" {
		t.Error("FundsCodeCredit raw string is not C")
	}
	if mt.FundsCodeCreditReversal.RawString() != "RC" {
		t.Error("FundsCodeCreditReversal raw string is not RC")
	}
	if mt.FundsCodeDebit.RawString() != "D" {
		t.Error("FundsCodeDebit raw string is not D")
	}
	if mt.FundsCodeDebitReversal.RawString() != "RD" {
		t.Error("FundsCodeDebitReversal raw string is not RD")
	}
}

func TestStatementLine(t *testing.T) {
	if (mt.StatementLine{Raw: "123"}).RawString() != "123" {
		t.Error("StatementLine raw string is not 123")
	}

	for _, test := range []struct {
		name                  string
		input                 string
		expectedErr           error
		expectedStatementLine mt.StatementLine
	}{
		{
			name:        "InvalidDate",
			input:       "0X10201020A20000,00FMSCNONREF//8327000090031789\nCard transaction",
			expectedErr: fmt.Errorf("statement line: invalid date"),
		},
		{
			name:        "InvalidEntryDate",
			input:       "0310201C20A20000,00FMSCNONREF//8327000090031789\nCard transaction",
			expectedErr: fmt.Errorf("statement line: invalid entry date"),
		},
		{
			name:        "InvalidFundsCode",
			input:       "0310201020A20000,00FMSCNONREF//8327000090031789\nCard transaction",
			expectedErr: fmt.Errorf("statement line: invalid or missing funds code"),
		},
		{
			name:  "ValidCredit",
			input: "0310201020C20000,00FMSCNONREF//8327000090031789\nCard transaction",
			expectedStatementLine: mt.StatementLine{
				Set: true,
				Raw: "0310201020C20000,00FMSCNONREF//8327000090031789\nCard transaction",
				Date: mt.Date{
					Set: true,
					Raw: "031020",
				},
				EntryDate: mt.Month{
					Set: true,
					Raw: "1020",
				},
				FundsCode:             mt.FundsCodeCredit,
				Amount:                20000.00,
				SwiftCode:             "FMSC",
				AccountOwnerReference: "NONREF",
				BankReference:         "//8327000090031789",
				Description:           "Card transaction",
			},
		},
		{
			name:  "ValidCreditReversal",
			input: "0310201020RC20000,00FMSCNONREF//8327000090031789\nCard transaction",
			expectedStatementLine: mt.StatementLine{
				Set: true,
				Raw: "0310201020RC20000,00FMSCNONREF//8327000090031789\nCard transaction",
				Date: mt.Date{
					Set: true,
					Raw: "031020",
				},
				EntryDate: mt.Month{
					Set: true,
					Raw: "1020",
				},
				FundsCode:             mt.FundsCodeCreditReversal,
				Amount:                20000.00,
				SwiftCode:             "FMSC",
				AccountOwnerReference: "NONREF",
				BankReference:         "//8327000090031789",
				Description:           "Card transaction",
			},
		},
		{
			name:  "ValidDebit",
			input: "0310201020D20000,00FMSCNONREF//8327000090031789\nCard transaction",
			expectedStatementLine: mt.StatementLine{
				Set: true,
				Raw: "0310201020D20000,00FMSCNONREF//8327000090031789\nCard transaction",
				Date: mt.Date{
					Set: true,
					Raw: "031020",
				},
				EntryDate: mt.Month{
					Set: true,
					Raw: "1020",
				},
				FundsCode:             mt.FundsCodeDebit,
				Amount:                20000.00,
				SwiftCode:             "FMSC",
				AccountOwnerReference: "NONREF",
				BankReference:         "//8327000090031789",
				Description:           "Card transaction",
			},
		},
		{
			name:  "ValidDebitReversal",
			input: "0310201020RD20000,00FMSCNONREF//8327000090031789\nCard transaction",
			expectedStatementLine: mt.StatementLine{
				Set: true,
				Raw: "0310201020RD20000,00FMSCNONREF//8327000090031789\nCard transaction",
				Date: mt.Date{
					Set: true,
					Raw: "031020",
				},
				EntryDate: mt.Month{
					Set: true,
					Raw: "1020",
				},
				FundsCode:             mt.FundsCodeDebitReversal,
				Amount:                20000.00,
				SwiftCode:             "FMSC",
				AccountOwnerReference: "NONREF",
				BankReference:         "//8327000090031789",
				Description:           "Card transaction",
			},
		},
	} {
		test := test

		t.Run("UnmarshalMT/"+test.name, func(t *testing.T) {
			t.Parallel()

			var statementLine mt.StatementLine
			err := statementLine.UnmarshalMT(test.input)
			mttest.ValidateError(t, test.expectedErr, err)
			mttest.ValidateStatementLine(t, test.expectedStatementLine, statementLine)
		})
	}
}

func TestBase(t *testing.T) {
	t.Parallel()

	input := mt.Base{AppHeaderInput: mt.AppHeaderInput{Set: true, MessagePriority: mt.PriorityNormal}}
	output := mt.Base{AppHeaderOutput: mt.AppHeaderOutput{Set: true, MessagePriority: mt.PriorityUrgent}}
	usr := mt.Base{UsrHeader: mt.UsrHeader{Set: true}}
	trl := mt.Base{Trailers: mt.Trailers{Set: true}}

	if !input.IsInput() {
		t.Error("expected input.IsInput to be true")
	}
	if input.IsOutput() {
		t.Error("expected input.IsInput to be false")
	}
	if input.Priority() != mt.PriorityNormal {
		t.Error("expected input.Priority() to be PriorityNormal")
	}
	if output.IsInput() {
		t.Error("expected output.IsInput to be false")
	}
	if !output.IsOutput() {
		t.Error("expected output.IsInput to be true")
	}
	if output.Priority() != mt.PriorityUrgent {
		t.Error("expected output.Priority() to be PriorityUrgent")
	}
	if input.HasUserHeader() {
		t.Error("expected input.HasUserHeader to be false")
	}
	if !usr.HasUserHeader() {
		t.Error("expected usr.HasUserHeader to be true")
	}
	if input.HasTrailers() {
		t.Error("expected input.HasTrailers to be false")
	}
	if !trl.HasTrailers() {
		t.Error("expected trl.HasTrailers to be true")
	}
}
