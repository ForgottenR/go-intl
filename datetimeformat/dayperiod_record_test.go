package datetimeformat

import (
	"encoding/json/v2"
	"slices"
	"testing"
	"time"

	"github.com/agentable/go-intl/internal/intltest"
)

func TestDayPeriodRecordAndParts(t *testing.T) {
	t.Parallel()

	morning := time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC)
	afternoon := morning.Add(12 * time.Hour)
	tests := []struct {
		name  string
		opts  Options
		width *string
		parts []Part
	}{
		{"automatic hour", Options{Hour: new("numeric")}, nil, []Part{{Type: PartHour, Value: "1"}, {Type: PartLiteral, Value: " "}, {Type: PartDayPeriod, Value: "AM"}}},
		{"automatic style", Options{TimeStyle: new("short")}, nil, []Part{{Type: PartHour, Value: "1"}, {Type: PartLiteral, Value: ":"}, {Type: PartMinute, Value: "00"}, {Type: PartLiteral, Value: " "}, {Type: PartDayPeriod, Value: "AM"}}},
		{"flexible long", Options{Hour: new("numeric"), DayPeriod: new("long")}, new("long"), []Part{{Type: PartHour, Value: "1"}, {Type: PartLiteral, Value: " "}, {Type: PartDayPeriod, Value: "in the morning"}}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			opts := tc.opts
			opts.TimeZone = new("UTC")
			f, err := New(intltest.LocaleList(t, "en-US"), opts)
			if err != nil {
				t.Fatal(err)
			}
			data, err := json.Marshal(f.ResolvedOptions())
			if err != nil {
				t.Fatal(err)
			}
			var record map[string]any
			if err := json.Unmarshal(data, &record); err != nil {
				t.Fatal(err)
			}
			gotWidth, present := record["dayPeriod"]
			if tc.width == nil && present {
				t.Fatalf("automatic AM/PM leaked into resolved dayPeriod: %s", data)
			}
			if tc.width != nil && (!present || gotWidth != *tc.width) {
				t.Fatalf("dayPeriod = %v, present=%v, want %q", gotWidth, present, *tc.width)
			}
			parts := mustDateFormatToParts(t, f, morning)
			if !slices.Equal(parts, tc.parts) {
				t.Fatalf("FormatToParts = %#v, want %#v", parts, tc.parts)
			}
			var joined string
			for _, part := range parts {
				joined += part.Value
			}
			if got := mustDateFormat(t, f, morning); got != joined {
				t.Fatalf("Format = %q, parts join = %q", got, joined)
			}
			rangeParts := mustDateFormatRangeToParts(t, f, morning, afternoon)
			var rangeText string
			var periods []RangePart
			for _, part := range rangeParts {
				rangeText += part.Value
				if part.Type == PartDayPeriod {
					periods = append(periods, part)
				}
			}
			wantPeriods := []RangePart{{Type: PartDayPeriod, Value: "AM", Source: SourceStartRange}, {Type: PartDayPeriod, Value: "PM", Source: SourceEndRange}}
			if tc.width != nil {
				wantPeriods[0].Value = "in the morning"
				wantPeriods[1].Value = "in the afternoon"
			}
			if !slices.Equal(periods, wantPeriods) {
				t.Fatalf("range dayPeriod parts = %#v, want %#v", periods, wantPeriods)
			}
			if got := mustDateFormatRange(t, f, morning, afternoon); got != rangeText {
				t.Fatalf("FormatRange = %q, parts join = %q", got, rangeText)
			}
		})
	}
}
