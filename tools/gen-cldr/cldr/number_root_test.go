package cldr

import (
	"bytes"
	"strings"
	"testing"
)

func TestNumberRootRejectsInvalidInput(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, version, context string
		pin, xml               []byte
	}{
		{"version", "49.0.0", "VERSION cldr", numberRootPin, numberRootXML},
		{"SHA", "48.1.0", "SHA-256", numberRootPin, append(bytes.Clone(numberRootXML), ' ')},
		{"pin JSON", "48.1.0", "source.json", []byte(`{`), numberRootXML},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := parseNumberRoot(tc.version, tc.pin, tc.xml)
			if err == nil || !strings.Contains(err.Error(), tc.context) {
				t.Fatalf("parseNumberRoot() error = %v, want %s", err, tc.context)
			}
		})
	}
}

func TestRootSymbolsRejectsInvalidData(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, old, replacement, context string }{
		{"missing group", `<group>٬</group>`, ``, "missing group"},
		{"alias source", `source="locale" path="../symbols[@numberSystem='latn']"`, `source="root" path="../symbols[@numberSystem='latn']"`, "unsupported alias"},
		{"alias path", `../symbols[@numberSystem='latn']`, `../symbols[@numberSystem='arab']`, "unsupported alias"},
		{"duplicate system", `numberSystem="arabext"`, `numberSystem="arab"`, "duplicate"},
		{"missing system", `numberSystem="adlm"`, `numberSystem="unknown"`, "symbols[adlm]"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			xml := strings.Replace(string(numberRootXML), tc.old, tc.replacement, 1)
			if xml == string(numberRootXML) {
				t.Fatalf("root source missing %q", tc.old)
			}
			_, err := parseRootSymbols([]byte(xml))
			if err == nil || !strings.Contains(err.Error(), tc.context) {
				t.Fatalf("parseRootSymbols() error = %v, want %s", err, tc.context)
			}
		})
	}
}
