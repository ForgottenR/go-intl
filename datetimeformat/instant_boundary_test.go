package datetimeformat

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/agentable/go-intl/internal/ecma402"
	"github.com/agentable/go-intl/internal/intltest"
	"github.com/agentable/go-intl/option"
)

func TestInstantBoundary(t *testing.T) {
	t.Parallel()
	f, err := New(intltest.LocaleList(t, "en-US"), Options{TimeZone: option.String("UTC")})
	if err != nil {
		t.Fatal(err)
	}
	positive := time.Unix(8640000000000, 0)
	negative := time.Unix(-8640000000000, 0)
	for _, tc := range []struct {
		name    string
		instant time.Time
		valid   bool
	}{
		{"positive bound", positive, true}, {"negative bound", negative, true},
		{"positive beyond submillisecond", positive.Add(time.Nanosecond), false},
		{"negative beyond submillisecond", negative.Add(-time.Nanosecond), false},
		{"negative submillisecond", time.Unix(0, -1), true},
		{"zero Go time", time.Time{}, true},
		{"BCE", time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := f.Format(tc.instant)
			if tc.valid != (err == nil) {
				t.Errorf("Format error = %v, valid = %v", err, tc.valid)
			}
			if !tc.valid && !errors.Is(err, ecma402.ErrInvalidValue) {
				t.Errorf("Format error = %v, want ErrInvalidValue", err)
			}
			_, err = f.FormatToParts(tc.instant)
			if tc.valid != (err == nil) {
				t.Errorf("FormatToParts error = %v, valid = %v", err, tc.valid)
			}
			_, err = f.FormatRange(tc.instant, positive)
			if tc.valid != (err == nil) {
				t.Errorf("FormatRange start error = %v, valid = %v", err, tc.valid)
			}
			if !tc.valid && (err == nil || !strings.Contains(err.Error(), "start")) {
				t.Errorf("start error = %v", err)
			}
			_, err = f.FormatRangeToParts(negative, tc.instant)
			if tc.valid != (err == nil) {
				t.Errorf("FormatRangeToParts end error = %v, valid = %v", err, tc.valid)
			}
			if !tc.valid && (err == nil || !strings.Contains(err.Error(), "end")) {
				t.Errorf("end error = %v", err)
			}
		})
	}
}
