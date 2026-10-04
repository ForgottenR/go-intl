package numberformat

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"testing"

	"github.com/agentable/go-intl/internal/testcontract"
	"github.com/agentable/go-intl/tools/conformance"
)

func TestNumberResolvedSnapshots(t *testing.T) {
	if testcontract.FixtureRunnerChild(t, runNumberConformanceFixture) {
		return
	}
	t.Parallel()

	defaults := `{"locale":"en","numberingSystem":"latn","style":"decimal","minimumIntegerDigits":1,"minimumFractionDigits":0,"maximumFractionDigits":3,"useGrouping":"auto","notation":"standard","signDisplay":"auto","roundingIncrement":1,"roundingMode":"halfExpand","roundingPriority":"auto","trailingZeroDisplay":"auto"}`
	tests := []struct {
		name    string
		options string
		remove  string
		key     string
		value   any
		failure string
	}{
		{name: "complete", options: "{}"},
		{name: "missing required key", options: "{}", remove: "roundingMode", failure: "ResolvedOptions"},
		{name: "extra currency", options: "{}", key: "currency", value: "USD", failure: "ResolvedOptions"},
		{name: "unknown key", options: "{}", key: "unknown", value: true, failure: "ResolvedOptions"},
		{name: "disabled boolean", options: `{"useGrouping":false}`, key: "useGrouping", value: false},
		{name: "disabled string is wrong", options: `{"useGrouping":false}`, key: "useGrouping", value: "false", failure: "ResolvedOptions"},
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
			f := conformance.Fixture{
				ID: "snapshot", Source: "manual:adapter", Locale: "en",
				Options: jsontext.Value(tc.options), Input: jsontext.Value(`{}`),
				ExpectedResolved: want,
			}
			testcontract.AssertFixtureRunner(t, "numberformat", "TestNumberResolvedSnapshots", f, tc.failure)
		})
	}
}
