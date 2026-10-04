package displaynames

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/agentable/go-intl/internal/testcontract"
	"github.com/agentable/go-intl/tools/conformance"
)

func TestDisplayNamesResolvedSnapshots(t *testing.T) {
	if testcontract.FixtureRunnerChild(t, runDisplayNamesConformanceFixture) {
		return
	}
	t.Parallel()

	// ECMA-402 DisplayNames resolvedOptions table; Node 26.10.0 witness.
	tests := []struct {
		name    string
		want    string
		failure string
	}{
		{"complete", `{"fallback":"code","type":"region","style":"long","locale":"en"}`, ""},
		{"missing key", `{"locale":"en","style":"long","type":"region"}`, "ResolvedOptions"},
		{"extra language display", `{"locale":"en","style":"long","type":"region","fallback":"code","languageDisplay":"dialect"}`, "ResolvedOptions"},
		{"unknown key", `{"locale":"en","style":"long","type":"region","fallback":"code","unexpected":true}`, "ResolvedOptions"},
		{"null versus absent", `{"locale":"en","style":"long","type":"region","fallback":"code","languageDisplay":null}`, "ResolvedOptions"},
		{"wrong type", `{"locale":"en","style":"long","type":"region","fallback":false}`, "ResolvedOptions"},
		{"wrong value", `{"locale":"en","style":"long","type":"region","fallback":"none"}`, "ResolvedOptions"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := conformance.Fixture{
				ID: "snapshot", Source: "manual:adapter", Locale: "en",
				Options: jsontext.Value(`{"type":"region"}`),
				Input:   jsontext.Value(`"US"`), Expected: new("United States"),
				ExpectedResolved: jsontext.Value(tc.want),
			}
			testcontract.AssertFixtureRunner(t, "displaynames", "TestDisplayNamesResolvedSnapshots", f, tc.failure)
		})
	}
}
