package datetimeformat

import (
	"encoding/json"
	"encoding/json/jsontext"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/agentable/go-intl/internal/intltest"
	"github.com/agentable/go-intl/internal/testcontract"
)

func TestStyleHourCycles(t *testing.T) {
	t.Parallel()
	for _, loc := range []string{"en-US", "en-GB", "ja-JP", "fr-CA"} {
		for _, cycle := range []string{"h11", "h12", "h23", "h24"} {
			for _, style := range []string{"short", "medium"} {
				for _, dateStyle := range []bool{false, true} {
					t.Run(fmt.Sprintf("%s/%s/%s/%t", loc, cycle, style, dateStyle), func(t *testing.T) {
						t.Parallel()
						opts := Options{TimeZone: new("UTC"), TimeStyle: new(style), HourCycle: new(cycle)}
						if dateStyle {
							opts.DateStyle = new("short")
						}
						f, err := New(intltest.LocaleList(t, loc), opts)
						if err != nil {
							t.Fatal(err)
						}
						record, err := json.Marshal(f.ResolvedOptions())
						if err != nil {
							t.Fatal(err)
						}
						var values map[string]any
						if err := json.Unmarshal(record, &values); err != nil {
							t.Fatal(err)
						}
						for _, key := range []string{"hour", "minute", "second", "dayPeriod"} {
							if _, ok := values[key]; ok {
								t.Errorf("style reports component %s", key)
							}
						}
						if values["hourCycle"] != cycle {
							t.Errorf("record = %s", record)
						}
						midnight := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
						assertHourCycle(t, f, cycle, midnight)
						assertHourCycle(t, f, cycle, midnight.Add(15*time.Hour))
						assertCycleRange(t, f, cycle, midnight, midnight.Add(3*time.Hour))
						assertCycleRange(t, f, cycle, midnight, midnight.Add(27*time.Hour))
					})
				}
			}
		}
	}
}

func assertCycleRange(t *testing.T, f *DateTimeFormat, cycle string, start, end time.Time) {
	t.Helper()
	parts, err := f.FormatRangeToParts(start, end)
	if err != nil {
		t.Fatal(err)
	}
	joined := ""
	var hours, periods int
	for _, p := range parts {
		joined += p.Value
		if p.Source != SourceStartRange && p.Source != SourceEndRange && p.Source != SourceShared {
			t.Errorf("invalid source %q", p.Source)
		}
		if p.Type == PartDayPeriod {
			periods++
		}
		if p.Type != PartHour {
			continue
		}
		hours++
		instant := start
		if p.Source == SourceEndRange {
			instant = end
		}
		want := instant.Hour()
		if cycle == "h11" || cycle == "h12" {
			want %= 12
		}
		if cycle == "h12" && want == 0 {
			want = 12
		}
		if cycle == "h24" && want == 0 {
			want = 24
		}
		if got, err := strconv.Atoi(p.Value); err != nil || got != want {
			t.Errorf("range hour = %q (%s), want %d: %v", p.Value, p.Source, want, parts)
		}
		if p.Source == SourceShared {
			t.Errorf("different hours marked shared: %v", parts)
		}
	}
	if hours != 2 {
		t.Errorf("range hour count %d: %v", hours, parts)
	}
	if want12 := cycle == "h11" || cycle == "h12"; (periods > 0) != want12 {
		t.Errorf("range day periods = %d for %s: %v", periods, cycle, parts)
	}
	if text, err := f.FormatRange(start, end); err != nil || text != joined {
		t.Errorf("range %q, %v differs from parts %q", text, err, joined)
	}
}

