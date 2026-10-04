package relativetimeformat

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/agentable/go-intl/internal/testcontract"
	"github.com/agentable/go-intl/tools/conformance"
)

func TestRelativeConformanceObservations(t *testing.T) {
	if testcontract.FixtureRunnerChild(t, runRelativeConformanceFixture) {
		return
	}
	t.Parallel()

	parts := []conformance.Part{
		{Type: "literal", Value: "in "},
		{Type: "integer", Value: "1", Unit: "day"},
		{Type: "literal", Value: " day"},
	}
	literal := []conformance.Part{{Type: "literal", Value: "tomorrow"}}
	tests := []struct {
		name    string
		fixture conformance.Fixture
		failure string
	}{
		{"numeric parts only", conformance.Fixture{ExpectedParts: parts}, ""},
		{"literal parts only", conformance.Fixture{Options: jsontext.Value(`{"numeric":"auto"}`), ExpectedParts: literal}, ""},
		{"numeric observations", conformance.Fixture{Expected: new("in 1 day"), ExpectedParts: parts}, ""},
		{"literal observations", conformance.Fixture{Options: jsontext.Value(`{"numeric":"auto"}`), Expected: new("tomorrow"), ExpectedParts: literal}, ""},
		{"wrong text", conformance.Fixture{Expected: new("wrong"), ExpectedParts: parts}, "Format("},
		{"wrong number parts", conformance.Fixture{Expected: new("in 1 day"), ExpectedParts: []conformance.Part{{Type: "integer", Value: "2", Unit: "day"}}}, "FormatToParts"},
		{"wrong literal parts", conformance.Fixture{Options: jsontext.Value(`{"numeric":"auto"}`), Expected: new("tomorrow"), ExpectedParts: []conformance.Part{{Type: "literal", Value: "yesterday"}}}, "FormatToParts"},
		{"empty expected parts", conformance.Fixture{Expected: new("in 1 day"), ExpectedParts: []conformance.Part{}}, "FormatToParts"},
		{"negative zero", conformance.Fixture{Input: jsontext.Value(`{"value":-0,"unit":"day"}`), ExpectedParts: []conformance.Part{{Type: "integer", Value: "0", Unit: "day"}, {Type: "literal", Value: " days ago"}}}, ""},
		{"positive zero", conformance.Fixture{Input: jsontext.Value(`{"value":0,"unit":"day"}`), ExpectedParts: []conformance.Part{{Type: "literal", Value: "in "}, {Type: "integer", Value: "0", Unit: "day"}, {Type: "literal", Value: " days"}}}, ""},
		{"supported empty", conformance.Fixture{Feature: "supportedLocalesOf", Input: jsontext.Value(`[]`), ExpectedLocales: []string{}}, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := tc.fixture
			f.ID = "observation"
			f.Source = "manual:adapter"
			f.Locale = "en"
			if f.Options == nil {
				f.Options = jsontext.Value(`{}`)
			}
			if f.Input == nil {
				f.Input = jsontext.Value(`{"value":1,"unit":"day"}`)
			}
			testcontract.AssertFixtureRunner(t, "relativetimeformat", "TestRelativeConformanceObservations", f, tc.failure)
		})
	}
}
