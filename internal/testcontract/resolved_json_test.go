package testcontract

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/agentable/go-intl/tools/conformance"
)

func TestResolvedOptionsJSONScalars(t *testing.T) {
	if FixtureRunnerChild(t, func(t *testing.T, f conformance.Fixture) {
		record := struct {
			Zero  int    `json:"zero"`
			False bool   `json:"false"`
			Empty string `json:"empty"`
		}{}
		AssertResolvedOptionsJSON(t, record, f.ExpectedResolved)
	}) {
		return
	}
	t.Parallel()

	tests := []struct {
		name    string
		want    string
		failure string
	}{
		{"preserve zero values and ignore order", `{"empty":"","false":false,"zero":0}`, ""},
		{"missing zero", `{"empty":"","false":false}`, "ResolvedOptions"},
		{"missing false", `{"empty":"","zero":0}`, "ResolvedOptions"},
		{"missing empty", `{"false":false,"zero":0}`, "ResolvedOptions"},
		{"number versus string", `{"empty":"","false":false,"zero":"0"}`, "ResolvedOptions"},
		{"boolean versus string", `{"empty":"","false":"false","zero":0}`, "ResolvedOptions"},
		{"string versus null", `{"empty":null,"false":false,"zero":0}`, "ResolvedOptions"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := conformance.Fixture{
				ID: "scalars", Source: "manual:adapter", Locale: "en",
				Options: jsontext.Value(`{}`), Input: jsontext.Value(`{}`),
				ExpectedResolved: jsontext.Value(tc.want),
			}
			AssertFixtureRunner(t, "testcontract", "TestResolvedOptionsJSONScalars", f, tc.failure)
		})
	}
}
