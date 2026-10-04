package durationformat

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"testing"

	"github.com/agentable/go-intl/internal/testcontract"
	"github.com/agentable/go-intl/tools/conformance"
)

func TestDurationResolvedSnapshots(t *testing.T) {
	t.Parallel()

	// Node 26.10.0; ECMA-402 DurationFormat resolvedOptions table.
	defaults := `{"locale":"en","numberingSystem":"latn","style":"long","years":"long","yearsDisplay":"auto","months":"long","monthsDisplay":"auto","weeks":"long","weeksDisplay":"auto","days":"long","daysDisplay":"auto","hours":"long","hoursDisplay":"auto","minutes":"long","minutesDisplay":"auto","seconds":"long","secondsDisplay":"auto","milliseconds":"long","millisecondsDisplay":"auto","microseconds":"long","microsecondsDisplay":"auto","nanoseconds":"long","nanosecondsDisplay":"auto"}`
	tests := []struct {
		name    string
		options string
		remove  string
		key     string
		value   any
		failure string
	}{
		{name: "complete"},
		{name: "missing style", remove: "style", failure: "ResolvedOptions"},
		{name: "unknown key", key: "unknown", value: true, failure: "ResolvedOptions"},
		{name: "null versus absent", key: "fractionalDigits", value: nil, failure: "ResolvedOptions"},
		{name: "fraction zero", options: `{"style":"long","fractionalDigits":0}`, key: "fractionalDigits", value: 0},
		{name: "unreported fraction zero", options: `{"style":"long","fractionalDigits":0}`, failure: "ResolvedOptions"},
		{name: "wrong type", key: "secondsDisplay", value: false, failure: "ResolvedOptions"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var record map[string]any
			if err := json.Unmarshal([]byte(defaults), &record); err != nil {
				t.Fatal(err)
			}
			delete(record, tc.remove)
			if tc.key != "" {
				record[tc.key] = tc.value
			}
			want, err := json.Marshal(record)
			if err != nil {
				t.Fatal(err)
			}
			f := durationObservationFixture(conformance.Fixture{ExpectedResolved: want})
			if tc.options != "" {
				f.Options = jsontext.Value(tc.options)
			}
			testcontract.AssertFixtureRunner(t, "durationformat", "TestDurationFormatConformance", f, tc.failure)
		})
	}
}