func TestExplicitHourCycles(t *testing.T) {
	t.Parallel()
	for _, loc := range []string{"en-US", "en-GB", "fr-CA"} {
		for _, cycle := range []string{"h11", "h12", "h23", "h24"} {
			for _, matcher := range []string{"basic", "best fit"} {
				for _, extension := range []bool{false, true} {
					t.Run(fmt.Sprintf("%s/%s/%s/%t", loc, cycle, matcher, extension), func(t *testing.T) {
						t.Parallel()
						opts := Options{TimeZone: new("UTC"), Hour: new("numeric"), FormatMatcher: new(matcher)}
						tag := loc
						if extension {
							tag += "-u-hc-" + cycle
						} else {
							opts.HourCycle = new(cycle)
						}
						f, err := New(intltest.LocaleList(t, tag), opts)
						if err != nil {
							t.Fatal(err)
						}
						resolved := f.ResolvedOptions()
						if resolved.HourCycle == nil || string(*resolved.HourCycle) != cycle {
							t.Errorf("hourCycle = %v, want %s", resolved.HourCycle, cycle)
						}
						if extension && !strings.Contains(resolved.Locale.String(), "hc-"+cycle) {
							t.Errorf("locale = %s, want hc retained", resolved.Locale)
						}
						assertHourCycle(t, f, cycle, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
					})
				}
			}
		}
	}
	for _, optionCycle := range []string{"h11", "h24"} {
		f, err := New(intltest.LocaleList(t, "en-US-u-hc-h23"), Options{TimeZone: new("UTC"), Hour: new("numeric"), HourCycle: new(optionCycle)})
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(f.ResolvedOptions().Locale.String(), "hc-") {
			t.Errorf("overridden hc retained: %s", f.ResolvedOptions().Locale)
		}
		assertHourCycle(t, f, optionCycle, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	}
}

func TestLocaleHour12Preference(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		locale string
		hour12 bool
		cycle  string
	}{
		{"fr-CA", true, "h12"}, {"fr-CA", false, "h23"}, {"ja", true, "h11"}, {"ja-JP", true, "h11"}, {"en-US", true, "h12"},
		{"ja", false, "h23"}, {"en-US", false, "h23"}, {"fr-FR", false, "h23"},
	} {
		for _, style := range []bool{false, true} {
			for _, extension := range []string{"", "-u-hc-h11", "-u-hc-h12", "-u-hc-h23", "-u-hc-h24"} {
				t.Run(fmt.Sprintf("%s/%t/%t/%s", tc.locale, tc.hour12, style, extension), func(t *testing.T) {
					t.Parallel()
					opts := Options{TimeZone: new("UTC"), Hour12: new(tc.hour12), HourCycle: new("h24")}
					if style {
						opts.TimeStyle = new("short")
					} else {
						opts.Hour = new("numeric")
						opts.Minute = new("2-digit")
					}
					f, err := New(intltest.LocaleList(t, tc.locale+extension), opts)
					if err != nil {
						t.Fatal(err)
					}
					r := f.ResolvedOptions()
					if r.HourCycle == nil || string(*r.HourCycle) != tc.cycle || r.Hour12 == nil || *r.Hour12 != tc.hour12 {
						t.Errorf("resolved = %#v, want %s/%t", r, tc.cycle, tc.hour12)
					}
					if strings.Contains(r.Locale.String(), "hc-") {
						t.Errorf("hour12 retained overridden hc: %s", r.Locale)
					}
					if tc.locale == "ja" && r.Locale.String() != "ja" {
						t.Errorf("maximization leaked into locale: %s", r.Locale)
					}
					midnight := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
					assertHourCycle(t, f, tc.cycle, midnight)
					assertCycleRange(t, f, tc.cycle, midnight, midnight.Add(3*time.Hour))
				})
			}
		}
	}
}

func assertHourCycle(t *testing.T, f *DateTimeFormat, cycle string, instant time.Time) {
	t.Helper()
	want := instant.Hour()
	if cycle == "h11" || cycle == "h12" {
		want %= 12
	}
	if cycle == "h12" && want == 0 {
		want = 12
	}
	if cycle == "h24" && want == 0 {
		want = 24
	}
	parts, err := f.FormatToParts(instant)
	if err != nil {
		t.Fatal(err)
	}
	var joined string
	var hours, periods int
	for _, p := range parts {
		joined += p.Value
		if p.Type == PartDayPeriod {
			periods++
		}
		if p.Type == PartHour {
			hours++
			if got, err := strconv.Atoi(p.Value); err != nil || got != want {
				t.Errorf("hour = %q, want %d (%s)", p.Value, want, cycle)
			}
			if width := f.ResolvedOptions().Hour; width != nil && ((*width == TwoDigitFieldStyle) != (len(p.Value) == 2)) && want < 10 {
				t.Errorf("hour width %q disagrees with %s", p.Value, *width)
			}
		}
	}
	if hours != 1 {
		t.Errorf("hours = %d, parts %v", hours, parts)
	}
	if want12 := cycle == "h11" || cycle == "h12"; (periods > 0) != want12 {
		t.Errorf("day periods = %d for %s: %v", periods, cycle, parts)
	}
	if text, err := f.Format(instant); err != nil || text != joined {
		t.Errorf("text %q, %v differs from parts %q", text, err, joined)
	}
}

func TestLanguageRegionDefaultHourCycle(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, tag, cycle, resolved string
		options                    Options
	}{
		{"French components", "fr-CA", "h23", `{"locale":"fr-CA","calendar":"gregory","numberingSystem":"latn","timeZone":"UTC","hourCycle":"h23","hour12":false,"hour":"2-digit"}`, Options{Hour: new("numeric")}},
		{"English components", "en-CA", "h12", `{"locale":"en-CA","calendar":"gregory","numberingSystem":"latn","timeZone":"UTC","hourCycle":"h12","hour12":true,"hour":"numeric"}`, Options{Hour: new("numeric")}},
		{"French style", "fr-CA", "h23", `{"locale":"fr-CA","calendar":"gregory","numberingSystem":"latn","timeZone":"UTC","hourCycle":"h23","hour12":false,"timeStyle":"short"}`, Options{TimeStyle: new("short")}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			opts := tc.options
			opts.TimeZone = new("UTC")
			f, err := New(intltest.LocaleList(t, tc.tag), opts)
			if err != nil {
				t.Fatal(err)
			}
			testcontract.AssertResolvedOptionsJSON(t, f.ResolvedOptions(), jsontext.Value(tc.resolved))
			midnight := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			assertHourCycle(t, f, tc.cycle, midnight)
			assertCycleRange(t, f, tc.cycle, midnight, midnight.Add(3*time.Hour))
			assertCycleRange(t, f, tc.cycle, midnight, midnight.Add(27*time.Hour))
		})
	}
}
