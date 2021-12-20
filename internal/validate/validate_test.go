// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package validate_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DennisVis/mt/internal/validate"
	mttest "github.com/DennisVis/mt/testdata"
)

type testSubStruct struct {
	SubStringVal string `mt:"O,16!x"`
}

type testRawStringer string

func (trs testRawStringer) RawString() string {
	return string(trs)
}

type testStruct struct {
	privateField         string          `mt:"0,O,1!a"`
	StringVal            string          `mt:"1,M,16!x"`
	StringValOptional    string          `mt:"1,O,16!x"`
	StructVal            testSubStruct   `mt:"2,O,dive"`
	StructSliceVal       []testSubStruct `mt:"3,O,dive"`
	StringSliceVal       []string        `mt:"4,O,16!x"`
	IntVal               int             `mt:"5,O,4!n"`
	UintVal              uint            `mt:"6,O,4!n"`
	Float32Val           float32         `mt:"7,O,4!d"`
	Float64Val           float64         `mt:"8,O,4!d"`
	BoolVal              bool            `mt:"9,O,1!n"`
	StringerVal          testRawStringer `mt:"10,O,3!a"`
	StringPtrVal         *string         `mt:"11,M,16!x"`
	StringPtrValOptional *string         `mt:"11,O,16!x"`
}

func newFilledTestStruct() testStruct {
	str16x := strings.Repeat("x", 16)
	ss := testSubStruct{
		SubStringVal: str16x,
	}

	return testStruct{
		privateField:         str16x,
		StringVal:            str16x,
		StringValOptional:    str16x,
		StructVal:            ss,
		StructSliceVal:       []testSubStruct{ss, ss, ss},
		StringSliceVal:       []string{str16x, str16x, str16x},
		IntVal:               1234,
		UintVal:              1234,
		Float32Val:           12.3,
		Float64Val:           12.3,
		BoolVal:              true,
		StringerVal:          testRawStringer("ABC"),
		StringPtrVal:         &str16x,
		StringPtrValOptional: &str16x,
	}
}

func createTestStruct(fn func(*testStruct)) testStruct {
	return func(ts testStruct) testStruct {
		if fn != nil {
			fn(&ts)
		}
		return ts
	}(newFilledTestStruct())
}

func createTestStructPtr(fn func(*testStruct)) *testStruct {
	ts := createTestStruct(fn)
	return &ts
}

func TestCreate(t *testing.T) {
	for _, test := range []struct {
		name        string
		createFrom  interface{}
		expectedErr error
	}{
		{
			name:        "NotAStruct",
			createFrom:  []string{"1"},
			expectedErr: fmt.Errorf("not a struct"),
		},
		{
			name: "TooFewStructTagParts",
			createFrom: struct {
				StringVal string `mt:"1"`
			}{
				StringVal: "1",
			},
			expectedErr: fmt.Errorf("tag for field StringVal needs at least 3 parts"),
		},
		{
			name: "InvalidMandatoryValue",
			createFrom: struct {
				StringVal string `mt:"1,S,1!a"`
			}{
				StringVal: "1",
			},
			expectedErr: fmt.Errorf("tag for field StringVal needs M or O as second part"),
		},
		{
			name: "InvalidPattern",
			createFrom: struct {
				StringVal string `mt:"1,M,1**"`
			}{
				StringVal: "1",
			},
			expectedErr: fmt.Errorf("tag for field StringVal contained invalid pattern"),
		},
		{
			name: "InvalidPrivateFieldPattern",
			createFrom: struct {
				stringVal string `mt:"1,M,1**"`
			}{
				stringVal: "1",
			},
		},
		{
			name: "InvalidSubStructStructTag",
			createFrom: struct {
				StructVal struct {
					SubStringVal string `mt:"S"`
				} `mt:"1,M,dive"`
			}{
				StructVal: struct {
					SubStringVal string `mt:"S"`
				}{
					SubStringVal: "1",
				},
			},
			expectedErr: fmt.Errorf("tag for sub field SubStringVal needs at least 2 parts"),
		},
		{
			name: "InvalidSubStructMandatoryValue",
			createFrom: struct {
				StructVal struct {
					SubStringVal string `mt:"S,1!a"`
				} `mt:"1,M,dive"`
			}{
				StructVal: struct {
					SubStringVal string `mt:"S,1!a"`
				}{
					SubStringVal: "1",
				},
			},
			expectedErr: fmt.Errorf("tag for field SubStringVal needs M or O as second part"),
		},
		{
			name: "InvalidSubStructStructPattern",
			createFrom: struct {
				StructVal struct {
					SubStringVal string `mt:"M,1**"`
				} `mt:"1,M,dive"`
			}{
				StructVal: struct {
					SubStringVal string `mt:"M,1**"`
				}{
					SubStringVal: "1",
				},
			},
			expectedErr: fmt.Errorf("tag for field SubStringVal contained invalid pattern"),
		},
		{
			name: "InvalidSubStructPrivateFieldPattern",
			createFrom: struct {
				StructVal struct {
					stringVal string `mt:"M,1**"`
				} `mt:"1,M,dive"`
			}{
				StructVal: struct {
					stringVal string `mt:"M,1**"`
				}{
					stringVal: "1",
				},
			},
		},
		{
			name: "InvalidStructSliceStructTag",
			createFrom: struct {
				StructSliceVal []struct {
					SliceStringVal string `mt:"S"`
				} `mt:"1,M,dive"`
			}{
				StructSliceVal: []struct {
					SliceStringVal string `mt:"S"`
				}{
					{
						SliceStringVal: "1",
					},
				},
			},
			expectedErr: fmt.Errorf("tag for sub field SliceStringVal needs at least 2 parts"),
		},
		{
			name:       "Valid",
			createFrom: testStruct{},
		},
		{
			name:       "ValidPtr",
			createFrom: &testStruct{},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := validate.CreateValidatorForStruct(test.createFrom)
			mttest.ValidateError(t, test.expectedErr, err)
		})
	}
}

