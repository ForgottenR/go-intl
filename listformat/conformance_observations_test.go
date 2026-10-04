package listformat

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/agentable/go-intl/internal/testcontract"
	"github.com/agentable/go-intl/tools/conformance"
)

func TestListConformanceObservations(t *testing.T) {
	if testcontract.FixtureRunnerChild(t, runListConformanceFixture) {
		return
	}
	t.Parallel()

	parts := []conformance.Part{
		{Type: "element", Value: "A"},
		{Type: "literal", Value: " and "},
		{Type: "element", Value: "B"},
	}
	tests := []struct {
		name    string
		fixture conformance.Fixture
		failure string
	}{
		{"parts only", conformance.Fixture{ExpectedParts: parts}, ""},
		{"empty parts only", conformance.Fixture{Input: jsontext.Value(`[]`), ExpectedParts: []conformance.Part{}}, ""},
		{"text and parts", conformance.Fixture{Expected: new("A and B"), ExpectedParts: parts}, ""},
		{"empty text and parts", conformance.Fixture{Input: jsontext.Value(`[]`), Expected: new(""), ExpectedParts: []conformance.Part{}}, ""},
		{"wrong text", conformance.Fixture{Expected: new("wrong"), ExpectedParts: parts}, "Format("},
		{"wrong element", conformance.Fixture{Expected: new("A and B"), ExpectedParts: []conformance.Part{{Type: "element", Value: "wrong"}}}, "FormatToParts"},
		{"wrong literal", conformance.Fixture{Expected: new("A and B"), ExpectedParts: []conformance.Part{{Type: "literal", Value: "A and B"}}}, "FormatToParts"},
		{"empty expected parts", conformance.Fixture{Expected: new("A and B"), ExpectedParts: []conformance.Part{}}, "FormatToParts"},
		{"supported empty", conformance.Fixture{Feature: "supportedLocalesOf", Input: jsontext.Value(`[]`), ExpectedLocales: []string{}}, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := tc.fixture
			f.ID = "observation"
			f.Source = "manual:adapter"
			f.Locale = "en"
			f.Options = jsontext.Value(`{}`)
			if f.Input == nil {
				f.Input = jsontext.Value(`["A","B"]`)
			}
			testcontract.AssertFixtureRunner(t, "listformat", "TestListConformanceObservations", f, tc.failure)
		})
	}
}
