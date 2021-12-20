// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mt

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/DennisVis/mt/internal/message"
)

const obsolescenceMinutesPerFactor = 5

var leadingZerosRegexp = regexp.MustCompile(`^0+`)

// basicHeaderBlockToBasicHeader parses the basic header block and returns a BasicHeader struct.
//
// The header block content should be in the following format:
//
// F01SCBLZAJJXXXX5712100002
//
// Which will be parsed as:
//
// F				<- application ID
// 01				<- service ID
// SCBLZAJJXXXX		<- logical terminal address
// 5712				<- session number
// 100002			<- sequence number
func basicHeaderBlockToBasicHeader(block message.Block) (BasicHeader, error) {
	msgBscHeader := BasicHeader{
		Raw: "{1:" + block.Content + "}",
	}

	if len(block.Content) != 25 {
		return msgBscHeader, fmt.Errorf("invalid basic header block content length: %d", len(block.Content))
	}

	switch block.Content[0:1] {
	case "F":
		msgBscHeader.AppID = ApplicationIDFinancial
	case "A":
		msgBscHeader.AppID = ApplicationIDGeneral
	case "L":
		msgBscHeader.AppID = ApplicationIDLogin
	default:
		return msgBscHeader, fmt.Errorf("unknown application id in basic header block content: %s", block.Content[0:1])
	}

	switch block.Content[1:3] {
	case "01":
		msgBscHeader.ServiceID = ServiceIDFINGPA
	case "21":
		msgBscHeader.ServiceID = ServiceIDACKNACK
	default:
		return msgBscHeader, fmt.Errorf("unknown service id in basic header block content: %s", block.Content[1:3])
	}

	msgBscHeader.LogicalTerminalAddress = block.Content[3:15]
	msgBscHeader.SessionNumber = block.Content[15:19]
	msgBscHeader.SequenceNumber = block.Content[19:]

	return msgBscHeader, nil
}

// 120811BANKFRPPAXXX2222123456
func stringToMessageInputReferenceDate(str string) (InputReference, error) {
	mird := InputReference{
		Set: true,
		Raw: str,
	}

	if len(str) != 28 {
		return mird, fmt.Errorf("invalid message input reference with date string length: %d", len(str))
	}

	dateStr := str[0:6]
	var date DateOrDateTime
	err := date.UnmarshalMT(dateStr)
	if err != nil {
		return mird, fmt.Errorf("invalid message input reference with date date string: %s: %w", dateStr, err)
	}
	mird.DateOrDateTime = date

	mird.LogicalTerminalAddress = str[6:18]
	mird.SessionNumber = str[18:22]
	mird.SequenceNumber = str[22:]

	return mird, nil
}

// 1806271539180626BANKFRPPAXXX2222123456
func stringToMessageReference(str string) (Reference, error) {
	mr := Reference{
		Set: true,
		Raw: str,
	}

	dateTimeStr := str[0:10]
	var dateTime DateTime
	err := dateTime.UnmarshalMT(dateTimeStr)
	if err != nil {
		return mr, fmt.Errorf("invalid message reference date/time string: %s: %w", dateTimeStr, err)
	}
	mr.DateTime = dateTime

	mir, err := stringToMessageInputReferenceDate(str[10:])
	if err != nil {
		return mr, fmt.Errorf("invalid message reference message input reference: %s: %w", str[10:], err)
	}
	mr.MessageInputReference = mir

	return mr, nil
}

