package conformance

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"testing"
)

func TestFixtureJSONPreservesEmptyExpectations(t *testing.T) {
	t.Parallel()

	f := Fixture{
		ID: "empty", Source: "manual:presence", Locale: "en",
		Options: jsontext.Value(`{}`), Input: jsontext.Value(`1`),
		Expected: new(""), ExpectedOK: new(false), ExpectedLocales: []string{},
		ExpectedParts: []Part{}, ExpectedRange: new(""), ExpectedRangeParts: []RangePart{},
		ExpectedResolved: jsontext.Value(`{}`),
	}
	data, err := json.Marshal([]Fixture{f})
	if err != nil {
		t.Fatal(err)
	}
	var records []map[string]jsontext.Value
	if err := json.Unmarshal(data, &records); err != nil {
		t.Fatal(err)
	}
	for field, want := range map[string]string{
		"expected": `""`, "expectedOk": `false`, "expectedLocales": `[]`,
		"expectedParts": `[]`, "expectedRange": `""`, "expectedRangeParts": `[]`,
		"expectedResolvedOptions": `{}`,
	} {
		if got := string(records[0][field]); got != want {
			t.Errorf("encoded %s = %q, want %q", field, got, want)
		}
	}
	if _, err := decodeFixtures("empty.json", data); err != nil {
		t.Fatalf("decodeFixtures() rejected encoded observations: %v", err)
	}
}

func TestFixtureJSONOmitsAbsentExpectations(t *testing.T) {
	t.Parallel()

	f := Fixture{
		ID: "text", Source: "manual:presence", Locale: "en",
		Options: jsontext.Value(`{}`), Input: jsontext.Value(`1`), Expected: new("1"),
	}
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	var record map[string]jsontext.Value
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"expectedOk", "expectedLocales", "expectedParts", "expectedRange", "expectedRangeParts", "expectedResolvedOptions"} {
		if value, ok := record[field]; ok {
			t.Errorf("absent %s encoded as %s", field, value)
		}
	}
}
