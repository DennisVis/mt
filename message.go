// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mt

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ApplicationID identifies the application within which the message is being sent or received. The available options
// are: F = FIN All user-to-user, FIN system and FIN service messages A = GPA (General Purpose Application) Most GPA
// system and service messages L = GPA Certain GPA service messages, for example, LOGIN, LAKs, ABORT These values are
// automatically assigned by the SWIFT system and the user's CBT.
type ApplicationID int

const (
	ApplicationIDFinancial ApplicationID = iota // F
	ApplicationIDGeneral                        // A
	ApplicationIDLogin                          // L
)

func (aid ApplicationID) String() string {
	switch aid {
	case ApplicationIDGeneral:
		return "A"
	case ApplicationIDLogin:
		return "L"
	// ApplicationIDFinancial
	default:
		return "F"
	}
}

func (aid ApplicationID) RawString() string {
	return aid.String()
}

// ServiceID consists of two numeric characters. It identifies the type of data that is being sent or received and, in
// doing so, whether the message which follows is one of the following: a user-to-user message, a system message, a
// service message, for example, a session control command, such as SELECT, or a logical acknowledgment, such as
// ACK/SAK/UAK. Possible values are 01 = FIN/GPA or 21 = ACK/NAK.
type ServiceID int

const (
	ServiceIDFINGPA  ServiceID = iota // 01
	ServiceIDACKNACK                  // 21
)

func (sid ServiceID) String() string {
	switch sid {
	case ServiceIDACKNACK:
		return "21"
	// ServiceIDFINGPA
	default:
		return "01"
	}
}

func (sid ServiceID) RawString() string {
	return sid.String()
}

// Priority is used within FIN Application Headers only, defines the priority with which a message is delivered.
// The possible values are:
// S = System
// U = Urgent
// N = Normal
type Priority int

const (
	PriorityNormal Priority = iota // N
	PrioritySystem                 // S
	PriorityUrgent                 // U
)

func (p Priority) String() string {
	switch p {
	case PrioritySystem:
		return "S"
	case PriorityUrgent:
		return "U"
	// PriorityNormal
	default:
		return "N"
	}
}

func (p Priority) RawString() string {
	return p.String()
}

// DeliveryMonitor applies only to FIN user-to-user messages. The chosen option is expressed as a single digit:
// 1 = Non-Delivery Warning
// 2 = Delivery Notification
// 3 = Non-Delivery Warning and Delivery Notification
// If the message has priority 'U', the user must request delivery monitoring option '1' or '3'. If the message has
// priority 'N', the user can request delivery monitoring option '2' or, by leaving the option blank, no delivery
// monitoring.
type DeliveryMonitor int

const (
	DeliveryMonitorNonDelivery DeliveryMonitor = iota // 1
	DeliveryMonitorDelivery                           // 2
	DeliveryMonitorBoth                               // 3
)

func (dm DeliveryMonitor) String() string {
	switch dm {
	case DeliveryMonitorNonDelivery:
		return "1"
	case DeliveryMonitorDelivery:
		return "2"
	// DeliveryMonitorBoth
	default:
		return "3"
	}
}

func (dm DeliveryMonitor) RawString() string {
	return dm.String()
}

// CreditDebit indicates whether the balance is a credit or a debit.
type CreditDebit int

const (
	Credit CreditDebit = iota
	Debit
)

func (cd CreditDebit) String() string {
	switch cd {

	case Debit:
		return "D"
	// Credit
	default:
		return "C"
	}
}

func (cd CreditDebit) RawString() string {
	return cd.String()
}

func creditDebitFromString(input string) (CreditDebit, error) {
	switch input {
	case "C":
		return Credit, nil
	case "D":
		return Debit, nil
	default:
		return 0, fmt.Errorf("credit/debit: invalid indicator: %s", input)
	}
}

// Balance represents the balance of a given account at a given date.
type Balance struct {
	Set         bool
	Raw         string
	CreditDebit CreditDebit `mt:"M,1!a"`
	Date        Date        `mt:"M,6!n"`
	Currency    string      `mt:"M,3!a"`
	Amount      float32     `mt:"M,15d"`
}

