// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mt

import "fmt"

// Error is used when parsing of an input encounters a problem.
//
// A parse error will generally not stop the parsing process, as the remaining messages will attempted to be parsed.
//
// Any and all parse errors per input should be aggregated by this library and returned to the caller.
type Error struct {
	line  int
	cause error
}

// NewError creates a new parse error.
func NewError(cause error, line int) Error {
	return Error{line: line, cause: cause}
}

// Cause returns the underlying error.
func (e Error) Cause() error {
	return e.cause
}

// Line returns the line in the input where the error occurred.
func (e Error) Line() int {
	return e.line
}

// String returns the string representation of the parse error.
func (e Error) String() string {
	return fmt.Sprintf("#%d: %s", e.line, e.Cause())
}

// Error implements the Error interface.
func (e Error) Error() string {
	return e.String()
}

// Errors is a custom error type that is used for aggregating Error's into one error.
type Errors []Error

// String returns the string representation of a group of parse errors.
func (es Errors) String() string {
	if len(es) == 0 {
		return ""
	}

	str := "mt: Parse errors per message line:"

	for _, pe := range es {
		str += "\n" + pe.String()
	}

	return str
}

// Error implements the Error interface.
func (es Errors) Error() string {
	return es.String()
}