// (1348)120811BANKFRPPAXXX2222123456
func stringToMessageOutputReference(str string) (OutputReference, error) {
	mor := OutputReference{
		Set: true,
		Raw: str,
	}

	switch len(str) {
	case 28:
		dateTimeStr := str[0:6]
		var dateTime DateOrDateTime
		err := dateTime.UnmarshalMT(dateTimeStr)
		if err != nil {
			return mor, fmt.Errorf("invalid message output reference date/time string: %s: %w", dateTime, err)
		}
		mor.DateOrDateTime = dateTime

		mor.LogicalTerminalAddress = str[6:18]
		mor.SessionNumber = str[18:23]
		mor.SequenceNumber = str[23:]
	case 32:
		dateTimeStr := str[4:10] + str[0:4]
		var dateTime DateOrDateTime
		err := dateTime.UnmarshalMT(dateTimeStr)
		if err != nil {
			return mor, fmt.Errorf("invalid message output reference date/time string: %s: %w", dateTime, err)
		}
		mor.DateOrDateTime = dateTime

		mor.LogicalTerminalAddress = str[10:22]
		mor.SessionNumber = str[22:27]
		mor.SequenceNumber = str[27:]
	}

	return mor, nil
}

// 1348120811BANKFRPPAXXX2222123456
func stringToPossibleDuplicateEmission(str string) (PossibleDuplicateEmission, error) {
	pde := PossibleDuplicateEmission{
		Raw: str,
	}

	if len(str) != 32 {
		return pde, fmt.Errorf("invalid possible duplicate emission string length: %d", len(str))
	}

	timeStr := str[0:4]
	var time Time
	err := time.UnmarshalMT(timeStr)
	if err != nil {
		return pde, fmt.Errorf("invalid possible duplicate emission time: %s: %w", timeStr, err)
	}
	pde.Time = time

	mir, err := stringToMessageInputReferenceDate(str[4:])
	if err != nil {
		return pde, fmt.Errorf("invalid possible duplicate emission message input reference: %s: %w", str[4:], err)
	}
	pde.MessageInputReference = mir

	return pde, nil
}

// 1213120811BANKFRPPAXXX2222123456
func stringToPossibleDuplicateMessage(str string) (PossibleDuplicateMessage, error) {
	pdm := PossibleDuplicateMessage{
		Raw: str,
	}

	if len(str) != 32 && len(str) != 36 {
		return pdm, fmt.Errorf("invalid possible duplicate message string length: %d", len(str))
	}

	timeStr := str[0:4]
	var time Time
	err := time.UnmarshalMT(timeStr)
	if err != nil {
		return pdm, fmt.Errorf("invalid possible duplicate message time: %s: %w", timeStr, err)
	}
	pdm.Time = time

	mor, err := stringToMessageOutputReference(str[4:])
	if err != nil {
		return pdm, fmt.Errorf("invalid possible duplicate message message output reference: %s: %w", str[4:], err)
	}
	pdm.MessageOutputReference = mor

	return pdm, nil
}

// 1454120811BANKFRPPAXXX2222123456
func stringToSystemOriginatedMessage(str string) (SystemOriginatedMessage, error) {
	som := SystemOriginatedMessage{
		Raw: str,
	}

	if len(str) != 32 {
		return som, fmt.Errorf("invalid system originated message string length: %d", len(str))
	}

	timeStr := str[0:4]
	var time Time
	err := time.UnmarshalMT(timeStr)
	if err != nil {
		return som, fmt.Errorf("invalid system originated message time: %s: %w", timeStr, err)
	}
	som.Time = time

	mir, err := stringToMessageInputReferenceDate(str[4:])
	if err != nil {
		return som, fmt.Errorf("invalid system originated message message input reference: %s: %w", str[4:], err)
	}
	som.MessageInputReference = mir

	return som, nil
}

