// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package validate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/DennisVis/mt/internal/pattern"
)

type Validator interface {
	Validate(interface{}) ValidationError
}

type validator struct {
	typeName string
	items    validationItems
}

type rawStringer interface {
	RawString() string
}

type validationItem struct {
	label     string
	field     string
	mandatory bool
	dive      bool
	pattern   pattern.Pattern
	items     validationItems
}

type validationItems map[string]validationItem

func createItem(label, mandatoryStr, patternStr, fieldName string) (validationItem, error) {
	i := validationItem{
		label: label,
		field: fieldName,
	}

	var mandatory bool
	switch {
	case mandatoryStr == "M":
		mandatory = true
	case mandatoryStr == "O":
		mandatory = false
	default:
		return i, fmt.Errorf("mt tag for field %s needs M or O as second part", fieldName)
	}

	i.mandatory = mandatory

	if patternStr == "dive" {
		i.dive = true
	} else {
		ptrn, err := pattern.Parse(patternStr)
		if err != nil {
			return i, fmt.Errorf("mt tag for field %s contained invalid pattern %q: %w", fieldName, patternStr, err)
		}

		i.pattern = ptrn
	}

	return i, nil
}

func diveIntoStruct(rv reflect.Value) (validationItems, error) {
	subItems := make(validationItems)

	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		sf := rt.Field(i)
		fv := rv.Field(i)

		if !fv.CanInterface() {
			continue
		}

		fieldName := sf.Name

		structTag, ok := sf.Tag.Lookup("mt")
		if !ok || structTag == "" {
			continue
		}

		tagSplit := strings.Split(structTag, ",")
		if len(tagSplit) < 2 {
			return subItems, fmt.Errorf("mt tag for sub field %s needs at least 2 parts: %s", fieldName, structTag)
		}

		i, err := createItem("", tagSplit[0], tagSplit[1], fieldName)
		if err != nil {
			return subItems, err
		}

		subItems[fieldName] = i
	}

	return subItems, nil
}

func CreateValidatorForStruct(strct interface{}) (Validator, error) {
	v := validator{}

	rt := reflect.TypeOf(strct)
	rv := reflect.ValueOf(strct)

	// if we were given a pointer, dereference it
	if rt.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rv.Type()
	}

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct: %s", rt)
	}

	v.typeName = rt.Name()
	v.items = make(validationItems)

	for i := 0; i < rv.NumField(); i++ {
		sf := rt.Field(i)
		fv := rv.Field(i)

		if !fv.CanInterface() {
			continue
		}

		fieldName := sf.Name

		structTag, ok := sf.Tag.Lookup("mt")
		if !ok || structTag == "" {
			continue
		}

		tagSplit := strings.Split(structTag, ",")
		if len(tagSplit) < 3 {
			return nil, fmt.Errorf("mt tag for field %s needs at least 3 parts: %s", fieldName, structTag)
		}

		i, err := createItem(tagSplit[0], tagSplit[1], tagSplit[2], fieldName)
		if err != nil {
			return nil, err
		}

		v.items[fieldName] = i

		switch {
		case fv.Kind() == reflect.Struct && i.dive:
			subItems, err := diveIntoStruct(fv)
			if err != nil {
				return nil, fmt.Errorf("could not dive into struct field %s: %w", fieldName, err)
			}

			i.items = subItems

			v.items[fieldName] = i
		case fv.Kind() == reflect.Slice && fv.Type().Elem().Kind() == reflect.Struct && i.dive:
			subItems, err := diveIntoStruct(reflect.New(fv.Type().Elem()).Elem())
			if err != nil {
				return nil, fmt.Errorf("could not dive into struct slice field %s: %w", fieldName, err)
			}

			i.items = subItems

			v.items[fieldName] = i
		}
	}

	return &v, nil
}

func MustCreateValidatorForStruct(strct interface{}) Validator {
	v, err := CreateValidatorForStruct(strct)
	if err != nil {
		panic(err)
	}

	return v
}

