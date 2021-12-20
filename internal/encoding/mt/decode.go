// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mt

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type MTUnmarshaler interface {
	UnmarshalMT(input string) error
}

func toUnmarshaler(rval reflect.Value) (MTUnmarshaler, bool) {
	switch {
	case !rval.CanAddr() || !rval.CanInterface():
		return nil, false
	case rval.Kind() == reflect.Ptr && rval.IsNil():
		return nil, false
	case rval.Kind() != reflect.Ptr && rval.Kind() != reflect.Interface && rval.Kind() != reflect.Struct:
		return nil, false
	case rval.Kind() == reflect.Interface && rval.Type().Name() == "MTUnmarshaler" && !rval.IsNil():
		return rval.Interface().(MTUnmarshaler), true
	default:
		um, ok := rval.Interface().(MTUnmarshaler)
		if ok {
			return um, true
		}

		um, ok = rval.Addr().Interface().(MTUnmarshaler)

		return um, ok
	}
}

func isUnmarshaler(rval reflect.Value) bool {
	_, ok := toUnmarshaler(rval)
	return ok
}

func useUnmarshaler(val string, rval reflect.Value) error {
	um, _ := toUnmarshaler(rval)

	err := um.UnmarshalMT(val)
	if err != nil {
		return fmt.Errorf("decoding failed: %w", err)
	}

	return nil
}

func unmarshalBool(val string, rval reflect.Value) error {
	b, err := strconv.ParseBool(val)
	if err != nil {
		return fmt.Errorf("invalid bool value: %w", err)
	}

	rval.SetBool(b)

	return nil
}

func unmarshalInt(val string, rval reflect.Value, bitSize int) error {
	i, err := strconv.ParseInt(val, 10, bitSize)
	if err != nil {
		return fmt.Errorf("invalid int value: %w", err)
	}

	rval.SetInt(i)

	return nil
}

func unmarshalUint(val string, rval reflect.Value, bitSize int) error {
	i, err := strconv.ParseUint(val, 10, bitSize)
	if err != nil {
		return fmt.Errorf("invalid uint value: %w", err)
	}

	rval.SetUint(i)

	return nil
}

func unmarshalFloat(val string, rval reflect.Value, bitSize int) error {
	f, err := strconv.ParseFloat(val, bitSize)
	if err != nil {
		return fmt.Errorf("invalid float value: %w", err)
	}

	rval.SetFloat(f)

	return nil
}

func unmarshalSlice(vals []string, itemName string, rval reflect.Value) error {
	elType := rval.Type().Elem()

	for _, v := range vals {
		ins := reflect.New(elType).Elem()

		err := unmarshalItem([]string{v}, itemName, ins)
		if err != nil {
			return fmt.Errorf("decoding failed for slice item: %w", err)
		}

		reflect.Append(rval, ins)
	}

	return nil
}

func unmarshalString(val string, rval reflect.Value) error {
	rval.SetString(val)
	return nil
}

func unmarshalItem(vals []string, itemName string, rval reflect.Value) error {
	if len(vals) > 1 && rval.Kind() != reflect.Slice {
		return fmt.Errorf("multiple values but field is not a slice")
	}

	var err error
	switch {
	case isUnmarshaler(rval):
		err = useUnmarshaler(vals[0], rval)
	case rval.Kind() == reflect.Bool:
		err = unmarshalBool(vals[0], rval)
	case rval.Kind() == reflect.Int:
		err = unmarshalInt(vals[0], rval, 0)
	case rval.Kind() == reflect.Int8:
		err = unmarshalInt(vals[0], rval, 8)
	case rval.Kind() == reflect.Int16:
		err = unmarshalInt(vals[0], rval, 16)
	case rval.Kind() == reflect.Int32:
		err = unmarshalInt(vals[0], rval, 32)
	case rval.Kind() == reflect.Int64:
		err = unmarshalInt(vals[0], rval, 64)
	case rval.Kind() == reflect.Uint:
		err = unmarshalUint(vals[0], rval, 0)
	case rval.Kind() == reflect.Uint8:
		err = unmarshalUint(vals[0], rval, 8)
	case rval.Kind() == reflect.Uint16:
		err = unmarshalUint(vals[0], rval, 16)
	case rval.Kind() == reflect.Uint32:
		err = unmarshalUint(vals[0], rval, 32)
	case rval.Kind() == reflect.Uint64:
		err = unmarshalUint(vals[0], rval, 64)
	case rval.Kind() == reflect.Float32:
		err = unmarshalFloat(vals[0], rval, 32)
	case rval.Kind() == reflect.Float64:
		err = unmarshalFloat(vals[0], rval, 64)
	case rval.Kind() == reflect.Slice:
		err = unmarshalSlice(vals, itemName, rval)
	case rval.Kind() == reflect.String:
		err = unmarshalString(vals[0], rval)
	default:
		err = fmt.Errorf("unsupported type: %v", rval.Type())
	}
	if err != nil {
		return fmt.Errorf("decoding failed: %w", err)
	}

	return nil
}

func UnmarshalMT(fields map[string][]string, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer: %s", reflect.TypeOf(v))
	}
	if rv.IsNil() {
		return fmt.Errorf("not a non-nil pointer: %s", reflect.TypeOf(v))
	}

	rdv := reflect.Indirect(rv)
	rdt := rdv.Type()
	if rdt.Kind() != reflect.Struct {
		return fmt.Errorf("not a pointer to a struct: %s", reflect.TypeOf(v))
	}

	for i := 0; i < rdv.NumField(); i++ {
		fv := rdv.Field(i)
		sf := rdt.Field(i)

		structTag, ok := sf.Tag.Lookup("mt")
		if !ok || structTag == "" {
			continue
		}

		tagSplit := strings.Split(structTag, ",")
		tag := tagSplit[0]

		vals, ok := fields[tag]
		if !ok || len(vals) < 1 {
			continue
		}

		err := unmarshalItem(vals, sf.Name, fv)
		if err != nil {
			return fmt.Errorf("decoding failed for tag %s, field %s: %w", tag, sf.Name, err)
		}
	}

	return nil
}