// appHeaderBlockToAppHeaderInput parses the app header block as a AppHeaderInput struct.
//
// The app header input block content should be in the following format:
//
// I940BOFAUS6BXBAMN
// I940SCBLZAJJXXXXN2020
//
// Which will be parsed as:
//
// I			<- Input/Output type, will be ignored from here on
// 940			<- Message type
// SCBLZAJJXXXX	<- Receiver's address
// N			<- Message priority (optional)
// 2			<- Delivery monitor (optional)
// 020			<- Obsolescence period in magnitudes of 5 minutes (003 - 15 minutes, 020 - 100 minutes) (optional)
func appHeaderBlockToAppHeaderInput(block message.Block) (AppHeaderInput, error) {
	msgAppHeaderIn := AppHeaderInput{
		Raw: "{2:" + block.Content + "}",
	}

	if len(block.Content) < 16 {
		return msgAppHeaderIn, fmt.Errorf("invalid app header input block content length: %d", len(block.Content))
	}

	msgAppHeaderIn.Set = true
	msgAppHeaderIn.MessageType = block.Content[1:4] // from 1 as we don't care about the I anymore, it's dropped
	msgAppHeaderIn.ReceiverAddress = block.Content[4:16]

	setPriority := func(char string) error {
		switch char {
		case "S":
			msgAppHeaderIn.MessagePriority = PrioritySystem
			return nil
		case "N":
			msgAppHeaderIn.MessagePriority = PriorityNormal
			return nil
		case "U":
			msgAppHeaderIn.MessagePriority = PriorityUrgent
			return nil
		default:
			return fmt.Errorf("unknown message priority in app header input block content: %s", char)
		}
	}

	setDeliveryMonitor := func(char string) error {
		switch char {
		case "1":
			msgAppHeaderIn.DeliveryMonitor = DeliveryMonitorNonDelivery
			return nil
		case "2":
			msgAppHeaderIn.DeliveryMonitor = DeliveryMonitorDelivery
			return nil
		case "3":
			msgAppHeaderIn.DeliveryMonitor = DeliveryMonitorBoth
			return nil
		default:
			return fmt.Errorf("invalid delivery monitor in app header input block content: %s", char)
		}
	}

	setPriorityOrDeliveryMonitor := func(char string) error {
		switch char {
		case "S", "N", "U":
			return setPriority(char)
		case "1", "2", "3":
			return setDeliveryMonitor(char)
		default:
			return fmt.Errorf("invalid priority or delivery monitor in app header input block content: %s", char)
		}
	}

	setObsolescencePeriod := func(chars string) error {
		factorString := string(leadingZerosRegexp.ReplaceAll([]byte(chars), []byte("")))
		factor, err := strconv.Atoi(factorString)
		if err != nil {
			return fmt.Errorf("invalid obsolescence period in app header input block content: %v: %w", factor, err)
		}

		msgAppHeaderIn.ObsolescencePeriodInMinutes = factor * obsolescenceMinutesPerFactor

		return nil
	}

	switch len(block.Content) {
	// optional fields not present, nothing left to do
	// I940SCBLZAJJXXXX
	case 16:
		break
	// of the optional fields only priority or delivery monitor present
	// I940SCBLZAJJXXXXN
	// I940SCBLZAJJXXXX2
	case 17:
		err := setPriorityOrDeliveryMonitor(block.Content[16:17])
		if err != nil {
			return msgAppHeaderIn, fmt.Errorf(
				"could not set priority or delivery monitor for app header input: %w",
				err,
			)
		}
	// of the optional fields only priority and delivery monitor present
	// I940SCBLZAJJXXXXN2
	case 18:
		err := setPriority(block.Content[16:17])
		if err != nil {
			return msgAppHeaderIn, fmt.Errorf("could not set priority for app header input: %w", err)
		}
		err = setDeliveryMonitor(block.Content[17:18])
		if err != nil {
			return msgAppHeaderIn, fmt.Errorf("could not set delivery monitor for app header input: %w", err)
		}
	// of the optional fields only obsolescence period present
	// I940SCBLZAJJXXXX020
	case 19:
		err := setObsolescencePeriod(block.Content[16:])
		if err != nil {
			return msgAppHeaderIn, fmt.Errorf("could not set obsolescence period for app header input: %w", err)
		}
	// of the optional fields priority or delivery monitor and obsolescence period present
	// I940SCBLZAJJXXXXN020
	// I940SCBLZAJJXXXX2020
	case 20:
		err := setPriorityOrDeliveryMonitor(block.Content[16:17])
		if err != nil {
			return msgAppHeaderIn, fmt.Errorf(
				"could not set priority or delivery monitor for app header input: %w",
				err,
			)
		}
		err = setObsolescencePeriod(block.Content[17:])
		if err != nil {
			return msgAppHeaderIn, fmt.Errorf("could not set obsolescence period for app header input: %w", err)
		}
	// optional fields priority, delivery monitor and obsolescence period all present
	// I940SCBLZAJJXXXXN2020
	case 21:
		err := setPriority(block.Content[16:17])
		if err != nil {
			return msgAppHeaderIn, fmt.Errorf("could not set priority for app header input: %w", err)
		}
		err = setDeliveryMonitor(block.Content[17:18])
		if err != nil {
			return msgAppHeaderIn, fmt.Errorf("could not set delivery monitor for app header input: %w", err)
		}
		err = setObsolescencePeriod(block.Content[18:])
		if err != nil {
			return msgAppHeaderIn, fmt.Errorf("could not set obsolescence period for app header input: %w", err)
		}
	default:
		return msgAppHeaderIn, fmt.Errorf("invalid app header input block content length: %d", len(block.Content))
	}

	return msgAppHeaderIn, nil
}

