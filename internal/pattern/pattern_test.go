// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package pattern_test

import (
	"fmt"
	"testing"

	"github.com/DennisVis/mt/internal/pattern"
	mttest "github.com/DennisVis/mt/testdata"
)

func TestPatternParse(t *testing.T) {
	t.Helper()

	for _, test := range []struct {
		pattern     string
		expectedErr error
	}{
		{
			pattern:     "(/",
			expectedErr: fmt.Errorf("unclosed optional expression"),
		},
		{
			pattern:     "(/(/)",
			expectedErr: fmt.Errorf("unclosed optional expression"),
		},
		{
			pattern:     "(/(/",
			expectedErr: fmt.Errorf("unclosed optional expression"),
		},
		{
			pattern:     "1!a|(",
			expectedErr: fmt.Errorf("unclosed optional expression"),
		},
		{
			pattern:     "2*(1!a",
			expectedErr: fmt.Errorf("unclosed optional expression"),
		},
		{
			pattern:     "2*(1!a",
			expectedErr: fmt.Errorf("unclosed optional expression"),
		},
		{
			pattern:     "2*1!a|(",
			expectedErr: fmt.Errorf("unclosed optional expression"),
		},
		{
			pattern:     "2*1!a|2*(",
			expectedErr: fmt.Errorf("unclosed optional expression"),
		},
		{
			pattern: "(1!a2*1!a)",
		},
		{
			pattern: "(1!a|1!n)",
		},
		{
			pattern: "1!z",
		},
		{
			pattern:     "2**1!z",
			expectedErr: fmt.Errorf("unexpected token *"),
		},
		{
			pattern:     "(2**1!a)",
			expectedErr: fmt.Errorf("unexpected token *"),
		},
		{
			pattern:     "1!a|1!n|2**1!x",
			expectedErr: fmt.Errorf("unexpected token *"),
		},
		{
			pattern:     "(1!n|2**1!a)",
			expectedErr: fmt.Errorf("unexpected token *"),
		},
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(test.pattern, func(t *testing.T) {
			t.Parallel()

			_, err := pattern.Parse(test.pattern)
			mttest.ValidateError(t, test.expectedErr, err)
		})
	}
}