func TestMustCreate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("must create did not panic")
		}
	}()

	validate.MustCreateValidatorForStruct([]string{"1"})
}

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		name        string
		createFrom  interface{}
		input       interface{}
		expectedErr error
	}{
		{
			name:        "NotAStruct",
			createFrom:  testStruct{},
			input:       []string{"1"},
			expectedErr: fmt.Errorf("not a struct"),
		},
		{
			name:        "DifferentStruct",
			createFrom:  testStruct{},
			input:       testSubStruct{},
			expectedErr: fmt.Errorf("validator is for type testStruct, given type testSubStruct"),
		},
		{
			name:       "MissingMandatory",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				ts.StringVal = ""
			}),
			expectedErr: fmt.Errorf("empty mandatory field"),
		},
		{
			name:       "MissingOptional",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				ts.StringValOptional = ""
			}),
		},
		{
			name:       "OptionalPtrNil",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				ts.StringPtrValOptional = nil
			}),
		},
		{
			name:       "MandatoryPtrNil",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				ts.StringPtrVal = nil
			}),
			expectedErr: fmt.Errorf("empty mandatory field"),
		},
		{
			name:       "InvalidPrivateField",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				ts.privateField = "123"
			}),
		},
		{
			name:       "TooManyStringCharsInStringVal",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				ts.StringVal = "123456789012345678901"
			}),
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			name:       "TooManyStringCharsInStringPtrVal",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				str := "123456789012345678901"
				ts.StringPtrVal = &str
			}),
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			name:       "TooManyStringCharsInStructValStringVal",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				ts.StructVal.SubStringVal = "123456789012345678901"
			}),
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			name:       "TooManyStringCharsInStructSliceValStringVal",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				ts.StructSliceVal = []testSubStruct{
					{
						SubStringVal: "123456789012345678901",
					},
				}
			}),
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			name:       "TooManyStringCharsInStingSliceVal",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				ts.StringSliceVal = []string{
					"1234567890123456",
					"123456789012345678901",
				}
			}),
			expectedErr: fmt.Errorf("incomplete match"),
		},
		{
			name:       "InvalidInt",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				ts.IntVal = 123
			}),
			expectedErr: fmt.Errorf("expected 4 characters within 'n' group, got 3"),
		},
		{
			name:       "InvalidUint",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				ts.UintVal = 123
			}),
			expectedErr: fmt.Errorf("expected 4 characters within 'n' group, got 3"),
		},
		{
			name:       "InvalidFloat32",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				ts.Float32Val = 1.2
			}),
			expectedErr: fmt.Errorf("expected 4 characters within 'd' group, got 3"),
		},
		{
			name:       "InvalidFloat64",
			createFrom: testStruct{},
			input: createTestStruct(func(ts *testStruct) {
				ts.Float64Val = 1.2
			}),
			expectedErr: fmt.Errorf("expected 4 characters within 'd' group, got 3"),
		},
		{
			name:       "Valid",
			createFrom: testStruct{},
			input:      newFilledTestStruct(),
		},
		{
			name:       "ValidPtr",
			createFrom: testStruct{},
			input:      createTestStructPtr(nil),
		},
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			v := validate.MustCreateValidatorForStruct(test.createFrom)

			err := v.Validate(test.input)
			mttest.ValidateError(t, test.expectedErr, err)

			if err != nil {
				t.Logf("%s", err.Error())
			}
		})
	}
}