// appHeaderBlockToAppHeaderOutput parses the app header block as AppHeaderOutput struct.
//
// The app header input block content should be in the following format:
//
// O9401157091028SCBLZAJJXXXX57121000020910281157N
//
// Which will be parsed as:
//
// O							<- Input/Output type, will be ignored from here on
// 940							<- Message type
// 1157							<- Input time with respect to the sender's time zone
// 091028SCBLZAJJXXXX5712100002	<- Message Input Reference (MIR) which is parsed as:
// - 091028						<-- Sender's date
// - SCBLZAJJXXXX				<-- Logical terminal address
// - 5712						<-- Session number
// - 100002						<-- Sequence number
// 091028						<- Output date
// 1157							<- Output time
// N							<- Message priority (optional)
func appHeaderBlockToAppHeaderOutput(block message.Block) (AppHeaderOutput, error) {
	msgAppHeaderOut := AppHeaderOutput{
		Raw: "{2:" + block.Content + "}",
	}

	if len(block.Content) < 46 {
		return msgAppHeaderOut, fmt.Errorf("invalid app header output block content length: %d", len(block.Content))
	}

	msgAppHeaderOut.Set = true
	msgAppHeaderOut.MessageType = block.Content[1:4] // from 1 as we don't care about the O anymore, it's dropped

	inputTimeStr := block.Content[4:8]
	var inputTime Time
	err := inputTime.UnmarshalMT(inputTimeStr)
	if err != nil {
		return msgAppHeaderOut, fmt.Errorf(
			"invalid input time in app header output block content: %v: %w",
			inputTimeStr,
			err,
		)
	}
	msgAppHeaderOut.InputTime = inputTime

	outputDateStr := block.Content[36:42]
	var outputDate Date
	err = outputDate.UnmarshalMT(outputDateStr)
	if err != nil {
		return msgAppHeaderOut, fmt.Errorf(
			"invalid output date in app header output block content: %v: %w",
			outputDateStr,
			err,
		)
	}
	msgAppHeaderOut.OutputDate = outputDate

	outputTimeStr := block.Content[42:46]
	var outputTime Time
	err = outputTime.UnmarshalMT(outputTimeStr)
	if err != nil {
		return msgAppHeaderOut, fmt.Errorf(
			"invalid output time in app header output block content: %v: %w",
			outputTimeStr,
			err,
		)
	}
	msgAppHeaderOut.OutputTime = outputTime

	mird, err := stringToMessageInputReferenceDate(block.Content[8:36])
	if err != nil {
		return msgAppHeaderOut, fmt.Errorf("could not parse message input reference with date: %w", err)
	}
	msgAppHeaderOut.MessageInputReference = mird

	if len(block.Content) == 47 {
		switch block.Content[46] {
		case 'N':
			msgAppHeaderOut.MessagePriority = PriorityNormal
		case 'S':
			msgAppHeaderOut.MessagePriority = PrioritySystem
		case 'U':
			msgAppHeaderOut.MessagePriority = PriorityUrgent
		default:
			return msgAppHeaderOut, fmt.Errorf("invalid message priority")
		}
	}

	return msgAppHeaderOut, nil
}