func isRawStringer(rv reflect.Value) bool {
	_, ok := rv.Interface().(rawStringer)
	return ok
}

func valueToString(rv reflect.Value) string {
	var v string

	kind := rv.Kind()
	switch {
	case isRawStringer(rv):
		v = rv.Interface().(rawStringer).RawString()
	case kind == reflect.Int, kind == reflect.Int8, kind == reflect.Int16, kind == reflect.Int32, kind == reflect.Int64:
		v = strconv.FormatInt(rv.Int(), 10)
	case kind == reflect.Uint, kind == reflect.Uint8, kind == reflect.Uint16, kind == reflect.Uint32, kind == reflect.Uint64:
		v = strconv.FormatUint(rv.Uint(), 10)
	case kind == reflect.Float32:
		v = strings.ReplaceAll(strconv.FormatFloat(rv.Float(), 'f', 2, 32), ".", ",")
	case kind == reflect.Float64:
		v = strings.ReplaceAll(strconv.FormatFloat(rv.Float(), 'f', 2, 64), ".", ",")
	case kind == reflect.String:
		v = rv.String()
	}

	return v
}

func validateValue(item validationItem, rv reflect.Value) ValidationError {
	val := valueToString(rv)
	if (val == "" || val == "0") && item.mandatory {
		return valueError{fmt.Errorf("empty mandatory field %s", item.field)}
	}
	if val == "" {
		return nil
	}

	err := item.pattern.Validate(val)
	if err != nil {
		return valueError{fmt.Errorf("pattern validation failed: %w", err)}
	}

	return nil
}

func isUnsupportedType(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Bool, reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.Map, reflect.UnsafePointer:
		return true
	}

	return false
}

func validateMember(item validationItem, name string, rv reflect.Value) ValidationError {
	rt := rv.Type()

	switch {
	case rt.Kind() == reflect.Ptr && rv.IsNil() && item.mandatory:
		return valueError{fmt.Errorf("empty mandatory field %s", name)}
	case rt.Kind() == reflect.Ptr && rv.IsNil():
		return nil
	case rt.Kind() == reflect.Ptr:
		rv = rv.Elem()
	}

	shouldDive := item.dive

	switch {
	case isUnsupportedType(rv):
		return nil
	case rv.Kind() == reflect.Struct && shouldDive:
		return validateStruct(item.items, rv)
	case rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array:
		return validateSlice(item, name, rv)
	default:
		return validateValue(item, rv)
	}
}

func validateSlice(item validationItem, name string, rv reflect.Value) ValidationError {
	errors := make(validationErrors, 0)

	for i := 0; i < rv.Len(); i++ {
		fv := rv.Index(i)

		err := validateMember(item, name, fv)
		if err != nil {
			errors = append(errors, validationError{
				field: item.field + "[" + strconv.Itoa(i) + "]",
				label: item.label,
				err:   err,
			})
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func validateStruct(items validationItems, rv reflect.Value) ValidationError {
	errors := make(validationErrors, 0)

	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		sf := rt.Field(i)

		item, ok := items[sf.Name]
		if !ok {
			continue
		}

		err := validateMember(item, sf.Name, fv)
		if err != nil {
			errors = append(errors, validationError{
				field: item.field,
				label: item.label,
				err:   err,
			})
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (v *validator) Validate(strct interface{}) ValidationError {
	rv := reflect.ValueOf(strct)
	rt := rv.Type()

	// if we were given a pointer, dereference it
	if rt.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rv.Type()
	}

	if rt.Kind() != reflect.Struct {
		return valueError{fmt.Errorf("not a struct: %s", rt)}
	}

	if rt.Name() != v.typeName {
		return valueError{fmt.Errorf("validator is for type %s, given type %s", v.typeName, rt.Name())}
	}

	err := validateStruct(v.items, rv)

	return err
}
