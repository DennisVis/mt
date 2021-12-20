// Copyright (c) 2021 Dennis Vis
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package mt_test

import (
	"testing"
	"time"

	"github.com/DennisVis/mt"
)

func TestTime(t *testing.T) {
	var d mt.Time
	err := d.UnmarshalMT("1504")
	if err != nil {
		t.Error(err)
	}
	if d.Set != true {
		t.Errorf("expected Set to be true, got %t", d.Set)
	}
	if d.Raw != "1504" {
		t.Errorf("expected Raw to be 1504, got %s", d.Raw)
	}
	if d.RawString() != "1504" {
		t.Errorf("expected RawString() to return 1504, got %s", d.RawString())
	}
	if d.String() != "1504" {
		t.Errorf("expected String() to return 1504, got %s", d.String())
	}
	if d.Raw != "1504" {
		t.Errorf("expected Raw to be 1504, got %s", d.Raw)
	}
	if d.Time.Hour() != 15 {
		t.Errorf("expected Hour to be 3, got %d", d.Time.Hour())
	}
	if d.Time.Minute() != 4 {
		t.Errorf("expected Minute to be 4, got %d", d.Time.Minute())
	}

	var d2 mt.Time
	err = d2.UnmarshalMT("150")
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestMonth(t *testing.T) {
	var m mt.Month
	err := m.UnmarshalMT("0102")
	if err != nil {
		t.Error(err)
	}
	if m.Set != true {
		t.Errorf("expected Set to be true")
	}
	if m.Raw != "0102" {
		t.Errorf("expected Raw to be 0102, got %s", m.Raw)
	}
	if m.RawString() != "0102" {
		t.Errorf("expected RawString() to return 0102, got %s", m.RawString())
	}
	if m.String() != "0102" {
		t.Errorf("expected String() to return 0102, got %s", m.String())
	}
	if m.Raw != "0102" {
		t.Errorf("expected Raw to be 0102, got %s", m.Raw)
	}
	if m.Time.Month() != time.January {
		t.Errorf("expected Month to be January, got %s", m.Time.Month().String())
	}
	if m.Time.Day() != 2 {
		t.Errorf("expected Day to be 2, got %d", m.Time.Day())
	}

	var d2 mt.DateTime
	err = d2.UnmarshalMT("010")
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestDate(t *testing.T) {
	var d mt.Date
	err := d.UnmarshalMT("080102")
	if err != nil {
		t.Error(err)
	}
	if d.Set != true {
		t.Errorf("expected Set to be true")
	}
	if d.Raw != "080102" {
		t.Errorf("expected Raw to be 080102, got %s", d.Raw)
	}
	if d.RawString() != "080102" {
		t.Errorf("expected RawString() to return 080102, got %s", d.RawString())
	}
	if d.String() != "080102" {
		t.Errorf("expected String() to return 080102, got %s", d.String())
	}
	if d.Raw != "080102" {
		t.Errorf("expected Raw to be 080102, got %s", d.Raw)
	}
	if d.Time.Year() != 2008 {
		t.Errorf("expected Year to be 2008, got %d", d.Time.Year())
	}
	if d.Time.Month() != time.January {
		t.Errorf("expected Month to be January, got %s", d.Time.Month())
	}
	if d.Time.Day() != 2 {
		t.Errorf("expected Day to be 2, got %d", d.Time.Day())
	}

	var d2 mt.Date
	err = d2.UnmarshalMT("08010")
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestDateTime(t *testing.T) {
	var d mt.DateTime
	err := d.UnmarshalMT("0801021504")
	if err != nil {
		t.Error(err)
	}
	if d.Set != true {
		t.Errorf("expected Set to be true")
	}
	if d.Raw != "0801021504" {
		t.Errorf("expected Raw to be 0801021504, got %s", d.Raw)
	}
	if d.RawString() != "0801021504" {
		t.Errorf("expected RawString() to return 0801021504, got %s", d.RawString())
	}
	if d.String() != "0801021504" {
		t.Errorf("expected String() to return 0801021504, got %s", d.String())
	}
	if d.Raw != "0801021504" {
		t.Errorf("expected Raw to be 01021504, got %s", d.Raw)
	}
	if d.Time.Year() != 2008 {
		t.Errorf("expected Year to be 2008, got %d", d.Time.Year())
	}
	if d.Time.Month() != time.January {
		t.Errorf("expected Month to be January, got %s", d.Time.Month().String())
	}
	if d.Time.Day() != 2 {
		t.Errorf("expected Day to be 2, got %d", d.Time.Day())
	}
	if d.Time.Hour() != 15 {
		t.Errorf("expected Hour to be 15, got %d", d.Time.Hour())
	}
	if d.Time.Minute() != 4 {
		t.Errorf("expected Minute to be 4, got %d", d.Time.Minute())
	}

	var d2 mt.DateTime
	err = d2.UnmarshalMT("080102150")
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestDateTimeSec(t *testing.T) {
	var d mt.DateTimeSec
	err := d.UnmarshalMT("080102150405")
	if err != nil {
		t.Error(err)
	}
	if d.Set != true {
		t.Errorf("expected Set to be true")
	}
	if d.Raw != "080102150405" {
		t.Errorf("expected Raw to be 080102150405, got %s", d.Raw)
	}
	if d.RawString() != "080102150405" {
		t.Errorf("expected RawString() to return 080102150405, got %s", d.RawString())
	}
	if d.String() != "080102150405" {
		t.Errorf("expected String() to return 080102150405, got %s", d.String())
	}
	if d.Raw != "080102150405" {
		t.Errorf("expected Raw to be 080102150405, got %s", d.Raw)
	}
	if d.Time.Year() != 2008 {
		t.Errorf("expected Year to be 2008, got %d", d.Time.Year())
	}
	if d.Time.Month() != time.January {
		t.Errorf("expected Month to be January, got %s", d.Time.Month())
	}
	if d.Time.Day() != 2 {
		t.Errorf("expected Day to be 2, got %d", d.Time.Day())
	}
	if d.Time.Hour() != 15 {
		t.Errorf("expected Hour to be 15, got %d", d.Time.Hour())
	}
	if d.Time.Minute() != 4 {
		t.Errorf("expected Minute to be 4, got %d", d.Time.Minute())
	}
	if d.Time.Second() != 5 {
		t.Errorf("expected Second to be 5, got %d", d.Time.Second())
	}

	var d2 mt.DateTimeSec
	err = d2.UnmarshalMT("08010215040")
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestDateTimeSecCent(t *testing.T) {
	var d mt.DateTimeSecCent
	err := d.UnmarshalMT("080102150405123")
	if err != nil {
		t.Error(err)
	}
	if d.Set != true {
		t.Errorf("expected Set to be true")
	}
	if d.Raw != "080102150405123" {
		t.Errorf("expected Raw to be 080102150405123, got %s", d.Raw)
	}
	if d.RawString() != "080102150405123" {
		t.Errorf("expected RawString() to return 080102150405123, got %s", d.RawString())
	}
	if d.String() != "080102150405123" {
		t.Errorf("expected String() to return 080102150405123, got %s", d.String())
	}
	if d.Raw != "080102150405123" {
		t.Errorf("expected Raw to be 080102150405123, got %s", d.Raw)
	}
	if d.Time.Year() != 2008 {
		t.Errorf("expected Year to be 2008, got %d", d.Time.Year())
	}
	if d.Time.Month() != time.January {
		t.Errorf("expected Month to be January, got %s", d.Time.Month())
	}
	if d.Time.Day() != 2 {
		t.Errorf("expected Day to be 2, got %d", d.Time.Day())
	}
	if d.Time.Hour() != 15 {
		t.Errorf("expected Hour to be 15, got %d", d.Time.Hour())
	}
	if d.Time.Minute() != 4 {
		t.Errorf("expected Minute to be 4, got %d", d.Time.Minute())
	}
	if d.Time.Second() != 5 {
		t.Errorf("expected Second to be 5, got %d", d.Time.Second())
	}
	if d.Time.Nanosecond() != 123000000 {
		t.Errorf("expected Nanosecond to be 60000000, got %d", d.Time.Nanosecond())
	}

	var d2 mt.DateTimeSecCent
	err = d2.UnmarshalMT("080102150405")
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestDateTimeSecOptCent(t *testing.T) {
	var d mt.DateTimeSecOptCent
	err := d.UnmarshalMT("080102150405123")
	if err != nil {
		t.Error(err)
	}
	if d.Set != true {
		t.Errorf("expected Set to be true")
	}
	if d.Raw != "080102150405123" {
		t.Errorf("expected Raw to be 080102150405123, got %s", d.Raw)
	}
	if d.RawString() != "080102150405123" {
		t.Errorf("expected RawString() to return 080102150405123, got %s", d.RawString())
	}
	if d.String() != "080102150405123" {
		t.Errorf("expected String() to return 080102150405123, got %s", d.String())
	}
	if d.Raw != "080102150405123" {
		t.Errorf("expected Raw to be 080102150405123, got %s", d.Raw)
	}
	if d.Time.Year() != 2008 {
		t.Errorf("expected Year to be 2008, got %d", d.Time.Year())
	}
	if d.Time.Month() != time.January {
		t.Errorf("expected Month to be January, got %s", d.Time.Month())
	}
	if d.Time.Day() != 2 {
		t.Errorf("expected Day to be 2, got %d", d.Time.Day())
	}
	if d.Time.Hour() != 15 {
		t.Errorf("expected Hour to be 15, got %d", d.Time.Hour())
	}
	if d.Time.Minute() != 4 {
		t.Errorf("expected Minute to be 4, got %d", d.Time.Minute())
	}
	if d.Time.Second() != 5 {
		t.Errorf("expected Second to be 5, got %d", d.Time.Second())
	}
	if d.Time.Nanosecond() != 123000000 {
		t.Errorf("expected Nanosecond to be 60000000, got %d", d.Time.Nanosecond())
	}

	var d2 mt.DateTimeSecOptCent
	err = d2.UnmarshalMT("080102150405")
	if err != nil {
		t.Error(err)
	}
	if d2.Set != true {
		t.Errorf("expected Set to be true")
	}
	if d2.Raw != "080102150405" {
		t.Errorf("expected Raw to be 080102150405, got %s", d2.Raw)
	}
	if d2.RawString() != "080102150405" {
		t.Errorf("expected RawString() to return 080102150405, got %s", d2.RawString())
	}
	if d2.String() != "080102150405" {
		t.Errorf("expected String() to return 080102150405, got %s", d2.String())
	}
	if d2.Raw != "080102150405" {
		t.Errorf("expected Raw to be 080102150405, got %s", d2.Raw)
	}
	if d2.Time.Year() != 2008 {
		t.Errorf("expected Year to be 2008, got %d", d2.Time.Year())
	}
	if d2.Time.Month() != time.January {
		t.Errorf("expected Month to be January, got %s", d2.Time.Month())
	}
	if d2.Time.Day() != 2 {
		t.Errorf("expected Day to be 2, got %d", d2.Time.Day())
	}
	if d2.Time.Hour() != 15 {
		t.Errorf("expected Hour to be 15, got %d", d2.Time.Hour())
	}
	if d2.Time.Minute() != 4 {
		t.Errorf("expected Minute to be 4, got %d", d2.Time.Minute())
	}
	if d2.Time.Second() != 5 {
		t.Errorf("expected Second to be 5, got %d", d2.Time.Second())
	}

	var d3 mt.DateTimeSecOptCent
	err = d3.UnmarshalMT("08X102150405")
	if err == nil {
		t.Errorf("expected error")
	}

	var d4 mt.DateTimeSecOptCent
	err = d4.UnmarshalMT("08X102150405123")
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestDateTimeOffset(t *testing.T) {
	var d mt.DateTimeOffset
	err := d.UnmarshalMT("0801021504+0100")
	if err != nil {
		t.Error(err)
	}
	if d.Set != true {
		t.Errorf("expected Set to be true")
	}
	if d.Raw != "0801021504+0100" {
		t.Errorf("expected Raw to be 0801021504+0100, got %s", d.Raw)
	}
	if d.RawString() != "0801021504+0100" {
		t.Errorf("expected RawString() to return 0801021504+0100, got %s", d.RawString())
	}
	if d.String() != "0801021504+0100" {
		t.Errorf("expected String() to return 0801021504+0100, got %s", d.String())
	}
	if d.Raw != "0801021504+0100" {
		t.Errorf("expected Raw to be 0801021504+0100, got %s", d.Raw)
	}
	if d.Time.Year() != 2008 {
		t.Errorf("expected Year to be 2008, got %d", d.Time.Year())
	}
	if d.Time.Month() != time.January {
		t.Errorf("expected Month to be January, got %s", d.Time.Month())
	}
	if d.Time.Day() != 2 {
		t.Errorf("expected Day to be 2, got %d", d.Time.Day())
	}
	if d.Time.Hour() != 15 {
		t.Errorf("expected Hour to be 15, got %d", d.Time.Hour())
	}
	if d.Time.Minute() != 4 {
		t.Errorf("expected Minute to be 4, got %d", d.Time.Minute())
	}

	var d2 mt.DateTimeOffset
	err = d2.UnmarshalMT("0801021504=0100")
	if err == nil {
		t.Errorf("expected error")
	}
}