// appHeaderBlockToAppHeader decides if the given app header block is an input or output app header block and then
// passes parsing on to either appHeaderBlockToAppHeaderInput or appHeaderBlockToAppHeaderOutput respectivally.
func appHeaderBlockToAppHeader(block message.Block) (AppHeaderInput, AppHeaderOutput, error) {
	var appHeaderIn AppHeaderInput
	var appHeaderOut AppHeaderOutput
	var errToReturn error

	if len(block.Content) < 4 {
		return appHeaderIn, appHeaderOut, fmt.Errorf(
			"invalid app header block content length: %d",
			len(block.Content),
		)
	}

	switch block.Content[0:1] {
	case "I":
		msgAppHeaderIn, err := appHeaderBlockToAppHeaderInput(block)
		if err != nil {
			errToReturn = fmt.Errorf(
				"could not parse app header block as app header input: %w",
				err,
			)
		}
		appHeaderIn = msgAppHeaderIn
	case "O":
		msgAppHeaderOut, err := appHeaderBlockToAppHeaderOutput(block)
		if err != nil {
			errToReturn = fmt.Errorf(
				"could not parse app header block as app header output: %w",
				err,
			)
		}
		appHeaderOut = msgAppHeaderOut
	default:
		errToReturn = fmt.Errorf("invalid app header message type: %s", block.Content[0:1])
	}

	return appHeaderIn, appHeaderOut, errToReturn
}

// usrHeaderBlockToUsrHeader parses the user header block and returns a UsrHeader struct.
//
// The header block should contain one or more sub blocks. Each block will be processed, its label will decide which
// member of the struct its content will populate.
//
// To see which block corresponds to which struct member see the switch statement.
func usrHeaderBlockToUsrHeader(block message.Block) (UsrHeader, []error) {
	msgUsrHeader := UsrHeader{
		Set: true,
		Raw: "{3:" + block.Content + "}",
	}
	errors := make([]error, 0)

	for _, sb := range block.Blocks {
		switch sb.Label {
		case "103":
			msgUsrHeader.ServiceID = sb.Content
		case "106":
			msgInReference, err := stringToMessageInputReferenceDate(sb.Content)
			if err != nil {
				errors = append(errors, fmt.Errorf("invalid message input reference: %w", err))
				continue
			}

			msgUsrHeader.MessageInputReference = msgInReference
		case "108":
			msgUsrHeader.MessageUserReference = sb.Content
		case "111":
			msgUsrHeader.ServiceTypeID = sb.Content
		case "113":
			msgUsrHeader.BankingPriority = sb.Content
		case "115":
			msgUsrHeader.AddresseeInformation = sb.Content
		case "119":
			msgUsrHeader.ValidationFlag = sb.Content
		case "121":
			msgUsrHeader.UniqueEndToEndTransactionReference = sb.Content
		case "165":
			msgUsrHeader.PaymentReleaseInformation = sb.Content
		case "423":
			var balanceCheckpointDateTime DateTimeSecOptCent
			err := balanceCheckpointDateTime.UnmarshalMT(sb.Content)
			if err != nil {
				errors = append(errors, fmt.Errorf(
					"invalid balance checkpoint time in usr header block content: %s: %w",
					sb.Content,
					err,
				))
				continue
			}

			msgUsrHeader.BalanceCheckpointDateTime = balanceCheckpointDateTime
		case "424":
			msgUsrHeader.RelatedReference = sb.Content
		case "433":
			msgUsrHeader.SanctionsScreeningInformation = sb.Content
		case "434":
			msgUsrHeader.PaymentControlsInformation = sb.Content
		default:
			errors = append(errors, fmt.Errorf("invalid usr header block sub block label: %s", sb.Label))
		}
	}

	if len(errors) > 0 {
		return msgUsrHeader, errors
	}

	return msgUsrHeader, nil
}

