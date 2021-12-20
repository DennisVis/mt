// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pattern

import (
	"fmt"
	"strings"
)

type CharSet func(r rune) bool

type runeSet []rune

func (rs runeSet) contains(r rune) bool {
	for _, rr := range rs {
		if rr == r {
			return true
		}
	}
	return false
}

func charsetsKeysAsRunes(charSets map[string]CharSet) []rune {
	var keys []rune
	for k := range charSets {
		// we control the map below so know this is safe
		//nolint
		r, _, _ := strings.NewReader(k).ReadRune()
		keys = append(keys, r)
	}
	return keys
}

var (
	numbers           CharSet = func(r rune) bool { return r >= '0' && r <= '9' }
	alphaLower        CharSet = func(r rune) bool { return r >= 'a' && r <= 'z' }
	alphaUpper        CharSet = func(r rune) bool { return r >= 'A' && r <= 'Z' }
	alphaNumericUpper CharSet = func(r rune) bool { return numbers(r) || alphaUpper(r) }
	floats            CharSet = func(r rune) bool { return numbers(r) || r == ',' }
	special           CharSet = func(r rune) bool {
		return r == '/' || r == '-' || r == '?' || r == ':' || r == '(' || r == ')' || r == '.' || r == ',' || r == '\'' || r == '+' || r == '{' || r == '}' || r == '\n' || r == ' '
	}
	any      CharSet = func(r rune) bool { return alphaNumericUpper(r) || alphaLower(r) || floats(r) || special(r) }
	charSets         = map[string]CharSet{
		"n": numbers,
		"a": alphaUpper,
		"c": alphaNumericUpper,
		"x": any,
		"d": floats,
	}
	charSetsKeys runeSet = charsetsKeysAsRunes(charSets)
)

type ValidatesPartially interface {
	ValidatePartial(input string, currLine int) (string, error)
}

type Literal struct {
	Chars string
}

func (l Literal) ValidatePartial(input string, currLine int) (string, error) {
	if strings.HasPrefix(input, l.Chars) {
		if len(input) > len(l.Chars) {
			return input[len(l.Chars):], nil
		}

		return "", nil
	}

	return input, fmt.Errorf("expected input to have literal %q", l.Chars)
}

type Optional struct {
	Pattern ValidatesPartially
}

func (o Optional) ValidatePartial(input string, currLine int) (string, error) {
	rest, err := o.Pattern.ValidatePartial(input, currLine)
	if err != nil {
		return input, nil
	}

	return rest, nil
}

type CharGroup struct {
	charSetKey  string
	CharSet     CharSet
	Count       int
	CountStrict bool
}

func (cg CharGroup) countAndStripFloats(input string) (int, string) {
	charCount := 0

	countBeforeDecimal := 0
BeforeDecimalLoop:
	for _, r := range input {
		if countBeforeDecimal == cg.Count {
			break
		}

		switch {
		case r == ',':
			charCount++
			break BeforeDecimalLoop
		case numbers(r):
			charCount++
			countBeforeDecimal++
		default:
			return 0, input
		}
	}

	countAfterDecimal := 0
AfterDecimalLoop:
	for _, r := range input[charCount:] {
		if countBeforeDecimal+countAfterDecimal == cg.Count {
			break
		}

		switch {
		case numbers(r):
			charCount++
			countAfterDecimal++
		default:
			charCount++
			break AfterDecimalLoop
		}
	}

	finalCount := countBeforeDecimal + countAfterDecimal

	// no decimals meaning invalid float
	if countAfterDecimal == 0 {
		return 0, input
	}

	return finalCount, input[charCount:]
}

func (cg CharGroup) countAndStripChars(input string) (int, string) {
	switch cg.charSetKey {
	case "d":
		return cg.countAndStripFloats(input)
	default:
		count := 0

		for _, r := range input {
			if count == cg.Count {
				break
			}
			if !cg.CharSet(r) {
				break
			}

			count++
		}

		return count, input[count:]
	}
}

func (cg CharGroup) ValidatePartial(input string, currLine int) (string, error) {
	count, newInput := cg.countAndStripChars(input)

	switch {
	case count < cg.Count && cg.CountStrict:
		return newInput, fmt.Errorf("expected %d characters within '%s' group, got %d", cg.Count, cg.charSetKey, count)
	default:
		return newInput, nil
	}
}

type Pattern []ValidatesPartially

func (p Pattern) ValidatePartial(input string, currLine int) (string, error) {
	var err error

	for _, v := range p {
		input, err = v.ValidatePartial(input, currLine)
		if err != nil {
			return input, err
		}
	}

	return input, nil
}

func (p Pattern) Validate(input string) error {
	var err error

	for _, pv := range p {
		input, err = pv.ValidatePartial(input, 1)
		if err != nil {
			return fmt.Errorf("input invalid: %w", err)
		}
	}

	if input != "" {
		return fmt.Errorf("incomplete match")
	}

	return nil
}

type LinePattern struct {
	InRange func(line int) bool
	Pattern ValidatesPartially
}

func (lp LinePattern) ValidatePartial(input string, currLine int) (string, error) {
	lines := strings.Split(input, "\n")

	for i := 0; i < len(lines) && lp.InRange(currLine); i++ {
		line := lines[i]

		rest, err := lp.Pattern.ValidatePartial(line, currLine)
		if err != nil {
			return input, fmt.Errorf("line %d: %w", currLine, err)
		}
		if rest != "" {
			return input, fmt.Errorf("line %d: incomplete match", currLine)
		}

		newLineIdx := len(line + "\n")
		if len(input) > newLineIdx {
			input = input[len(line+"\n"):]
		} else {
			input = ""
		}

		currLine++
	}

	return input, nil
}

type OrPattern struct {
	Left  ValidatesPartially
	Right ValidatesPartially
}

func (op OrPattern) ValidatePartial(input string, currLine int) (string, error) {
	errStr := ""

	restLeft, errLeft := op.Left.ValidatePartial(input, currLine)
	if errLeft == nil && restLeft == "" {
		// no need to try the right side, the left covered the entire input successfully
		return restLeft, nil
	}

	restRight, errRight := op.Right.ValidatePartial(input, currLine)

	switch {
	case errLeft == nil && errRight != nil:
		return restLeft, nil
	case errLeft != nil && errRight == nil:
		return restRight, nil
	case errLeft == nil && errRight == nil && len(restLeft) < len(restRight):
		return restLeft, nil
	case errLeft == nil && errRight == nil && len(restLeft) > len(restRight):
		return restRight, nil
	default:
		errStr = fmt.Sprintf("left: %s, right: %s", errLeft, errRight)
		return input, fmt.Errorf("input invalid for or: %s", errStr)
	}
}

// TODO - Support escaping of reserved characters for use as literals
func Parse(input string) (Pattern, error) {
	tokens := lex(input)

	astPattern, err := parse(tokens)
	if err != nil {
		return nil, fmt.Errorf("could not parse pattern: %w", err)
	}

	return process(astPattern), nil
}
