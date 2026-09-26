package tz

import (
	"errors"
	"testing"
	"time"
)

func TestDefaultLocationUsesIANAName(t *testing.T) {
	t.Parallel()

	local, err := time.LoadLocation("US/Eastern")
	if err != nil {
		t.Fatal(err)
	}
	name, loc, err := defaultLocation(local)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := name, "America/New_York"; got != want {
		t.Fatalf("defaultLocation(US/Eastern) name = %q, want %q", got, want)
	}
	if got := LookupAt(loc, time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)).OffsetMs; got != -5*3600*1000 {
		t.Fatalf("defaultLocation(US/Eastern) offset = %d, want EST", got)
	}
}

func TestDefaultLocationFallsBackToUTCForNilLocal(t *testing.T) {
	t.Parallel()

	name, loc, err := defaultLocation(nil)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := name, "UTC"; got != want {
		t.Fatalf("defaultLocation(nil) name = %q, want %q", got, want)
	}
	if loc != time.UTC {
		t.Fatalf("defaultLocation(nil) location = %v, want UTC", loc)
	}
}

func TestDefaultLocationReturnsUsableLocationForLocal(t *testing.T) {
	t.Parallel()

	name, loc, err := defaultLocation(time.Local)
	if err != nil {
		t.Fatal(err)
	}
	if name == "" {
		t.Fatal("defaultLocation(Local) name is empty")
	}
	if loc == nil {
		t.Fatal("defaultLocation(Local) location is nil")
	}
	if info := LookupAt(loc, time.Now().UTC()); info.Name == "" {
		t.Fatal("LookupAt(default local) name is empty")
	}
}

func TestDefaultLocationUsesFixedOffsetName(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		seconds  int
		wantName string
		wantMs   int64
	}{
		{name: "+05:30", seconds: 5*3600 + 30*60, wantName: "+05:30", wantMs: (5*3600 + 30*60) * 1000},
		{name: "+0530", seconds: 5*3600 + 30*60, wantName: "+05:30", wantMs: (5*3600 + 30*60) * 1000},
		{name: "-00", seconds: 0, wantName: "+00:00", wantMs: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			name, loc, err := defaultLocation(time.FixedZone(tc.name, tc.seconds))
			if err != nil {
				t.Fatal(err)
			}
			if name != tc.wantName {
				t.Fatalf("defaultLocation(%s) name = %q, want %q", tc.name, name, tc.wantName)
			}
			if got := LookupAt(loc, time.Unix(0, 0)).OffsetMs; got != tc.wantMs {
				t.Fatalf("defaultLocation(%s) offset = %d, want %d", tc.name, got, tc.wantMs)
			}
		})
	}
}

func TestDefaultLocationRejectsInvalidNamedLocation(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"Mars/Phobos", "+24:00"} {
		got, loc, err := defaultLocation(time.FixedZone(name, 1234))
		if got != name || loc != nil || !errors.Is(err, ErrUnsupportedTimeZone) {
			t.Errorf("defaultLocation(%q) = %q, %v, %v", name, got, loc, err)
		}
	}
}

func TestDefaultOverrideForTest(t *testing.T) {
	restore, err := OverrideDefaultForTest("Asia/Shanghai")
	if err != nil {
		t.Fatalf("OverrideDefaultForTest() error = %v", err)
	}
	t.Cleanup(restore)

	name, loc, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := name, "Asia/Shanghai"; got != want {
		t.Fatalf("Default() name = %q, want %q", got, want)
	}
	if got := LookupAt(loc, time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)).OffsetMs; got != 8*3600*1000 {
		t.Fatalf("Default() offset = %d, want CST", got)
	}
}

func TestDefaultOverrideCanonicalizesLink(t *testing.T) {
	restore, err := OverrideDefaultForTest("US/Eastern")
	if err != nil {
		t.Fatalf("OverrideDefaultForTest() error = %v", err)
	}
	t.Cleanup(restore)

	name, _, err := Default()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := name, "America/New_York"; got != want {
		t.Fatalf("Default() name = %q, want %q", got, want)
	}
}

func TestDefaultOverrideRejectsUnsupportedTimeZone(t *testing.T) {
	t.Parallel()

	restore, err := OverrideDefaultForTest("Mars/Olympus")
	if err == nil {
		t.Fatal("OverrideDefaultForTest() error = nil, want unsupported time zone")
	}
	if restore != nil {
		t.Fatal("OverrideDefaultForTest() restore is non-nil, want nil")
	}
	if !errors.Is(err, ErrUnsupportedTimeZone) {
		t.Fatalf("OverrideDefaultForTest() error = %v, want ErrUnsupportedTimeZone", err)
	}
}
