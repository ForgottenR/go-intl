package datetimeformat

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"testing"

	"github.com/agentable/go-intl/internal/testcontract"
	"github.com/agentable/go-intl/tools/conformance"
)

func TestDateTimeResolvedSnapshots(t *testing.T) {
	if testcontract.FixtureRunnerChild(t, runDateTimeConformanceFixture) {
		return
	}
	t.Parallel()

	defaults := `{"locale":"en-US","calendar":"gregory","numberingSystem":"latn","timeZone":"UTC","hourCycle":"h12","hour12":true,"hour":"numeric"}`
	tests := []struct {
		name    string
		options string
		remove  string
		key     string
		value   any
		failure string
	}{
		{name: "complete hour", options: `{"hour":"numeric"}`},
		{name: "missing hour cycle", options: `{"hour":"numeric"}`, remove: "hourCycle", failure: "ResolvedOptions"},
		{name: "missing hour12", options: `{"hour":"numeric"}`, remove: "hour12", failure: "ResolvedOptions"},
		{name: "unreported flexible day period", options: `{"hour":"numeric","dayPeriod":"long"}`, failure: "ResolvedOptions"},
		{name: "null versus absent", options: `{"hour":"numeric"}`, key: "dayPeriod", value: nil, failure: "ResolvedOptions"},
		{name: "unknown key", options: `{"hour":"numeric"}`, key: "extra", value: true, failure: "ResolvedOptions"},
		{name: "style complete", options: `{"timeStyle":"short"}`, remove: "hour", key: "timeStyle", value: "short"},
		{name: "style missing", options: `{"timeStyle":"short"}`, remove: "hour", failure: "ResolvedOptions"},
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
			var options map[string]any
			if err := json.Unmarshal([]byte(tc.options), &options); err != nil {
				t.Fatal(err)
			}
			options["timeZone"] = "UTC"
			encodedOptions, err := json.Marshal(options)
			if err != nil {
				t.Fatal(err)
			}
			f := conformance.Fixture{
				ID: "snapshot", Source: "manual:adapter", Locale: "en-US",
				Options: encodedOptions, Input: jsontext.Value(`"2020-01-01T01:00:00Z"`),
				ExpectedResolved: want,
			}
			testcontract.AssertFixtureRunner(t, "datetimeformat", "TestDateTimeResolvedSnapshots", f, tc.failure)
		})
	}
}
