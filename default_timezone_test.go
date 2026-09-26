package gointl

import (
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/agentable/go-intl/datetimeformat"
	"github.com/agentable/go-intl/internal/intltest"
)

func TestDefaultTimeZoneUsesExplicitValidation(t *testing.T) {
	if name := os.Getenv("GO_INTL_TEST_DEFAULT_ZONE"); name != "" {
		// The test owns this fresh process; no parallel test observes time.Local.
		time.Local = time.FixedZone(name, 1234)
		locales := intltest.LocaleList(t, "en")
		implicit, implicitErr := datetimeformat.New(locales, datetimeformat.Options{})
		explicit, explicitErr := datetimeformat.New(locales, datetimeformat.Options{TimeZone: String(name)})
		if (implicitErr == nil) != (explicitErr == nil) {
			t.Fatalf("default error = %v, explicit = %v", implicitErr, explicitErr)
		}
		if explicitErr != nil {
			if !errors.Is(implicitErr, ErrUnsupportedOption) || !strings.Contains(implicitErr.Error(), name) {
				t.Fatalf("default error lacks category/candidate: %v", implicitErr)
			}
			return
		}
		if implicit.ResolvedOptions().TimeZone != explicit.ResolvedOptions().TimeZone {
			t.Fatal("resolved time zones differ")
		}
		instant := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)
		a, err := implicit.Format(instant)
		if err != nil {
			t.Fatal(err)
		}
		b, err := explicit.Format(instant)
		if err != nil {
			t.Fatal(err)
		}
		if a != b {
			t.Fatal("default and explicit output differ")
		}
		return
	}
	t.Parallel()
	for _, name := range []string{"America/New_York", "US/Eastern", "+05:30", "+0530", "UTC", "Mars/Phobos", "+24:00"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestDefaultTimeZoneUsesExplicitValidation$")
			cmd.Env = append(os.Environ(), "GO_INTL_TEST_DEFAULT_ZONE="+name)
			if output, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("subprocess: %v\n%s", err, output)
			}
		})
	}
}
