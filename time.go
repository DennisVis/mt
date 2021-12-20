// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mt

import (
	"fmt"
	"time"
)

const (
	TimeFormatTime            = "1504"
	TimeFormatMonth           = "0102"
	TimeFormatDate            = "060102"
	TimeFormatDateTime        = "0601021504"
	TimeFormatDateTimeSec     = "060102150405"
	TimeFormatDateTimeSecCent = "060102150405.999"
	TimeFormatDateTimeOffset  = "0601021504-0700"
)

type Time struct {
	Set  bool
	Raw  string
	Time time.Time
}

func (d *Time) UnmarshalMT(input string) error {
	t, err := time.Parse(TimeFormatTime, input)
	if err != nil {
		return fmt.Errorf("invalid Time: %w", err)
	}

	d.Set = true
	d.Raw = input
	d.Time = t

	return nil
}

func (d Time) RawString() string {
	return d.Raw
}

func (d Time) String() string {
	return d.RawString()
}

type Month struct {
	Set  bool
	Raw  string
	Time time.Time
}

func (m *Month) UnmarshalMT(input string) error {
	t, err := time.Parse(TimeFormatMonth, input)
	if err != nil {
		return fmt.Errorf("invalid Month: %w", err)
	}

	m.Set = true
	m.Raw = input
	m.Time = t

	return nil
}

func (m Month) RawString() string {
	return m.Raw
}

func (m Month) String() string {
	return m.RawString()
}

type Date struct {
	Set  bool
	Raw  string
	Time time.Time
}

func (d *Date) UnmarshalMT(input string) error {
	t, err := time.Parse(TimeFormatDate, input)
	if err != nil {
		return fmt.Errorf("invalid Date: %w", err)
	}

	d.Set = true
	d.Raw = input
	d.Time = t

	return nil
}

func (d Date) RawString() string {
	return d.Raw
}

func (d Date) String() string {
	return d.RawString()
}

type DateTime struct {
	Set  bool
	Raw  string
	Time time.Time
}

func (d *DateTime) UnmarshalMT(input string) error {
	t, err := time.Parse(TimeFormatDateTime, input)
	if err != nil {
		return fmt.Errorf("invalid DateTime: %w", err)
	}

	d.Set = true
	d.Raw = input
	d.Time = t

	return nil
}

func (d DateTime) RawString() string {
	return d.Raw
}

func (d DateTime) String() string {
	return d.RawString()
}

type DateOrDateTime struct {
	Set  bool
	Raw  string
	Time time.Time
}

func (d *DateOrDateTime) UnmarshalMT(input string) error {
	var t time.Time
	var err error

	if len(input) == 10 {
		t, err = time.Parse(TimeFormatDateTime, input)
		if err != nil {
			return fmt.Errorf("invalid DateOrDateTime date/time: %w", err)
		}
	} else {
		t, err = time.Parse(TimeFormatDate, input)
		if err != nil {
			return fmt.Errorf("invalid DateOrDateTime date: %w", err)
		}
	}

	d.Set = true
	d.Raw = input
	d.Time = t

	return nil
}

func (d DateOrDateTime) RawString() string {
	return d.Raw
}

func (d DateOrDateTime) String() string {
	return d.RawString()
}

type DateTimeSec struct {
	Set  bool
	Raw  string
	Time time.Time
}

func (d *DateTimeSec) UnmarshalMT(input string) error {
	t, err := time.Parse(TimeFormatDateTimeSec, input)
	if err != nil {
		return fmt.Errorf("invalid DateTimeSec: %w", err)
	}

	d.Set = true
	d.Raw = input
	d.Time = t

	return nil
}

func (d DateTimeSec) RawString() string {
	return d.Raw
}

func (d DateTimeSec) String() string {
	return d.RawString()
}

type DateTimeSecCent struct {
	Set  bool
	Raw  string
	Time time.Time
}

func (d *DateTimeSecCent) UnmarshalMT(input string) error {
	// time.Parse needs a decimal point to be able to parse sub-seconds.
	t, err := time.Parse(TimeFormatDateTimeSecCent, input[:12]+"."+input[12:])
	if err != nil {
		return fmt.Errorf("invalid DateTimeSecCent: %w", err)
	}

	d.Set = true
	d.Raw = input
	d.Time = t

	return nil
}

func (d DateTimeSecCent) RawString() string {
	return d.Raw
}

func (d DateTimeSecCent) String() string {
	return d.RawString()
}

type DateTimeSecOptCent struct {
	Set  bool
	Raw  string
	Time time.Time
}

func (d *DateTimeSecOptCent) UnmarshalMT(input string) error {
	var t time.Time
	var err error

	if len(input) == 15 {
		// time.Parse needs a decimal point to be able to parse sub-seconds.
		t, err = time.Parse(TimeFormatDateTimeSecCent, input[:12]+"."+input[12:])
		if err != nil {
			return fmt.Errorf("invalid DateTimeSecOptCent: %w", err)
		}
	} else {
		t, err = time.Parse(TimeFormatDateTimeSec, input)
		if err != nil {
			return fmt.Errorf("invalid DateTimeSecOptCent: %w", err)
		}
	}

	d.Set = true
	d.Raw = input
	d.Time = t

	return nil
}

func (d DateTimeSecOptCent) RawString() string {
	return d.Raw
}

func (d DateTimeSecOptCent) String() string {
	return d.RawString()
}

type DateTimeOffset struct {
	Set  bool
	Raw  string
	Time time.Time
}

func (d *DateTimeOffset) UnmarshalMT(input string) error {
	t, err := time.Parse(TimeFormatDateTimeOffset, input)
	if err != nil {
		return fmt.Errorf("invalid DateTimeSecOffset: %w", err)
	}

	d.Set = true
	d.Raw = input
	d.Time = t

	return nil
}

func (d DateTimeOffset) RawString() string {
	return d.Raw
}

func (d DateTimeOffset) String() string {
	return d.RawString()
}
