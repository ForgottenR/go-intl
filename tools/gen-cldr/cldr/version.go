// Package cldr ingests the unicode-org/cldr-json npm distribution into
// typed Go structs for the generator extract layer to consume. It is
// generator-only; the runtime package shares the name but lives at
// internal/cldr/ in the main module.
package cldr

import (
	"encoding/json/v2"
	"fmt"
	"path/filepath"
	"strings"
)

// Versions captures the CLDR / ICU / tzdata pin from internal/cldr/VERSION.
type Versions struct {
	CLDR   string
	ICU    string
	TZData string
}

// ReadVersionFile parses the three-line pin file. Unknown lines are ignored;
// any of the required keys (cldr/icu/tzdata) being absent is an error so the
// generator never runs with a partially-specified pin.
func ReadVersionFile(path string) (Versions, error) {
	raw, err := readRequiredFile(path)
	if err != nil {
		return Versions{}, err
	}
	var v Versions
	for line := range strings.SplitSeq(string(raw), "\n") {
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch strings.TrimSpace(key) {
		case "cldr":
			v.CLDR = strings.TrimSpace(val)
		case "icu":
			v.ICU = strings.TrimSpace(val)
		case "tzdata":
			v.TZData = strings.TrimSpace(val)
		}
	}
	if v.CLDR == "" || v.ICU == "" || v.TZData == "" {
		return Versions{}, fmt.Errorf("incomplete pin in %s: %+v", path, v)
	}
	return v, nil
}

// CrossCheck verifies every CLDR package consumed by the generator against
// the pinned package identity before any output is written.
func CrossCheck(cldrJSONRoot string, want Versions) error {
	for _, name := range requiredPackages {
		pkgPath := filepath.Join(cldrJSONRoot, name, "package.json")
		raw, err := readRequiredFile(pkgPath)
		if err != nil {
			return fmt.Errorf("read %s: %w", pkgPath, err)
		}
		var meta struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}
		if err := json.Unmarshal(raw, &meta); err != nil {
			return fmt.Errorf("parse %s: %w", pkgPath, err)
		}
		if meta.Name != name || meta.Version != want.CLDR {
			return fmt.Errorf("CLDR package %s: expected name %q version %q, got name %q version %q", pkgPath, name, want.CLDR, meta.Name, meta.Version)
		}
	}
	return nil
}
