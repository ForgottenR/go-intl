package datetimeformat

import (
	"strings"
	"testing"
	"time"

	"github.com/agentable/go-intl/internal/intltest"
)

func TestOffsetNameMinuteWidth(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ locale, zone, form, want string }{
		{"en-US", "+05:03", "shortOffset", "GMT+5:03"}, {"en-US", "-05:05", "shortOffset", "GMT-5:05"},
		{"en-US", "+05:00", "shortOffset", "GMT+5"}, {"en-US", "+00:00", "shortOffset", "GMT"},
		{"en-US", "+00:00", "longOffset", "GMT+00:00"}, {"en-US", "-05:03", "longOffset", "GMT-05:03"},
		{"fi", "+05:03", "shortOffset", "UTC+5.03"}, {"fi", "-05:05", "shortOffset", "UTC-5.05"},
		{"fi", "+05:00", "shortOffset", "UTC+5"}, {"fi", "-05:00", "longOffset", "UTC-05.00"},
	} {
		t.Run(tc.locale+tc.zone+tc.form, func(t *testing.T) {
			t.Parallel()
			f, err := New(intltest.LocaleList(t, tc.locale), Options{TimeZone: new(tc.zone), TimeZoneName: new(tc.form), Hour: new("numeric"), Minute: new("2-digit")})
			if err != nil {
				t.Fatal(err)
			}
			assertZoneNameMethods(t, f, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), tc.want)
		})
	}
}

func assertZoneNameMethods(t *testing.T, f *DateTimeFormat, instant time.Time, want string) {
	t.Helper()
	parts, err := f.FormatToParts(instant)
	if err != nil {
		t.Fatal(err)
	}
	joined := ""
	found := false
	for _, p := range parts {
		joined += p.Value
		if p.Type == PartTimeZoneName {
			found = true
			if p.Value != want {
				t.Errorf("timeZoneName = %q, want %q", p.Value, want)
			}
		}
	}
	if !found {
		t.Fatal("no zone part")
	}
	text, err := f.Format(instant)
	if err != nil || text != joined {
		t.Fatalf("Format = %q, %v, parts %q", text, err, joined)
	}
	end := instant.Add(3 * time.Hour)
	rangeParts, err := f.FormatRangeToParts(instant, end)
	if err != nil {
		t.Fatal(err)
	}
	joined = ""
	found = false
	for _, p := range rangeParts {
		joined += p.Value
		if p.Type == PartTimeZoneName {
			found = true
			if p.Value != want {
				t.Errorf("range timeZoneName = %q, want %q", p.Value, want)
			}
		}
	}
	if !found {
		t.Fatalf("no range zone part: %v", rangeParts)
	}
	text, err = f.FormatRange(instant, end)
	if err != nil || text != joined || !strings.Contains(text, want) {
		t.Errorf("FormatRange = %q, %v, parts %q", text, err, joined)
	}
}

func TestOffsetNameNumberingSystem(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ tag, nu, zone, form, want string }{
		{"en-US-u-nu-arab", "", "+05:30", "longOffset", "GMT+٠٥:٣٠"},
		{"en-US", "arab", "-05:03", "shortOffset", "GMT-٥:٠٣"},
		{"en-US", "arab", "+05:03", "shortOffset", "GMT+٥:٠٣"},
		{"en-US", "arab", "-05:30", "longOffset", "GMT-٠٥:٣٠"},
		{"en-US-u-nu-arab", "latn", "+05:30", "longOffset", "GMT+05:30"},
		{"en-US", "arab", "Etc/GMT-5", "long", "GMT+٠٥:٠٠"},
		{"en-US", "arab", "America/New_York", "long", "Eastern Standard Time"},
	} {
		t.Run(tc.tag+tc.nu+tc.zone+tc.form, func(t *testing.T) {
			t.Parallel()
			opts := Options{TimeZone: new(tc.zone), TimeZoneName: new(tc.form), Hour: new("numeric"), Minute: new("2-digit")}
			if tc.nu != "" {
				opts.NumberingSystem = new(tc.nu)
			}
			f, err := New(intltest.LocaleList(t, tc.tag), opts)
			if err != nil {
				t.Fatal(err)
			}
			assertZoneNameMethods(t, f, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), tc.want)
		})
	}
}

func TestHistoricalOffsetSeconds(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ form, nu, want string }{{"longOffset", "latn", "GMT+00:09:21"}, {"shortOffset", "latn", "GMT+0:09:21"}, {"longOffset", "arab", "GMT+٠٠:٠٩:٢١"}} {
		t.Run(tc.form+tc.nu, func(t *testing.T) {
			t.Parallel()
			f, err := New(intltest.LocaleList(t, "en-US"), Options{TimeZone: new("Europe/Paris"), TimeZoneName: new(tc.form), NumberingSystem: new(tc.nu), Hour: new("numeric"), Minute: new("2-digit")})
			if err != nil {
				t.Fatal(err)
			}
			assertZoneNameMethods(t, f, time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), tc.want)
		})
	}
}
