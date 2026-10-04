package main

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunGeneratesTimezoneDomain asserts the generator emits the timezone domain
// as a self-contained const-only payload at internal/cldr/timezone/data.go, no
// longer emits the retired root metazones.go literal renderer or the timezone
// alias facade or identity payload; identifier truth is generated separately
// into internal/tz.
func TestRunGeneratesTimezoneDomain(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := filepath.Join(dir, "node_modules")
	writeRuntimeCLDRFixtures(t, root)
	out := filepath.Join(dir, "out")
	if err := os.MkdirAll(out, 0o777); err != nil {
		t.Fatalf("mkdir out: %v", err)
	}
	versionPath := filepath.Join(out, "VERSION")
	if err := os.WriteFile(versionPath, []byte("cldr=48.1.0\nicu=78\ntzdata=2025b\n"), 0o666); err != nil {
		t.Fatalf("write VERSION: %v", err)
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	if err := Run(context.Background(), Config{CLDRDir: root, OutDir: out, VersionFile: versionPath, ProfileFile: writeLocaleProfileFixture(t, dir), TZDataLock: writeTZDataFixture(t)}, log); err != nil {
		t.Fatalf("Run: %v", err)
	}

	// The retired root literal renderer must not produce metazones.go anymore.
	if _, err := os.Stat(filepath.Join(out, "metazones.go")); !os.IsNotExist(err) {
		t.Fatalf("root metazones.go should no longer be generated, stat err = %v", err)
	}
	// The alias facade timezone/names.go is retired in favour of hand-written
	// accessors.go; the generator must not emit it.
	if _, err := os.Stat(filepath.Join(out, "timezone", "names.go")); !os.IsNotExist(err) {
		t.Fatalf("timezone/names.go should no longer be generated, stat err = %v", err)
	}

	// The timezone domain payload is a const-only data.go carrying the blobs and
	// the private _data table the decoder reads.
	payload, err := os.ReadFile(filepath.Join(out, "timezone", "data.go"))
	if err != nil {
		t.Fatalf("read timezone/data.go: %v", err)
	}
	if !containsAll(string(payload), "package timezone", "const _data",
		"_tzMetazonePeriodBlob", "_tzNamesBlob", "_tzFormatsBlob") || containsAll(string(payload), "_tzSupportedBlob") {
		t.Fatalf("timezone/data.go missing expected const payload:\n%s", payload)
	}
	joined := dewrapStringLiterals(string(payload))
	if !containsAll(joined, "America_Pacific", "Eastern Time", "British Summer Time", "GMT{0}") {
		t.Fatalf("timezone/data.go _data missing expected time-zone strings:\n%s", payload)
	}

	// The retired root supported.go literal file must no longer be generated.
	if _, err := os.Stat(filepath.Join(out, "supported.go")); !os.IsNotExist(err) {
		t.Fatalf("root supported.go should no longer be generated, stat err = %v", err)
	}
}

func TestRunAcceptsSingleMetazoneObject(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	root := filepath.Join(dir, "node_modules")
	writeRuntimeCLDRFixtures(t, root)
	supp := filepath.Join(root, "cldr-core", "supplemental")
	metaZones := `{"supplemental":{"metaZones":{"metazoneInfo":{"timezone":{"America":{"Los_Angeles":{"usesMetazone":{"_mzone":"America_Pacific"}}}}}}}}`
	if err := os.WriteFile(filepath.Join(supp, "metaZones.json"), []byte(metaZones), 0o666); err != nil {
		t.Fatalf("write metaZones: %v", err)
	}
	out := filepath.Join(dir, "out")
	if err := os.MkdirAll(out, 0o777); err != nil {
		t.Fatalf("mkdir out: %v", err)
	}
	versionPath := filepath.Join(out, "VERSION")
	if err := os.WriteFile(versionPath, []byte("cldr=48.1.0\nicu=78\ntzdata=2025b\n"), 0o666); err != nil {
		t.Fatalf("write VERSION: %v", err)
	}

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	if err := Run(context.Background(), Config{CLDRDir: root, OutDir: out, VersionFile: versionPath, ProfileFile: writeLocaleProfileFixture(t, dir), TZDataLock: writeTZDataFixture(t)}, log); err != nil {
		t.Fatalf("Run: %v", err)
	}
}

func writeTimeZoneCLDRFixture(t *testing.T, root string) {
	t.Helper()
	locales := []string{"en", "en-US", "zh", "zh-Hans", "zh-Hans-CN"}
	for _, loc := range locales {
		dir := filepath.Join(root, "cldr-dates-full", "main", loc)
		if err := os.MkdirAll(dir, 0o777); err != nil {
			t.Fatalf("mkdir timezone names %s: %v", loc, err)
		}
		if err := os.WriteFile(filepath.Join(dir, "timeZoneNames.json"), []byte(timeZoneNamesJSON(loc)), 0o666); err != nil {
			t.Fatalf("write timeZoneNames %s: %v", loc, err)
		}
	}
	supp := filepath.Join(root, "cldr-core", "supplemental")
	metaZones := `{"supplemental":{"metaZones":{"metazoneInfo":{"timezone":{"America":{"Los_Angeles":[{"usesMetazone":{"_mzone":"America_Pacific"}}],"New_York":[{"usesMetazone":{"_mzone":"America_Eastern"}}]},"Europe":{"London":[{"usesMetazone":{"_mzone":"GMT"}}],"Moscow":[{"usesMetazone":{"_mzone":"Moscow","_to":"2011-01-01 00:00"}},{"usesMetazone":{"_mzone":"Europe_Further_Eastern","_from":"2011-01-01 00:00","_to":"2013-12-01 00:00"}},{"usesMetazone":{"_mzone":"Moscow","_from":"2013-12-01 00:00"}}]}}}}}}`
	if err := os.WriteFile(filepath.Join(supp, "metaZones.json"), []byte(metaZones), 0o666); err != nil {
		t.Fatalf("write metaZones: %v", err)
	}
	aliases := `{"keyword":{"u":{"tz":{"utc":{"_alias":"Etc/UTC UTC"},"inccu":{"_alias":"Asia/Calcutta Asia/Kolkata","_iana":"Asia/Kolkata"}}}}}`
	if err := os.WriteFile(filepath.Join(root, "cldr-bcp47", "bcp47", "timezone.json"), []byte(aliases), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestRunDisplayIdentityDoesNotDependOnRegistryOutput(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	root := filepath.Join(dir, "node_modules")
	writeRuntimeCLDRFixtures(t, root)
	versionPath := filepath.Join(dir, "VERSION")
	if err := os.WriteFile(versionPath, []byte("cldr=48.1.0\nicu=78\ntzdata=2025b\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg := Config{CLDRDir: root, OutDir: filepath.Join(dir, "with-registry"), VersionFile: versionPath, ProfileFile: writeLocaleProfileFixture(t, dir), TZDataLock: writeTZDataFixture(t), TimeZoneOut: filepath.Join(dir, "registry_data.go")}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	if err := Run(t.Context(), cfg, log); err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join(cfg.OutDir, "timezone", "data.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !containsAll(dewrapStringLiterals(string(want)), "Asia/Kolkata", "Asia/Calcutta") {
		t.Fatal("full generation omitted the pinned display relation")
	}
	cfg.OutDir, cfg.TimeZoneOut = filepath.Join(dir, "without-registry"), ""
	if err := Run(t.Context(), cfg, log); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(cfg.OutDir, "timezone", "data.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("omitting optional registry output changed the runtime display payload")
	}
	cfg.OutDir, cfg.TZDataLock = filepath.Join(dir, "missing-input"), ""
	if err := Run(t.Context(), cfg, log); err == nil || !strings.Contains(err.Error(), "-tzdata-lock") {
		t.Fatalf("missing identity source error = %v, want -tzdata-lock context", err)
	}
	if _, err := os.Stat(filepath.Join(cfg.OutDir, "timezone", "data.go")); !os.IsNotExist(err) {
		t.Fatalf("failed generation wrote a display payload: %v", err)
	}
}

func timeZoneNamesJSON(locale string) string {
	return `{"main":{"` + locale + `":{"dates":{"timeZoneNames":{"hourFormat":"+HH:mm;-HH:mm","gmtFormat":"GMT{0}","gmtZeroFormat":"GMT","regionFormat":"{0} Time","metazone":{"America_Pacific":{"long":{"generic":"Pacific Time","standard":"Pacific Standard Time","daylight":"Pacific Daylight Time"},"short":{"generic":"PT","standard":"PST","daylight":"PDT"}},"America_Eastern":{"long":{"generic":"Eastern Time","standard":"Eastern Standard Time","daylight":"Eastern Daylight Time"}}},"zone":{"America":{"Los_Angeles":{"exemplarCity":"Los Angeles"},"New_York":{"exemplarCity":"New York"}},"Europe":{"London":{"long":{"daylight":"British Summer Time"},"short":{"daylight":"BST"}}},"Etc":{"UTC":{"long":{"standard":"Coordinated Universal Time"},"short":{"standard":"UTC"}}}}}}}}}`
}
