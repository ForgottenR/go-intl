package datetimeformat

import (
	"testing"
	"time"
)

func mustDateFormat(t *testing.T, f *DateTimeFormat, instant time.Time) string {
	t.Helper()
	value, err := f.Format(instant)
	if err != nil {
		t.Fatal(err)
	}
	return value
}
func mustDateFormatToParts(t *testing.T, f *DateTimeFormat, instant time.Time) []Part {
	t.Helper()
	value, err := f.FormatToParts(instant)
	if err != nil {
		t.Fatal(err)
	}
	return value
}
func mustDateFormatRange(t *testing.T, f *DateTimeFormat, start, end time.Time) string {
	t.Helper()
	value, err := f.FormatRange(start, end)
	if err != nil {
		t.Fatal(err)
	}
	return value
}
func mustDateFormatRangeToParts(t *testing.T, f *DateTimeFormat, start, end time.Time) []RangePart {
	t.Helper()
	value, err := f.FormatRangeToParts(start, end)
	if err != nil {
		t.Fatal(err)
	}
	return value
}
