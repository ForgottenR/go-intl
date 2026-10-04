package durationformat

import (
	"encoding/json/jsontext"
	"strings"
	"testing"

	"github.com/agentable/go-intl/internal/testcontract"
	"github.com/agentable/go-intl/tools/conformance"
)

func TestDurationConformanceObservations(t *testing.T) {
	t.Parallel()

	resolved := jsontext.Value(`{"locale":"en","numberingSystem":"latn","style":"long","years":"long","yearsDisplay":"auto","months":"long","monthsDisplay":"auto","weeks":"long","weeksDisplay":"auto","days":"long","daysDisplay":"auto","hours":"long","hoursDisplay":"auto","minutes":"long","minutesDisplay":"auto","seconds":"long","secondsDisplay":"auto","milliseconds":"long","millisecondsDisplay":"auto","microseconds":"long","microsecondsDisplay":"auto","nanoseconds":"long","nanosecondsDisplay":"auto"}`)
	parts := []conformance.Part{{Type: "integer", Value: "1", Unit: "second"}, {Type: "literal", Value: " ", Unit: "second"}, {Type: "unit", Value: "second", Unit: "second"}}
	tests := []struct {
		name    string
		fixture conformance.Fixture
		failure string
	}{
		{"resolved only", conformance.Fixture{Input: jsontext.Value(`"not a duration"`), ExpectedResolved: resolved}, ""},
		{"parts only", conformance.Fixture{ExpectedParts: parts}, ""},
		{"empty parts only", conformance.Fixture{Input: jsontext.Value(`{"seconds":0}`), ExpectedParts: []conformance.Part{}}, ""},
		{"observations", conformance.Fixture{Expected: new("1 second"), ExpectedParts: parts, ExpectedResolved: resolved}, ""},
		{"wrong text", conformance.Fixture{Expected: new("wrong"), ExpectedParts: parts}, "Format("},
		{"wrong parts", conformance.Fixture{Expected: new("1 second"), ExpectedParts: []conformance.Part{{Type: "unit", Value: "minute", Unit: "minute"}}}, "FormatToParts"},
		{"empty expected parts", conformance.Fixture{Expected: new("1 second"), ExpectedParts: []conformance.Part{}}, "FormatToParts"},
		{"wrong resolved", conformance.Fixture{Expected: new("1 second"), ExpectedResolved: jsontext.Value(`{"locale":"fr"}`)}, "ResolvedOptions"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := durationObservationFixture(tc.fixture)
			testcontract.AssertFixtureRunner(t, "durationformat", "TestDurationFormatConformance", f, tc.failure)
		})
	}
}

func TestDurationConformanceSuiteLedgers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		file    string
		content string
		failure string
	}{
		{"malformed divergence", "divergences.md", "malformed divergence line\n", "malformed divergence line"},
		{"malformed xfail", "xfail.json", "{", "xfail.json"},
		{"expired xfail", "xfail.json", `[{"id":"observation","reason":"pending fix","expires_at":"2000-01-01","tracking_issue":"SPEC-70"}]`, "expired"},
		{"valid xfail", "xfail.json", `[{"id":"observation","reason":"pending fix","expires_at":"2999-01-01","tracking_issue":"SPEC-70"}]`, ""},
		{"valid divergence", "divergences.md", "id: observation\nsource: manual:adapter\nowner: durationformat\nstatus: accepted\nreason: reference output differs\nreview_after: 2999-01-01\nremoval_path: refresh reference\n", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := durationObservationFixture(conformance.Fixture{Expected: new("wrong")})
			if tc.failure != "" {
				// A callback would report this option, making premature execution visible.
				f.Options = jsontext.Value(`{"callback-marker":"must not execute"}`)
			}
			output := testcontract.AssertFixtureSuite(t, "durationformat", "TestDurationFormatConformance", f, map[string]string{tc.file: tc.content}, tc.failure)
			if strings.Contains(output, "unsupported option") {
				t.Fatalf("fixture callback ran before ledger rejection:\n%s", output)
			}
		})
	}
}

func durationObservationFixture(f conformance.Fixture) conformance.Fixture {
	f.ID = "observation"
	f.Source = "manual:adapter"
	f.Locale = "en"
	f.Options = jsontext.Value(`{"style":"long"}`)
	if f.Input == nil {
		f.Input = jsontext.Value(`{"seconds":1}`)
	}
	return f
}
