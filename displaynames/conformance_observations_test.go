package displaynames

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/agentable/go-intl/internal/testcontract"
	"github.com/agentable/go-intl/tools/conformance"
)

func TestDisplayNamesConformanceObservations(t *testing.T) {
	if testcontract.FixtureRunnerChild(t, runDisplayNamesConformanceFixture) {
		return
	}
	t.Parallel()

	resolved := jsontext.Value(`{"locale":"en","style":"long","type":"language","fallback":"code","languageDisplay":"dialect"}`)
	tests := []struct {
		name    string
		fixture conformance.Fixture
		failure string
	}{
		{"resolved only", conformance.Fixture{Input: jsontext.Value(`{}`), ExpectedResolved: resolved}, ""},
		{"ok only", conformance.Fixture{ExpectedOK: new(true)}, ""},
		{"missing ok only", conformance.Fixture{Options: jsontext.Value(`{"type":"language","fallback":"none"}`), Input: jsontext.Value(`"zz"`), ExpectedOK: new(false)}, ""},
		{"text only", conformance.Fixture{Expected: new("English")}, ""},
		{"missing text only", conformance.Fixture{Options: jsontext.Value(`{"type":"language","fallback":"none"}`), Input: jsontext.Value(`"zz"`), Expected: new("")}, ""},
		{"observations", conformance.Fixture{Expected: new("English"), ExpectedOK: new(true), ExpectedResolved: resolved}, ""},
		{"wrong ok", conformance.Fixture{Expected: new("English"), ExpectedOK: new(false)}, "ok ="},
		{"wrong text", conformance.Fixture{Expected: new("wrong"), ExpectedOK: new(true)}, "Of("},
		{"wrong resolved", conformance.Fixture{Expected: new("English"), ExpectedResolved: jsontext.Value(`{"locale":"fr"}`)}, "ResolvedOptions"},
		{"constructor error", conformance.Fixture{Options: jsontext.Value(`{"type":""}`), ErrorCode: "invalidOption"}, ""},
		{"of error", conformance.Fixture{Input: jsontext.Value(`"not@language"`), ErrorCode: "invalidCode"}, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := tc.fixture
			f.ID = "observation"
			f.Source = "manual:adapter"
			f.Locale = "en"
			if f.Options == nil {
				f.Options = jsontext.Value(`{"type":"language"}`)
			}
			if f.Input == nil {
				f.Input = jsontext.Value(`"en"`)
			}
			testcontract.AssertFixtureRunner(t, "displaynames", "TestDisplayNamesConformanceObservations", f, tc.failure)
		})
	}
}
