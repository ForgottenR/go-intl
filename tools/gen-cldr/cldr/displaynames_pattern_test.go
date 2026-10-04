package cldr

import (
	"encoding/json/v2"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDisplayNamesPatternValidation(t *testing.T) {
	t.Parallel()
	for _, field := range []string{"localePattern", "localeSeparator"} {
		for _, value := range []string{"{0}", "{1}", "{0}/{2}", "{0}/{1", "{0}/{{1}", "{0} ({literal}{1})"} {
			t.Run(field+"/"+value, func(t *testing.T) {
				t.Parallel()
				root := t.TempDir()
				mustWriteDisplayNamesFixture(t, root, "en")
				path := changeDisplayPattern(t, root, field, value)
				_, err := loadDisplayNames(root, []string{"en"})
				if err == nil {
					t.Fatal("loader accepted damaged pattern")
				}
				for _, context := range []string{path, "en", field} {
					if !strings.Contains(err.Error(), context) {
						t.Errorf("error lacks %q: %v", context, err)
					}
				}
			})
		}
	}
}

func TestDisplayNamesPatternDefaultsAndRepeatedArguments(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	mustWriteDisplayNamesFixture(t, root, "en")
	changeDisplayPattern(t, root, "localePattern", "{0}（{1}、{1}）")
	changeDisplayPattern(t, root, "localeSeparator", "{1}／{0}／{1}")
	data, err := loadDisplayNames(root, []string{"en", "en-US"})
	if err != nil {
		t.Fatal(err)
	}
	if data["en-US"].LocalePattern != "{0}（{1}、{1}）" || data["en-US"].LocaleSeparator != "{1}／{0}／{1}" {
		t.Errorf("inherited patterns = %#v", data["en-US"])
	}
	changeDisplayPattern(t, root, "localeSeparator", "")
	data, err = loadDisplayNames(root, []string{"en"})
	if err != nil || data["en"].LocaleSeparator != "" {
		t.Errorf("missing separator default = %#v, %v", data, err)
	}
}

func changeDisplayPattern(t *testing.T, root, field, value string) string {
	t.Helper()
	path := filepath.Join(root, "cldr-localenames-full", "main", "en", "localeDisplayNames.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	patterns := doc["main"].(map[string]any)["en"].(map[string]any)["localeDisplayNames"].(map[string]any)["localeDisplayPattern"].(map[string]any)
	if value == "" {
		delete(patterns, field)
	} else {
		patterns[field] = value
	}
	raw, err = json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
