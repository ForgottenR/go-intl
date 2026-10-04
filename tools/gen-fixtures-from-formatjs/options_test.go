package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunFormatJSOptionsPreserveSourceSemantics(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		options string
		want    string
	}{
		{name: "variable", options: "{minimumFractionDigits: digits}"},
		{name: "arithmetic", options: "{minimumFractionDigits: 1 + 1}"},
		{name: "spread", options: "{minimumFractionDigits: 2, ...other}"},
		{name: "computed key", options: "{['minimumFractionDigits']: 2}"},
		{name: "shorthand", options: "{minimumFractionDigits}"},
		{name: "nested value", options: "{minimumFractionDigits: {value: 2}}"},
		{name: "duplicate number", options: "{minimumFractionDigits: 1, minimumFractionDigits: 2}", want: `{"minimumFractionDigits":2}`},
		{name: "duplicate types", options: "{useGrouping: 'auto', useGrouping: false}", want: `{"useGrouping":false}`},
		{name: "comments", options: "{ /* minimumFractionDigits: 9 */ minimumFractionDigits: 2, // trailing comment\n}", want: `{"minimumFractionDigits":2}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			source := t.TempDir()
			out := t.TempDir()
			for _, name := range []string{
				formatJSNumberFormatPackageDir, formatJSPluralRulesPackageDir,
				formatJSDateTimeFormatPackageDir, formatJSLocalePackageDir,
				formatJSCanonicalLocalesPackageDir, formatJSListFormatPackageDir,
				formatJSRelativeTimeFormatPackageDir, formatJSDurationFormatPackageDir,
			} {
				if err := os.MkdirAll(formatJSPackageTestsRoot(source, name), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			mustWriteFile(t, filepath.Join(formatJSPackageTestsRoot(source, formatJSNumberFormatPackageDir), "options.test.ts"),
				"const nf = new NumberFormat('en', "+tc.options+");\nexpect(nf.format(1)).toBe('1.00');\n")
			if err := run([]string{"-formatjs", source, "-out", out}); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(formatJSConformanceRoot(out, "numberformat"), "options-test-ts.json")
			if tc.want == "" {
				assertPathAbsent(t, "ununderstood options fixture", path)
				assertFileContainsAll(t, "source-owned skip", filepath.Join(out, repositorySkipListFile),
					formatJSNumberFormatTestSourcePrefix+"options.test.ts", skipCategoryUnsupportedExtractorShape)
				return
			}
			fixtures := mustReadFixtures(t, path)
			if len(fixtures) != 1 {
				t.Fatalf("fixtures = %d, want 1", len(fixtures))
			}
			assertJSONEqual(t, "exported options", fixtures[0].Options, tc.want)
			assertJSONEqual(t, "skip list", mustReadSkips(t, filepath.Join(out, repositorySkipListFile)), "[]")
		})
	}
}

func TestUnunderstoodConstructorDoesNotReuseEarlierOptions(t *testing.T) {
	t.Parallel()
	for _, options := range []string{"{minimumFractionDigits: digits}", "options", "{minimumFractionDigits: 2} as any"} {
		t.Run(options, func(t *testing.T) {
			t.Parallel()
			source := `const nf = new NumberFormat('en', {minimumFractionDigits: 2});
expect(nf.format(1)).toBe('1.00');
const nf = new NumberFormat('en', ` + options + `);
expect(nf.format(2)).toBe('2.00');`
			fixtures := extractNumberFormatFixtures("options.test.ts", source)
			if len(fixtures) != 1 {
				t.Fatalf("extracted %d fixtures, want only the first understood declaration", len(fixtures))
			}
			assertJSONEqual(t, "understood constructor", fixtures[0].Options, `{"minimumFractionDigits":2}`)
		})
	}
}

func TestRunUnsupportedLocaleConstructorBlocksEarlierDeclaration(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ constructor, packageDir, outputPackage, options, assertion string }{
		{"NumberFormat", formatJSNumberFormatPackageDir, "numberformat", "{}", "expect(f.format(1)).toBe('1');"},
		{"PluralRules", formatJSPluralRulesPackageDir, "pluralrules", "{}", "expect(f.select(1)).toBe('one');"},
		{"DateTimeFormat", formatJSDateTimeFormatPackageDir, "datetimeformat", "{year: 'numeric', month: 'short', day: 'numeric', timeZone: 'UTC'}", "expect(f.format(instant)).toBe('Jan 1, 2020');"},
	} {
		for _, localeExpression := range []string{"localeName", "['fr', 'en']"} {
			t.Run(tc.constructor+"/"+localeExpression, func(t *testing.T) {
				t.Parallel()
				source, out := t.TempDir(), t.TempDir()
				for _, packageDir := range []string{
					formatJSNumberFormatPackageDir, formatJSPluralRulesPackageDir,
					formatJSDateTimeFormatPackageDir, formatJSLocalePackageDir,
					formatJSCanonicalLocalesPackageDir, formatJSListFormatPackageDir,
					formatJSRelativeTimeFormatPackageDir, formatJSDurationFormatPackageDir,
				} {
					if err := os.MkdirAll(formatJSPackageTestsRoot(source, packageDir), 0o755); err != nil {
						t.Fatal(err)
					}
				}
				data := "const instant = new Date('2020-01-01T00:00:00Z');\n" +
					"it('supported', () => { const f = new " + tc.constructor + "('en', " + tc.options + "); " + tc.assertion + " });\n" +
					"it('unsupported', () => { const localeName = 'fr'; const f = new " + tc.constructor + "(" + localeExpression + ", " + tc.options + "); " + tc.assertion + " });\n"
				mustWriteFile(t, filepath.Join(formatJSPackageTestsRoot(source, tc.packageDir), "locale.test.ts"), data)
				if err := run([]string{"-formatjs", source, "-out", out}); err != nil {
					t.Fatal(err)
				}
				fixtures := mustReadFixtures(t, filepath.Join(formatJSConformanceRoot(out, tc.outputPackage), "locale-test-ts.json"))
				if len(fixtures) != 1 || fixtures[0].Locale != "en" {
					t.Fatalf("fixtures = %#v, want only the first understood constructor", fixtures)
				}
				assertFileContainsAll(t, "source-owned extraction audit", filepath.Join(out, repositorySkipListFile), formatJSTestSourcePrefix(tc.packageDir)+"locale.test.ts", skipCategoryPartialExtraction)
			})
		}
	}
}
