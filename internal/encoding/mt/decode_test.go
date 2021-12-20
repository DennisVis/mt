// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mt_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DennisVis/mt/internal/encoding/mt"
	mttest "github.com/DennisVis/mt/testdata"
)

type MTUnmarshaler interface {
	UnmarshalMT(input string) error
}

var errUnmarshalFail = fmt.Errorf("unmarshal fail")

type testSubStruct struct {
	set       bool
	processed bool
}

func (tss *testSubStruct) UnmarshalMT(input string) error {
	tss.processed = true
	return nil
}

type testSubStructInvalid struct {
	set       bool
	processed bool
}

func (tss *testSubStructInvalid) UnmarshalMT(input string) error {
	tss.processed = true
	return errUnmarshalFail
}

type testStruct struct {
	SubField       MTUnmarshaler   `mt:"1"`
	BoolField      bool            `mt:"2"`
	IntField       int             `mt:"3"`
	Int8Field      int8            `mt:"4"`
	Int16Field     int16           `mt:"5"`
	Int32Field     int32           `mt:"6"`
	Int64Field     int64           `mt:"7"`
	UintField      uint            `mt:"8"`
	Uint8Field     uint8           `mt:"9"`
	Uint16Field    uint16          `mt:"10"`
	Uint32Field    uint32          `mt:"11"`
	Uint64Field    uint64          `mt:"12"`
	Float32Field   float32         `mt:"13"`
	Float64Field   float64         `mt:"14"`
	SliceField     []string        `mt:"15"`
	SliceSubField  []MTUnmarshaler `mt:"16"`
	StringField    string          `mt:"17"`
	StringPtrField *string         `mt:"18"`
}

func TestUnmarshalMT(t *testing.T) {
	str := "1"

	for _, test := range []struct {
		name           string
		input          map[string][]string
		factory        func() interface{}
		expectedStruct testStruct
		expectedError  error
	}{
		{
			name: "NotAPointer",
			factory: func() interface{} {
				return testStruct{}
			},
			expectedError: fmt.Errorf("not a pointer"),
		},
		{
			name: "NotANonNilPointer",
			factory: func() interface{} {
				var ts *testStruct
				return ts
			},
			expectedError: fmt.Errorf("not a non-nil pointer"),
		},
		{
			name: "NotAStructPointer",
			factory: func() interface{} {
				str := "1"
				return &str
			},
			expectedError: fmt.Errorf("not a pointer to a struct"),
		},
		{
			name: "MissingStructTag",
			factory: func() interface{} {
				strct := struct {
					Field string
				}{}
				return &strct
			},
		},
		{
			name: "UnknownStructTag",
			factory: func() interface{} {
				strct := struct {
					Field string `mt:"unknown"`
				}{}
				return &strct
			},
			input: map[string][]string{
				"known": {"1"},
			},
		},
		{
			name: "MultiValuesForNonSliceField",
			factory: func() interface{} {
				strct := struct {
					Field string `mt:"1"`
				}{}
				return &strct
			},
			input: map[string][]string{
				"1": {"1", "2"},
			},
			expectedError: fmt.Errorf("multiple values but field is not a slice"),
		},
		{
			name: "InvalidBool",
			factory: func() interface{} {
				strct := struct {
					Field bool `mt:"1"`
				}{}
				return &strct
			},
			input: map[string][]string{
				"1": {"123"},
			},
			expectedError: fmt.Errorf("invalid bool value"),
		},
		{
			name: "InvalidInt",
			factory: func() interface{} {
				strct := struct {
					Field int `mt:"1"`
				}{}
				return &strct
			},
			input: map[string][]string{
				"1": {"bla"},
			},
			expectedError: fmt.Errorf("invalid int value"),
		},
		{
			name: "InvalidUint",
			factory: func() interface{} {
				strct := struct {
					Field uint `mt:"1"`
				}{}
				return &strct
			},
			input: map[string][]string{
				"1": {"bla"},
			},
			expectedError: fmt.Errorf("invalid uint value"),
		},
		{
			name: "InvalidFloat",
			factory: func() interface{} {
				strct := struct {
					Field float32 `mt:"1"`
				}{}
				return &strct
			},
			input: map[string][]string{
				"1": {"bla"},
			},
			expectedError: fmt.Errorf("invalid float value"),
		},
		{
			name: "AllFieldsValid",
			input: map[string][]string{
				"1":  {"test"},
				"2":  {"true"},
				"3":  {"1"},
				"4":  {"1"},
				"5":  {"1"},
				"6":  {"1"},
				"7":  {"1"},
				"8":  {"1"},
				"9":  {"1"},
				"10": {"1"},
				"11": {"1"},
				"12": {"1"},
				"13": {"1.1"},
				"14": {"2.2"},
				"15": {"test1", "test2"},
				"16": {"test1", "test2"},
				"17": {"1"},
				"18": {"1"},
			},
			factory: func() interface{} {
				return &testStruct{
					SubField: &testSubStruct{
						set: true,
					},
				}
			},
			expectedStruct: testStruct{
				SubField: &testSubStruct{
					set:       true,
					processed: true,
				},
				BoolField:    true,
				IntField:     1,
				Int8Field:    1,
				Int16Field:   1,
				Int32Field:   1,
				UintField:    1,
				Uint8Field:   1,
				Uint16Field:  1,
				Uint32Field:  1,
				Float32Field: 1.1,
				Float64Field: 2.2,
				SliceField:   []string{"test1", "test2"},
				SliceSubField: []MTUnmarshaler{
					&testSubStruct{
						set:       true,
						processed: true,
					},
					&testSubStruct{
						set:       true,
						processed: true,
					},
				},
				StringField:    "1",
				StringPtrField: &str,
			},
		},
		{
			name: "SubFieldInvalid",
			input: map[string][]string{
				"1": {"test"},
			},
			factory: func() interface{} {
				return &testStruct{
					SubField: &testSubStructInvalid{
						set: true,
					},
				}
			},
			expectedError: errUnmarshalFail,
		},
		// {
		// 	name: "SliceSubField",
		// 	input: map[string][]string{
		// 		"16": {"test"},
		// 	},
		// 	factory: func() interface{} {
		// 		return &testStruct{
		// 			SliceSubField: []MTUnmarshaler{
		// 				&testSubStruct{
		// 					set: true,
		// 				},
		// 			},
		// 		}
		// 	},
		// 	expectedError: errUnmarshalFail,
		// },
	} {
		// rebind to make sure we can run in parallel
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			v := test.factory()
			err := mt.UnmarshalMT(test.input, v)
			mttest.ValidateError(t, test.expectedError, err)

			if test.expectedError == nil {
				if reflect.DeepEqual(v, test.expectedStruct) {
					t.Errorf("unexpected result: %v", v)
				}
			}
		})
	}
}