func (b *Balance) UnmarshalMT(input string) error {
	// example:
	// C031002PLN40000,00

	// min: all fixed length fields plus at least 1 for amount
	// max: all fixed length fields plus max 15 for amount
	if len(input) < 11 || len(input) > 25 {
		return fmt.Errorf("balance: invalid input length: %d", len(input))
	}

	// mandatory, 1!a
	creditDebitStr := input[0:1]
	creditDebit, err := creditDebitFromString(creditDebitStr)
	if err != nil {
		return fmt.Errorf("balance: %w", err)
	}
	b.CreditDebit = creditDebit

	// mandatory, 6!n
	dateStr := input[1:7]
	d := Date{}
	err = d.UnmarshalMT(dateStr)
	if err != nil {
		return fmt.Errorf("balance: invalid date")
	}
	b.Date = d

	// mandatory, 3!a
	b.Currency = input[7:10]

	// mandatory, 15d
	amountStr := input[10:]
	amount, err := strconv.ParseFloat(strings.ReplaceAll(amountStr, ",", "."), 32)
	if err != nil {
		return fmt.Errorf("balance: invalid amount")
	}
	b.Amount = float32(amount)

	b.Set = true
	b.Raw = input

	return nil
}

func (b Balance) RawString() string {
	return b.Raw
}

type FundsCode int

const (
	FundsCodeCredit         FundsCode = iota // C
	FundsCodeCreditReversal                  // RC
	FundsCodeDebit                           // D
	FundsCodeDebitReversal                   // RD
)

func (fc FundsCode) String() string {
	switch fc {
	case FundsCodeCreditReversal:
		return "RC"
	case FundsCodeDebit:
		return "D"
	case FundsCodeDebitReversal:
		return "RD"
	// FundsCodeCredit
	default:
		return "C"
	}
}

func (fc FundsCode) RawString() string {
	return fc.String()
}

type StatementLine struct {
	Set                   bool
	Raw                   string
	Date                  Date      `mt:"M,6!n"`
	EntryDate             Month     `mt:"O,4!n"`
	FundsCode             FundsCode `mt:"M,2a"`
	Amount                float64   `mt:"M,15d"`
	SwiftCode             string    `mt:"M,1!a3!c"`
	AccountOwnerReference string    `mt:"M,16x"`
	BankReference         string    `mt:"O,//20x"`
	Description           string    `mt:"O,34a"`
}

func (sl *StatementLine) UnmarshalMT(input string) error {
	// example:
	// 0310201020C20000,00FMSCNONREF//8327000090031789
	// Card transaction

	lines := strings.Split(input, "\n")

	line1 := lines[0]

	// mandatory, 6!n
	dateStr := line1[0:6]
	d := Date{}
	err := d.UnmarshalMT(dateStr)
	if err != nil {
		return fmt.Errorf("statement line: invalid date")
	}
	sl.Date = d
	line1 = line1[6:]

	hasEntryDate := !strings.HasPrefix(line1, "C") && !strings.HasPrefix(line1, "D")

	if hasEntryDate {
		// optional, 4!n
		entryDateStr := line1[0:4]
		month := Month{}
		err := month.UnmarshalMT(entryDateStr)
		if err != nil {
			return fmt.Errorf("statement line: invalid entry date")
		}
		sl.EntryDate = month
		line1 = line1[4:]
	}

	// mandatory, 2a
	switch {
	case strings.HasPrefix(line1, "RC"):
		sl.FundsCode = FundsCodeCreditReversal
		line1 = line1[2:]
	case strings.HasPrefix(line1, "RD"):
		sl.FundsCode = FundsCodeDebitReversal
		line1 = line1[2:]
	case strings.HasPrefix(line1, "C"):
		sl.FundsCode = FundsCodeCredit
		line1 = line1[1:]
	case strings.HasPrefix(line1, "D"):
		sl.FundsCode = FundsCodeDebit
		line1 = line1[1:]
	default:
		return fmt.Errorf("statement line: invalid or missing funds code")
	}

	amountNrOfDigits := 0
	for unicode.IsDigit(rune(line1[amountNrOfDigits])) || line1[amountNrOfDigits] == ',' {
		amountNrOfDigits++
	}

	// mandatory, 15d
	amountStr := line1[0:amountNrOfDigits]
	// above we've made sure to only regard digits and commas
	// therefore we cane safely ignore the error
	//nolint
	amount, _ := strconv.ParseFloat(strings.ReplaceAll(amountStr, ",", "."), 32)
	sl.Amount = amount
	line1 = line1[amountNrOfDigits:]

	// mandatory, 1!a3!c
	sl.SwiftCode = line1[0:4]
	line1 = line1[4:]

	split := strings.Split(line1, "//")
	if len(split) == 2 {
		sl.AccountOwnerReference = split[0]
		sl.BankReference = "//" + split[1]
	} else {
		sl.AccountOwnerReference = line1
	}

	if len(lines) > 1 {
		sl.Description = lines[1]
	}

	sl.Set = true
	sl.Raw = input

	return nil
}

