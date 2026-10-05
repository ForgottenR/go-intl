package cldr

import (
	"crypto/sha256"
	_ "embed"
	"encoding/json/v2"
	"encoding/xml"
	"fmt"

	"github.com/agentable/go-intl/internal/numbering"
)

// These generator-only inputs fill the symbols and aliases omitted from npm
// JSON. The original XML bytes retain upstream provenance and licensing.
//
//go:embed number-root/root.xml
var numberRootXML []byte

//go:embed number-root/source.json
var numberRootPin []byte

type numberRootData struct {
	symbols map[string]NumberSymbols
	hashes  map[string]string
}

func loadNumberRoot(version string) (numberRootData, error) {
	return parseNumberRoot(version, numberRootPin, numberRootXML)
}

func parseNumberRoot(version string, pinBytes, xmlBytes []byte) (numberRootData, error) {
	const pinName = "tools/gen-cldr/cldr/number-root/source.json"
	const xmlName = "tools/gen-cldr/cldr/number-root/root.xml"
	var pin struct {
		CLDR   string `json:"cldr"`
		SHA256 string `json:"sha256"`
	}
	if err := json.Unmarshal(pinBytes, &pin); err != nil {
		return numberRootData{}, fmt.Errorf("parse %s: %w", pinName, err)
	}
	if pin.CLDR != version {
		return numberRootData{}, fmt.Errorf("%s: cldr %q, want VERSION cldr %q", pinName, pin.CLDR, version)
	}
	digest := fmt.Sprintf("%x", sha256.Sum256(xmlBytes))
	if digest != pin.SHA256 {
		return numberRootData{}, fmt.Errorf("%s: SHA-256 %s, want %s from %s", xmlName, digest, pin.SHA256, pinName)
	}
	symbols, err := parseRootSymbols(xmlBytes)
	if err != nil {
		return numberRootData{}, fmt.Errorf("parse %s: %w", xmlName, err)
	}
	return numberRootData{
		symbols: symbols,
		hashes: map[string]string{
			pinName: fmt.Sprintf("%x", sha256.Sum256(pinBytes)),
			xmlName: digest,
		},
	}, nil
}

func parseRootSymbols(data []byte) (map[string]NumberSymbols, error) {
	var doc struct {
		XMLName xml.Name `xml:"ldml"`
		Rows    []struct {
			System string `xml:"numberSystem,attr"`
			Alias  *struct {
				Source string `xml:"source,attr"`
				Path   string `xml:"path,attr"`
			} `xml:"alias"`
			Fields []struct {
				XMLName xml.Name
				Value   string `xml:",chardata"`
				Alt     string `xml:"alt,attr"`
			} `xml:",any"`
		} `xml:"numbers>symbols"`
	}
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	out := make(map[string]NumberSymbols)
	seen := make(map[string]bool)
	for _, row := range doc.Rows {
		if row.System == "" || seen[row.System] {
			return nil, fmt.Errorf("numbers/symbols: missing or duplicate numberSystem %q", row.System)
		}
		seen[row.System] = true
		if row.Alias != nil {
			if row.Alias.Source != "locale" || row.Alias.Path != "../symbols[@numberSystem='latn']" || len(row.Fields) != 0 || row.System == "latn" {
				return nil, fmt.Errorf("numbers/symbols[%s]: unsupported alias source=%q path=%q", row.System, row.Alias.Source, row.Alias.Path)
			}
			continue
		}
		fields := make(map[string]string)
		for _, field := range row.Fields {
			if field.Alt == "" {
				fields[field.XMLName.Local] = field.Value
			}
		}
		for _, field := range []string{"decimal", "group", "percentSign", "plusSign", "minusSign", "nan", "infinity", "approximatelySign", "perMille", "exponential", "superscriptingExponent", "timeSeparator"} {
			if fields[field] == "" {
				return nil, fmt.Errorf("numbers/symbols[%s]: missing %s", row.System, field)
			}
		}
		// Reuse the JSON symbol projection; root XML contributes symbol facts only.
		raw, err := json.Marshal(fields)
		if err != nil {
			return nil, err
		}
		symbols, err := parseNumberSymbols(raw)
		if err != nil {
			return nil, fmt.Errorf("numbers/symbols[%s]: %w", row.System, err)
		}
		out[row.System] = symbols
	}
	for _, system := range numbering.SimpleNumberingSystems() {
		if !seen[system] {
			return nil, fmt.Errorf("numbers/symbols[%s]: no concrete row or locale-latn alias", system)
		}
	}
	if _, ok := out["latn"]; !ok {
		return nil, fmt.Errorf("numbers/symbols[latn]: concrete row required")
	}
	return out, nil
}

func (root numberRootData) fill(numbers map[string]Numbers) {
	for _, data := range numbers {
		for system, symbols := range root.symbols {
			if _, ok := data.Symbols[system]; !ok {
				symbols.RangeSign = data.Symbols[data.DefaultNumberingSystem].RangeSign
				data.Symbols[system] = symbols
			}
		}
	}
}
