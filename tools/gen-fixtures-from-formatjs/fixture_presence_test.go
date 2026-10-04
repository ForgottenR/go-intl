package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"os"
	"path/filepath"
	"testing"

	"github.com/agentable/go-intl/tools/conformance"
)

func TestWriteFixturesPreservesEmptyExpectations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		field string
		value string
		set   func(*fixture)
	}{
		{"expected", `""`, func(f *fixture) { f.Expected = new("") }},
		{"expectedOk", `false`, func(f *fixture) { f.ExpectedOK = new(false) }},
		{"expectedLocales", `[]`, func(f *fixture) { f.ExpectedLocales = []string{} }},
		{"expectedParts", `[]`, func(f *fixture) { f.ExpectedParts = []conformance.Part{} }},
		{"expectedRange", `""`, func(f *fixture) { f.ExpectedRange = new("") }},
		{"expectedRangeParts", `[]`, func(f *fixture) { f.ExpectedRangeParts = []conformance.RangePart{} }},
		{"expectedResolvedOptions", `{}`, func(f *fixture) { f.ExpectedResolved = map[string]any{} }},
	}
	for _, tc := range tests {
		t.Run(tc.field, func(t *testing.T) {
			t.Parallel()

			root := filepath.Join(t.TempDir(), "numberformat")
			path := filepath.Join(root, "testdata", "conformance", "manual", "empty.json")
			f := fixture{ID: "empty", Source: "manual:presence", Locale: "en", Options: map[string]any{}, Input: 1}
			tc.set(&f)
			if err := writeJSON(path, []fixture{f}); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var records []map[string]jsontext.Value
			if err := json.Unmarshal(data, &records); err != nil {
				t.Fatal(err)
			}
			if got := string(records[0][tc.field]); got != tc.value {
				t.Fatalf("written %s = %q, want %q", tc.field, got, tc.value)
			}
			loaded, err := conformance.LoadFixtures(root)
			if err != nil {
				t.Fatalf("LoadFixtures() rejected written observation: %v", err)
			}
			if len(loaded) != 1 {
				t.Fatalf("LoadFixtures() returned %d fixtures, want 1", len(loaded))
			}
		})
	}
}
