// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package mt

import (
	"context"
	"fmt"
	"io"

	"github.com/DennisVis/mt/internal/encoding/mt"
	"github.com/DennisVis/mt/internal/validate"
)

const MessageTypeMT940 = "940"

var mt940Validator = validate.MustCreateValidatorForStruct(MT940{})

func MTxToMT940(mtx MTx) (MT940, error) {
	mt940 := MT940{}

	if mtx.Type() != MessageTypeMT940 {
		return mt940, fmt.Errorf("expected message type %s, got %s", MessageTypeMT940, mtx.Type())
	}

	mt940.Base = mtx.Base

	err := mt.UnmarshalMT(mtx.Body, &mt940)
	if err != nil {
		return mt940, fmt.Errorf("could not unmarshal MT%s message: %w", MessageTypeMT940, err)
	}

	err = mt940Validator.Validate(mt940)
	if err != nil {
		return mt940, fmt.Errorf("validation failed for MT%s message:\n%s", MessageTypeMT940, err)
	}

	return mt940, nil
}

func ValidateMT940(mt940 MT940) error {
	err := mt940Validator.Validate(mt940)
	if err != nil {
		return fmt.Errorf("validation failed for MT%s message:\n%w", MessageTypeMT940, err)
	}

	return nil
}

func parseAndValidateMT940(mtx MTx, skipValidation, lax bool) (MT940, error) {
	mt940, err := MTxToMT940(mtx)
	if err != nil || skipValidation {
		return mt940, err
	}

	err = ValidateMT940(mt940)
	if err != nil && !lax {
		return mt940, err
	}

	return mt940, nil
}

// ParseMT940 parses and validates MTx messages from ParseMTx into MT940 messages.
// Invalid messages are discarded unless the option Lax is passed.
func ParseMT940(ctx context.Context, rd io.Reader, options ...option) (chan MT940, chan Error) {
	cfg := optionsToConfig(options)

	genericMessages, parseErrors := ParseMTx(ctx, rd, options...)

	mt940Ch := make(chan MT940)

	go func() {
		for mtx := range genericMessages {
			mt940, err := parseAndValidateMT940(mtx, cfg.SkipValidation, cfg.Lax)
			if err != nil {
				parseErrors <- NewError(err, mtx.Line)

				if !cfg.Lax {
					continue
				}
			}

			mt940Ch <- mt940
		}
	}()

	return mt940Ch, parseErrors
}

// ParseAllMT940 parses and validates MTx messages from ParseAllMTx into MT940 messages.
// Invalid messages are discarded unless the option Lax is passed.
func ParseAllMT940(ctx context.Context, rd io.Reader, options ...option) ([]MT940, error) {
	cfg := optionsToConfig(options)

	genericMessages, pes := ParseAllMTx(ctx, rd, options...)

	mt940s := make([]MT940, 0)

	var parseErrors Errors
	if pes != nil {
		parseErrors = pes.(Errors)
	}

	for _, mtx := range genericMessages {
		mt940, err := parseAndValidateMT940(mtx, cfg.SkipValidation, cfg.Lax)
		if err != nil {
			parseErrors = append(parseErrors, NewError(err, mtx.Line))

			if !cfg.Lax {
				continue
			}
		}

		mt940s = append(mt940s, mt940)
	}

	return mt940s, parseErrors
}
