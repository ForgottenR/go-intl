package datetimeformat

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/agentable/go-intl/internal/intltest"
)

func TestFractionSeparator(t *testing.T) {
	t.Parallel()
	for _, loc := range []string{"fr-FR", "ar-EG"} {
		for _, systems := range [][]string{{"latn", "arab"}, {"arab", "latn"}} {
			for _, nu := range systems {
				for digits := 1; digits <= 3; digits++ {
					f, err := New(intltest.LocaleList(t, loc), Options{TimeZone: new("UTC"), Second: new("numeric"), FractionalSecondDigits: new(digits), NumberingSystem: new(nu)})
					if err != nil {
						t.Fatal(err)
					}
					t.Run(fmt.Sprintf("%s/%s/%d", loc, nu, digits), func(t *testing.T) {
						t.Parallel()
						assertFractionMethods(t, f, loc, nu, digits)
						fresh, err := New(intltest.LocaleList(t, loc), Options{TimeZone: new("UTC"), Second: new("numeric"), FractionalSecondDigits: new(digits), NumberingSystem: new(nu)})
						if err != nil {
							t.Fatal(err)
						}
						assertFractionMethods(t, fresh, loc, nu, digits)
					})
				}
			}
		}
	}
}

func assertFractionMethods(t *testing.T, f *DateTimeFormat, loc, nu string, digits int) {
	t.Helper()
	separator := "."
	if loc == "fr-FR" {
		separator = ","
	}
	wantDigits := "987"[:digits]
	if nu == "arab" {
		separator = "٫"
		wantDigits = string([]rune("٩٨٧")[:digits])
	}
	instant := time.Date(2024, 1, 1, 0, 0, 7, 987_999_999, time.UTC)
	parts, err := f.FormatToParts(instant)
	if err != nil {
		t.Fatal(err)
	}
	joined := ""
	fractions := 0
	for i, p := range parts {
		joined += p.Value
		if p.Type != PartFractionalSecond {
			continue
		}
		fractions++
		if p.Value != wantDigits {
			t.Errorf("fraction %q, want %q", p.Value, wantDigits)
		}
		if i == 0 || parts[i-1].Type != PartLiteral || parts[i-1].Value != separator {
			t.Errorf("separator want %q: %v", separator, parts)
		}
	}
	if fractions != 1 {
		t.Errorf("fraction count %d: %v", fractions, parts)
	}
	if text, err := f.Format(instant); err != nil || text != joined {
		t.Errorf("Format %q, %v, parts %q", text, err, joined)
	}
	for _, delta := range []time.Duration{3 * time.Second, 24*time.Hour + 3*time.Second} {
		parts, err := f.FormatRangeToParts(instant, instant.Add(delta))
		if err != nil {
			t.Fatal(err)
		}
		joined = ""
		fractions := 0
		for i, p := range parts {
			joined += p.Value
			if p.Type != PartFractionalSecond {
				continue
			}
			fractions++
			if p.Value != wantDigits {
				t.Errorf("range fraction %q, want %q", p.Value, wantDigits)
			}
			if i == 0 || parts[i-1].Type != PartLiteral || !strings.HasSuffix(parts[i-1].Value, separator) {
				t.Errorf("range separator want %q: %v", separator, parts)
			}
		}
		if fractions != 2 {
			t.Errorf("fraction count %d: %v", fractions, parts)
		}
		if text, err := f.FormatRange(instant, instant.Add(delta)); err != nil || text != joined {
			t.Errorf("range %q, %v, parts %q", text, err, joined)
		}
	}
}
