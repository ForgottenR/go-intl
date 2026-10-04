package testcontract

import (
	"encoding/json/v2"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/agentable/go-intl/tools/conformance"
)

const fixtureAdapterRoot = "GO_INTL_FIXTURE_ADAPTER_ROOT"

// FixtureRunnerRoot selects the suite supplied by a child-test invocation.
func FixtureRunnerRoot(defaultRoot string) string {
	if root := os.Getenv(fixtureAdapterRoot); root != "" {
		return root
	}
	return defaultRoot
}

// FixtureRunnerChild runs the package's actual fixture callback in a child test.
func FixtureRunnerChild(t *testing.T, run func(*testing.T, conformance.Fixture)) bool {
	t.Helper()

	root := FixtureRunnerRoot("")
	if root == "" {
		return false
	}
	conformance.RunFixtures(t, root, run)
	return true
}

// AssertFixtureRunner observes success or an assertion failure through RunFixtures.
func AssertFixtureRunner(t *testing.T, packageName, testName string, fixture conformance.Fixture, wantFailure string) {
	t.Helper()

	AssertFixtureSuite(t, packageName, testName, fixture, nil, wantFailure)
}

// AssertFixtureSuite checks the actual test entry with optional suite ledgers.
func AssertFixtureSuite(t *testing.T, packageName, testName string, fixture conformance.Fixture, files map[string]string, wantFailure string) string {
	t.Helper()

	root := filepath.Join(t.TempDir(), packageName)
	dir := filepath.Join(root, "testdata", "conformance", "manual")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal([]conformance.Fixture{fixture})
	if err != nil {
		t.Fatal(err)
	}
	name := "observation.json"
	if fixture.ErrorCode != "" {
		name = "errors.json"
	}
	if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(root, "testdata", name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	//nolint:gosec // Test helper re-execs the current test binary with a fixed flag and an escaped test name.
	cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^"+regexp.QuoteMeta(testName)+"$")
	cmd.Env = append(os.Environ(), fixtureAdapterRoot+"="+root)
	output, err := cmd.CombinedOutput()
	if wantFailure == "" {
		if err != nil {
			t.Fatalf("fixture runner rejected declared observations: %v\n%s", err, output)
		}
		return string(output)
	}
	if err == nil {
		t.Fatalf("fixture runner accepted incorrect observation, want %q\n%s", wantFailure, output)
	}
	if !strings.Contains(string(output), wantFailure) {
		t.Fatalf("fixture runner failure = %s, want %q", output, wantFailure)
	}
	return string(output)
}