func TestPattern(t *testing.T) {
	t.Helper()

	for _, test := range []struct {
		pattern     string
		input       string
		expectedErr error
	}{
		{
			pattern:     "x16x",
			input:       "y1234567890",
			expectedErr: fmt.Errorf("expected input to have literal \"x\""),
		},
		{
			pattern: "x16x",
			input:   "x1234567890",
		},
		{
			pattern: "16",
			input:   "16",
		},
		{
			pattern: "//(//)",
			input:   "//",
		},
		{
			pattern: "//(//)",
			input:   "////",
		},
		{
			pattern: "//((/)/)",
			input:   "////",
		},
		{
			pattern: "//|^^",
			input:   "//",
		},
		{
			pattern: "//|^^",
			input:   "^^",
		},
		{
			pattern:     "16x",
			input:       "abc123,*",
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			pattern: "16x",
			input:   "1234567890ABCDEF",
		},
		{
			pattern:     "3!a",
			input:       "ABc",
			expectedErr: fmt.Errorf("expected 3 characters within 'a' group, got 2"),
		},
		{
			pattern: "/3!a",
			input:   "/ABC",
		},
		{
			pattern:     "3d",
			input:       "0,,",
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			pattern:     "3d",
			input:       "0,aa",
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			pattern:     "3d",
			input:       "0,000,00",
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			pattern: "2d",
			input:   "0,0",
		},
		{
			pattern:     "2d",
			input:       "0,00",
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			pattern: "3d",
			input:   "0,00",
		},
		{
			pattern:     "3!d",
			input:       "00,00",
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			pattern: "3!d",
			input:   "0,00",
		},
		{
			pattern:     "3!d",
			input:       "0,000,00",
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			pattern: "3!d3!d",
			input:   "0,000,00",
		},
		{
			pattern:     "/3!a",
			input:       "ABC",
			expectedErr: fmt.Errorf("expected input to have literal \"/\""),
		},
		{
			pattern: "/",
			input:   "/",
		},
		{
			pattern: "(/)3!a",
			input:   "/ABC",
		},
		{
			pattern: "(/)3!a",
			input:   "ABC",
		},
		{
			pattern: "(/)",
			input:   "/",
		},
		{
			pattern: "(/)",
			input:   "",
		},
		{
			pattern: "(/(/))",
			input:   "",
		},
		{
			pattern: "(/(/))",
			input:   "/",
		},
		{
			pattern: "(/(/))",
			input:   "//",
		},
		{
			pattern: "(3!a)",
			input:   "",
		},
		{
			pattern: "(3!a)",
			input:   "ABC",
		},
		{
			pattern: "(2*3!a)",
			input:   "",
		},
		{
			pattern: "(2*3!a)",
			input:   "ABC",
		},
		{
			pattern: "(2*3!a)",
			input:   "ABC\nDEF",
		},
		{
			pattern: "2!c26!n",
			input:   "PL25106000760000888888888888",
		},
		{
			pattern: "8!c/12!n",
			input:   "BPHKPLPK/320000752973",
		},
		{
			pattern: "1!a6!n3!a15d",
			input:   "C020628PLN3481,35",
		},
		{
			pattern:     "5!n(/)3!n",
			input:       "somethingelse",
			expectedErr: fmt.Errorf("expected 5 characters within 'n' group, got 0"),
		},
		{
			pattern: "2!c26!n|8!c/12!n",
			input:   "PL25106000760000888888888888",
		},
		{
			pattern: "2!c26!n|8!c/12!n",
			input:   "BPHKPLPK/320000752973",
		},
		{
			pattern:     "2!c26!n|8!c/12!n",
			input:       "BPHKPLPK320000752973",
			expectedErr: fmt.Errorf("input invalid for or"),
		},
		{
			pattern: "6*65x",
			input:   "abc\nefg\nhij",
		},
		{
			pattern: "1*6!n4!n2a|8n1!a3!c1*(//)16x",
			input:   "1234561234AB\n//1010001272972001",
		},
		{
			pattern: "1*6!n4!n2a|8n1!a3!c1*(//)16x",
			input:   "12345678AABC\n//1010001272972001",
		},
		{
			pattern:     "1*6!n4!n2a|8n1!a3!c1*(//)16x",
			input:       "12345678AABC\n//10100012729720011",
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			pattern:     "1*6!n4!n2a|8n1!a3!c1*(//)16x",
			input:       "12345678AAB\n//1010001272972001",
			expectedErr: fmt.Errorf("input invalid for or"),
		},
		{
			pattern: "2*3!a2*3!n",
			input:   "ABC\nDEF\n123\n456",
		},
		{
			pattern: "1!a|2!n|3!d1*1!a|2!n|3!d",
			input:   "A\nB",
		},
		{
			pattern: "2*1!a|2*1!n",
			input:   "A\nB",
		},
		{
			pattern: "2*1!a|2*1!n",
			input:   "1\n2",
		},
		{
			pattern: "1!a|(2*1!n|2*1!a)",
			input:   "A",
		},
		{
			pattern: "1!a|(2*1!n|2*1!a)",
			input:   "1\n2",
		},
		{
			pattern: "1!a|(2*1!n|2*1!a)",
			input:   "A\nB",
		},
		{
			pattern: "2!a|(1!n)1!a",
			input:   "AB",
		},
		{
			pattern: "2!a|1!a",
			input:   "AB",
		},
		{
			pattern: "1!n|2!a",
			input:   "AB",
		},
		{
			pattern:     "1!n|2!a",
			input:       "12",
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			pattern:     "2!n|1!n",
			input:       "123",
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			pattern: "2!a|(1!n)1!a",
			input:   "A",
		},
		{
			pattern: "2!a|(1!n)1!a",
			input:   "1A",
		},
		{
			pattern: "1!a|2!n|3!d1*1!a|2!n|3!d",
			input:   "12\n32",
		},
		{
			pattern:     "1!d",
			input:       "1,",
			expectedErr: fmt.Errorf("expected 1 characters within 'd' group, got 0"),
		},
		{
			pattern:     "1!d",
			input:       "x,0",
			expectedErr: fmt.Errorf("expected 1 characters within 'd' group, got 0"),
		},
		{
			pattern: "1!a|2!n|3!d1*1!a|2!n|3!d",
			input:   "1,23\n4,56",
		},
		{
			pattern:     "1!a|2!n|3!d1*1!a|2!n|3!d",
			input:       "1\n32",
			expectedErr: fmt.Errorf("input invalid for or"),
		},
		{
			pattern:     "1!a|2!n|3!d1*1!a|2!n|3!d",
			input:       "12\n3",
			expectedErr: fmt.Errorf("input invalid for or"),
		},
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(fmt.Sprintf("%q:%q", test.pattern, test.input), func(t *testing.T) {
			t.Parallel()

			ptrn, err := pattern.Parse(test.pattern)
			if err != nil {
				t.Fatal(err)
			}

			err = ptrn.Validate(test.input)
			mttest.ValidateError(t, test.expectedErr, err)
		})
	}
}
