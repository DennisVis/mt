// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mt_test

import (
	"fmt"
	"testing"

	"github.com/DennisVis/mt"
)

func TestErrorString(t *testing.T) {
	for _, test := range []struct {
		name        string
		parseErr    mt.Error
		expectedStr string
	}{
		{
			name:        "SimpleErrorLine1",
			parseErr:    mt.NewError(fmt.Errorf("simple error"), 1),
			expectedStr: "#1: simple error",
		},
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if test.parseErr.Error() != test.expectedStr {
				t.Errorf("expected %s, got %s", test.expectedStr, test.parseErr.Error())
			}
		})
	}
}

func TestErrorsString(t *testing.T) {
	for _, test := range []struct {
		name        string
		parseErr    mt.Errors
		expectedStr string
	}{
		{
			name:        "EmptyErrors",
			parseErr:    mt.Errors{},
			expectedStr: "",
		},
		{
			name: "SimpleErrorLine1",
			parseErr: mt.Errors{
				mt.NewError(fmt.Errorf("simple error"), 1),
				mt.NewError(fmt.Errorf("and another"), 2),
			},
			expectedStr: "mt: Parse errors per message line:\n#1: simple error\n#2: and another",
		},
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if test.parseErr.Error() != test.expectedStr {
				t.Errorf("expected %s, got %s", test.expectedStr, test.parseErr.Error())
			}
		})
	}
}
