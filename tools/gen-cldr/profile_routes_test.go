package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunRejectsUnbackedProfileBeforeWriting(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, profile, remove, domain string }{
		{"no source", `{"locales":["xx-YY"]}`, "", "number"},
		{"missing parent domain", `{"locales":["en-US"]}`, "cldr-misc-full/main/en/listPatterns.json", "list"},
		{"missing domain", `{"locales":["en"]}`, "cldr-misc-full/main/en/listPatterns.json", "list"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			root := filepath.Join(dir, "node_modules")
			writeRuntimeCLDRFixtures(t, root)
			writeListPatternCLDRFixture(t, root)
			writeRelativeTimeCLDRFixture(t, root)
			writeDisplayNamesCLDRFixture(t, root)
			if tc.remove != "" {
				if err := os.Remove(filepath.Join(root, tc.remove)); err != nil {
					t.Fatal(err)
				}
			}
			version := filepath.Join(dir, "VERSION")
			mustWriteGenCLDRFile(t, version, "cldr=48.1.0\nicu=78\ntzdata=2025b\n")
			profile := filepath.Join(dir, "profile.json")
			mustWriteGenCLDRFile(t, profile, tc.profile)
			out := filepath.Join(dir, "out")
			sentinel := filepath.Join(out, "existing.go")
			mustWriteGenCLDRFile(t, sentinel, "unchanged")
			err := Run(context.Background(), Config{CLDRDir: root, OutDir: out, VersionFile: version, ProfileFile: profile, TZDataLock: writeTZDataFixture(t)}, slog.New(slog.NewTextHandler(io.Discard, nil)))
			if err == nil {
				t.Fatal("Run accepted an unbacked profile")
			}
			for _, part := range []string{tc.domain, root} {
				if !strings.Contains(err.Error(), part) {
					t.Errorf("error %q missing %q", err, part)
				}
			}
			entries, err := os.ReadDir(out)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 1 {
				t.Fatalf("failure wrote output: %v", entries)
			}
			got, err := os.ReadFile(sentinel)
			if err != nil || string(got) != "unchanged" {
				t.Fatalf("existing output = %q, %v", got, err)
			}
		})
	}
}

func TestRunLoadsProfileLookupParents(t *testing.T) {
	t.Parallel()
	for _, locale := range []string{"en-US", "zh-Hans-CN"} {
		t.Run(locale, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			root := filepath.Join(dir, "node_modules")
			writeRuntimeCLDRFixtures(t, root)
			if locale == "zh-Hans-CN" {
				for _, pkg := range []string{"cldr-numbers-full", "cldr-dates-full", "cldr-misc-full", "cldr-units-full", "cldr-localenames-full"} {
					parent := filepath.Join(root, pkg, "main", "en")
					entries, err := os.ReadDir(parent)
					if err != nil {
						t.Fatal(err)
					}
					for _, entry := range entries {
						data, err := os.ReadFile(filepath.Join(parent, entry.Name()))
						if err != nil {
							t.Fatal(err)
						}
						mustWriteGenCLDRFile(t, filepath.Join(root, pkg, "main", "zh-Hans", entry.Name()), strings.ReplaceAll(string(data), `"en"`, `"zh-Hans"`))
					}
				}
				mustWriteGenCLDRFile(t, filepath.Join(root, "cldr-core", "availableLocales.json"), `{"availableLocales":{"modern":["zh-Hans"]}}`)
			}
			version := filepath.Join(dir, "VERSION")
			mustWriteGenCLDRFile(t, version, "cldr=48.1.0\nicu=78\ntzdata=2025b\n")
			profile := filepath.Join(dir, "profile.json")
			mustWriteGenCLDRFile(t, profile, `{"locales":["`+locale+`"]}`)
			out := filepath.Join(dir, "out")
			if err := Run(context.Background(), Config{CLDRDir: root, OutDir: out, VersionFile: version, ProfileFile: profile, TZDataLock: writeTZDataFixture(t)}, slog.New(slog.NewTextHandler(io.Discard, nil))); err != nil {
				t.Fatal(err)
			}
			for _, domain := range []string{"number", "date", "currency", "timezone", "unit", "list", "relativetime", "displaynames", "locale"} {
				if _, err := os.Stat(filepath.Join(out, domain, "data.go")); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