func (sl StatementLine) RawString() string {
	return sl.Raw
}

// OutputReference is a reference to an output message containing both the send date and time of said message.
type OutputReference struct {
	Set                    bool
	Raw                    string
	LogicalTerminalAddress string
	SessionNumber          string
	SequenceNumber         string
	DateOrDateTime         DateOrDateTime
}

// InputReference is a reference to an input message containing only the send date of said message.
type InputReference struct {
	Set                    bool
	Raw                    string
	LogicalTerminalAddress string
	SessionNumber          string
	SequenceNumber         string
	DateOrDateTime         DateOrDateTime
}

// Reference is a reference to an original user message.
type Reference struct {
	Set                   bool
	Raw                   string
	DateTime              DateTime
	MessageInputReference InputReference
}

// BasicHeader is the only mandatory block; block 1. The basic header contains the general information that identifies
// the message, and some additional control information. The FIN interface automatically builds the basic header.
type BasicHeader struct {
	Raw                    string
	AppID                  ApplicationID
	ServiceID              ServiceID
	SessionNumber          string
	SequenceNumber         string
	LogicalTerminalAddress string
}

// AppHeaderInput contains information, from block 2, that is specific to the application. The application
// header is required for messages that users, or the system and users, exchange. Exceptions are session establishment
// and session closure.
//
// It is filled if and only if the message is of the input variety and therefore might not have been set. It is advised
// to verify whether it has been set before accessing the data within. This can be done with the Set field or the
// IsInput function on the containing message.
type AppHeaderInput struct {
	Set                         bool
	Raw                         string
	ObsolescencePeriodInMinutes int
	MessageType                 string
	ReceiverAddress             string
	MessagePriority             Priority
	DeliveryMonitor             DeliveryMonitor
}

// AppHeaderOutput contains information, from block 2, that is specific to the application. The application header is
// required for messages that users, or the system and users, exchange. Exceptions are session establishment and session
// closure.
//
// It is filled if and only if the message is of the output variety and herefore might not have been set. It is advised
// to verify whether it has been set before accessing the data within. This can be done with the Set field or the
// IsOutput function on the containing message.
type AppHeaderOutput struct {
	Set                   bool
	Raw                   string
	MessagePriority       Priority
	MessageType           string
	MessageInputReference InputReference
	InputTime             Time
	OutputDate            Date
	OutputTime            Time
}