// trailersBlockToTrailers parses the trailers block and returns a MessageTrailers struct.
//
// The trailers block should contain one or more sub blocks. Each block will be processed, its label will decide which
// member of the struct its content will populate.
//
// To see which block corresponds to which struct member see the switch statement.
func trailersBlockToTrailers(block message.Block) (Trailers, []error) {
	msgTrailers := Trailers{
		Set:                true,
		AdditionalTrailers: make(map[string]string),
	}
	errors := make([]error, 0)
	raw := "{5:"

	for _, sb := range block.Blocks {
		raw += "{" + sb.Label + ":" + sb.Content + "}"

		switch sb.Label {
		case "CHK":
			msgTrailers.Checksum = sb.Content
		case "TNG":
			msgTrailers.TestAndTrainingMessage = true
		case "PDE":
			pde, err := stringToPossibleDuplicateEmission(sb.Content)
			if err != nil {
				errors = append(errors, fmt.Errorf("invalid possible duplicate emission: %w", err))
			}
			msgTrailers.PossibleDuplicateEmission = pde
		case "DLM":
			msgTrailers.DelayedMessage = true
		case "MRF":
			mr, err := stringToMessageReference(sb.Content)
			if err != nil {
				errors = append(errors, fmt.Errorf("invalid message reference: %w", err))
			}
			msgTrailers.MessageReference = mr
		case "PDM":
			mor, err := stringToPossibleDuplicateMessage(sb.Content)
			if err != nil {
				errors = append(errors, fmt.Errorf("invalid possible duplicate message: %w", err))
			}
			msgTrailers.PossibleDuplicateMessage = mor
		case "SYS":
			som, err := stringToSystemOriginatedMessage(sb.Content)
			if err != nil {
				errors = append(errors, fmt.Errorf("invalid system originated message: %w", err))
			}
			msgTrailers.SystemOriginatedMessage = som
		default:
			msgTrailers.AdditionalTrailers[sb.Label] = sb.Content
		}
	}

	msgTrailers.Raw = raw + "}"

	if len(errors) > 0 {
		return msgTrailers, errors
	}

	return msgTrailers, nil
}

func messageToMTx(msg message.Message) (MTx, Errors) {
	mtx := MTx{}

	mtx.Raw = msg.Raw
	mtx.Body = msg.Body
	mtx.Line = msg.Line

	errors := make(Errors, 0)

	msgHeader, err := basicHeaderBlockToBasicHeader(msg.BasicHeader)
	if err != nil {
		errors = append(errors, NewError(fmt.Errorf("invalid basic header: %w", err), msg.Line))
	}
	mtx.BasicHeader = msgHeader

	appHeaderInput, appHeaderOutput, err := appHeaderBlockToAppHeader(msg.AppHeader)
	if err != nil {
		errors = append(errors, NewError(fmt.Errorf("invalid app header: %w", err), msg.Line))
	}
	mtx.AppHeaderInput = appHeaderInput
	mtx.AppHeaderOutput = appHeaderOutput

	usrHeader, errs := usrHeaderBlockToUsrHeader(msg.UsrHeader)
	for _, err := range errs {
		errors = append(errors, NewError(fmt.Errorf("invalid user header: %w", err), msg.Line))
	}
	mtx.UsrHeader = usrHeader

	trailers, errs := trailersBlockToTrailers(msg.Trailers)
	for _, err := range errs {
		errors = append(errors, NewError(fmt.Errorf("invalid trailers: %w", err), msg.Line))
	}
	mtx.Trailers = trailers

	if len(errors) > 0 {
		return mtx, errors
	}

	return mtx, nil
}
