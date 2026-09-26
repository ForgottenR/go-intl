package datetimeformat

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/agentable/go-intl/internal/intltest"
	"github.com/agentable/go-intl/option"
)

func TestMonthPatternContext(t *testing.T) {
	t.Parallel()
	start := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 2, 5, 0, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		locale, width                       string
		alone, format, second, secondFormat string
	}{
		{"ru", "long", "январь", "января", "февраль", "февраля"},
		{"ru", "short", "янв.", "янв.", "февр.", "февр."},
		{"ru", "narrow", "Я", "Я", "Ф", "Ф"},
		{"ru", "numeric", "1", "01", "2", "02"},
		{"pl", "long", "styczeń", "stycznia", "luty", "lutego"},
		{"pl", "short", "sty", "sty", "lut", "lut"},
		{"pl", "narrow", "S", "s", "L", "l"},
		{"pl", "numeric", "1", "01", "2", "02"},
	} {
		for _, context := range []string{"alone", "day", "year"} {
			t.Run(tc.locale+tc.width+context, func(t *testing.T) {
				t.Parallel()
				opts := Options{Month: option.String(tc.width), TimeZone: option.String("UTC")}
				want, second := tc.alone, tc.second
				if context == "day" {
					opts.Day = option.String("numeric")
					want, second = tc.format, tc.secondFormat
				}
				if context == "year" {
					opts.Year = option.String("numeric")
					if tc.width == "numeric" {
						want, second = "01", "02"
					}
				}
				f, err := New(intltest.LocaleList(t, tc.locale), opts)
				if err != nil {
					t.Fatal(err)
				}
				parts := mustDateFormatToParts(t, f, start)
				var text strings.Builder
				var months []string
				for _, part := range parts {
					text.WriteString(part.Value)
					if part.Type == PartMonth {
						months = append(months, part.Value)
					}
				}
				if !slices.Equal(months, []string{want}) {
					t.Errorf("month = %q, want %q", months, want)
				}
				if got := mustDateFormat(t, f, start); got != text.String() {
					t.Errorf("text = %q, parts = %q", got, text.String())
				}
				text.Reset()
				months = nil
				for _, part := range mustDateFormatRangeToParts(t, f, start, end) {
					text.WriteString(part.Value)
					if part.Type == PartMonth {
						months = append(months, part.Value)
					}
				}
				if !slices.Equal(months, []string{want, second}) {
					t.Errorf("range months = %q, want %q", months, []string{want, second})
				}
				if got := mustDateFormatRange(t, f, start, end); got != text.String() {
					t.Errorf("range text = %q, parts = %q", got, text.String())
				}
			})
		}
	}
}