// UsrHeader is an optional header that contains the information from block 3.
//
// This header is optional and therefore might not have been set. It is advised to verify whether it has been set
// before accessing the data within. This can be done with the Set field or the HasUserHeader function on the message.
type UsrHeader struct {
	Set                                bool
	Raw                                string
	ServiceID                          string
	AddresseeInformation               string
	BankingPriority                    string
	MessageUserReference               string
	ValidationFlag                     string
	RelatedReference                   string
	ServiceTypeID                      string
	UniqueEndToEndTransactionReference string
	PaymentReleaseInformation          string
	SanctionsScreeningInformation      string
	PaymentControlsInformation         string
	BalanceCheckpointDateTime          DateTimeSecOptCent
	MessageInputReference              InputReference
}

// PossibleDuplicateEmission is added if user thinks the same message was sent previously.
type PossibleDuplicateEmission struct {
	Raw                   string
	Time                  Time
	MessageInputReference InputReference
}

// PossibleDuplicateMessage is added by the system to any output message (GPA and FIN with a Service Identifier of 01)
// being resent because a prior delivery may not be valid. If a system PLT receives a report request with a PDM trailer,
// the response has a plain PDM (without the optional delivery reference). Other PDMs may be added because of
// unsuccessful delivery attempts to the user.
type PossibleDuplicateMessage struct {
	Raw                    string
	Time                   Time
	MessageOutputReference OutputReference
}

// SystemOriginatedMessage is the system message or service message.
type SystemOriginatedMessage struct {
	Raw                   string
	Time                  Time
	MessageInputReference InputReference
}

// Trailers contains the information from block 5. The trailer either indicates special circumstances that relate
// to message handling or contains security information.
//
// Trailers are optional and therefore might not have been set. It is advised to verify whether they have been set
// before accessing the data within. This can be done with the Set field or the HasTrailers function on the message.
type Trailers struct {
	Set                       bool
	Raw                       string
	DelayedMessage            bool
	TestAndTrainingMessage    bool
	Checksum                  string
	MessageReference          Reference
	PossibleDuplicateEmission PossibleDuplicateEmission
	PossibleDuplicateMessage  PossibleDuplicateMessage
	SystemOriginatedMessage   SystemOriginatedMessage
	AdditionalTrailers        map[string]string
}

// Base holds the basic structure all MT messages adhere to, excluding the body.
type Base struct {
	Raw             string
	Line            int
	BasicHeader     BasicHeader
	AppHeaderInput  AppHeaderInput
	AppHeaderOutput AppHeaderOutput
	UsrHeader       UsrHeader
	Trailers        Trailers
}

// IsInput returns true if the message is of the input variety. If so it will contain an input type app header.
// It is advised to use this function before accessing information in the AppHeaderInput struct.
func (b Base) IsInput() bool {
	return b.AppHeaderInput.Set
}

// IsOutput returns true if the message is of the output variety. If so it will contain an output type app header.
// It is advised to use this function before accessing information in the AppHeaderOutput struct.
func (b Base) IsOutput() bool {
	return b.AppHeaderOutput.Set
}

// Priority takes the message type from the app header, taking into account whether the message is input or output.
func (b Base) Type() string {
	if b.IsInput() {
		return b.AppHeaderInput.MessageType
	}
	return b.AppHeaderOutput.MessageType
}

// Priority takes the priority from the app header, taking into account whether the message is input or output.
func (b Base) Priority() Priority {
	if b.IsInput() {
		return b.AppHeaderInput.MessagePriority
	}
	return b.AppHeaderOutput.MessagePriority
}

// HasUserHeader returns true if the user header of the message was set/filled.
// It is advised to use this function before accessing information in the UsrHeader struct.
func (b Base) HasUserHeader() bool {
	return b.UsrHeader.Set
}

// HasTrailers returns true if the trailers of the message were set/filled.
// It is advised to use this function before accessing information in the Trailers struct.
func (b Base) HasTrailers() bool {
	return b.Trailers.Set
}

// MTx represents a complete message including headers and a body. The body has not been further processes or validated.
// The specific type of MT message this holds can be determined by the Type() function.
//
// If parsing into a more specifc struct is desireed the MTxToMT... functions be used.
// After parsing into a more specific type the ValidateMT... functions can be used to validate the message.
type MTx struct {
	Base
	Body map[string][]string
}
