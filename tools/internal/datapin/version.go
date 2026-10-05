// Package datapin parses generation input pins shared by maintainer tools.
package datapin

import (
	"fmt"
	"os"
	"strings"
)

// Versions records source, reference, and identity versions; environment checks remain with each tool.
type Versions struct{ CLDR, ICU, TZData string }

// ReadVersions rejects malformed or duplicate lines and requires all three pins.
// Comments and blank lines are ignored; well-formed unknown keys are allowed.
func ReadVersions(path string) (Versions, error) {
	//nolint:gosec // G304: maintainer tools explicitly select the pin file.
	data, err := os.ReadFile(path)
	if err != nil {
		return Versions{}, fmt.Errorf("read %s: %w", path, err)
	}
	values := map[string]string{}
	for lineNumber, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
			return Versions{}, fmt.Errorf("%s:%d: malformed version pin %q", path, lineNumber+1, line)
		}
		key = strings.TrimSpace(key)
		if _, exists := values[key]; exists {
			return Versions{}, fmt.Errorf("%s:%d: duplicate version pin %q", path, lineNumber+1, key)
		}
		values[key] = strings.TrimSpace(value)
	}
	for _, key := range []string{"cldr", "icu", "tzdata"} {
		if values[key] == "" {
			return Versions{}, fmt.Errorf("%s: missing %s version pin", path, key)
		}
	}
	if !isDottedNumericVersion(values["cldr"], 3) {
		return Versions{}, fmt.Errorf("%s: invalid cldr version pin cldr=%s", path, values["cldr"])
	}
	return Versions{CLDR: values["cldr"], ICU: values["icu"], TZData: values["tzdata"]}, nil
}

func isDottedNumericVersion(version string, parts int) bool {
	components := strings.Split(version, ".")
	if len(components) != parts {
		return false
	}
	for _, component := range components {
		if component == "" {
			return false
		}
		for _, r := range component {
			if r < '0' || r > '9' {
				return false
			}
		}
	}
	return true
}
