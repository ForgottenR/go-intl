package main

import (
	"io"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunRejectsInvalidFixtureShapes(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		fields string
		want   string
	}{
		{"numeric options", `"options":42,"expected":""`, "options"},
		{"array options", `"options":[],"expected":""`, "options"},
		{"null options", `"options":null,"expected":""`, "options"},
		{"unknown field", `"options":{},"expected":"","unknownField":true`, "unknownField"},
		{"no observation", `"options":{}`, "expectation"},
		{"null text", `"options":{},"expected":null,"expectedParts":[]`, "expected"},
		{"null parts", `"options":{},"expected":"","expectedParts":null`, "expectedParts"},
		{"null range parts", `"options":{},"expected":"","expectedRangeParts":null`, "expectedRangeParts"},
		{"null locales", `"options":{},"expected":"","expectedLocales":null`, "expectedLocales"},
		{"null resolved", `"options":{},"expected":"","expectedResolvedOptions":null`, "expectedResolvedOptions"},
		{"array resolved", `"options":{},"expected":"","expectedResolvedOptions":[]`, "expectedResolvedOptions"},
		{"null ok", `"options":{},"expected":"","expectedOk":null`, "expectedOk"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			packageDir := writeFixtureFile(t, root, "listformat", "listformat/testdata/conformance/manual/shape.json",
				`[{"id":"shape-probe","source":"manual:probe","locale":"en","input":[],`+tc.fields+`}]`)
			err := run([]string{packageDir}, io.Discard)
			if err == nil {
				t.Fatalf("run() accepted %s", tc.name)
			}
			if !strings.Contains(err.Error(), filepath.Join("manual", "shape.json")) || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("run() error = %v, want file path and %q", err, tc.want)
			}
		})
	}
}

func TestRunPreservesEmptyAndIndependentExpectations(t *testing.T) {
	t.Parallel()
	for _, expectation := range []string{
		`"expected":""`, `"expectedOk":false`, `"expectedParts":[]`,
		`"expectedRangeParts":[]`, `"expectedLocales":[]`, `"expectedResolvedOptions":{}`,
	} {
		t.Run(expectation, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			packageDir := writeFixtureFile(t, root, "listformat", "listformat/testdata/conformance/manual/shape.json",
				`[{"id":"shape-probe","source":"manual:probe","locale":"en","options":{},"input":[],`+expectation+`}]`)
			if err := run([]string{packageDir}, io.Discard); err != nil {
				t.Fatalf("run() rejected present expectation: %v", err)
			}
		})
	}
}
