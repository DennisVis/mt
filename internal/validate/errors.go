// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package validate

import "strings"

type ValidationError interface {
	Error() string
	IndentError(indent string) string
}

type valueError struct {
	err error
}

func (e valueError) IndentError(indent string) string {
	return indent + e.err.Error()
}

func (e valueError) Error() string {
	return e.IndentError("")
}

type validationError struct {
	field string
	label string
	err   ValidationError
}

func (ve validationError) IndentError(indent string) string {
	str := ""

	if ve.field != "" {
		str += ve.field
	}

	if ve.label != "" {
		str += "|" + ve.label + "|"
	}

	_, ok := ve.err.(valueError)
	if ok {
		return indent + str + ": " + ve.err.Error()
	} else {
		return indent + str + ":\n" + ve.err.IndentError(indent+indent)
	}
}

type validationErrors []validationError

func (ves validationErrors) IndentError(indent string) string {
	str := ""
	for _, err := range ves {
		str += err.IndentError(indent) + "\n"
	}

	return strings.TrimRight(str, "\n")
}

func (ves validationErrors) Error() string {
	return ves.IndentError("\t")
}
