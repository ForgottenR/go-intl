package tz

import (
	"os"
	"strings"
	"sync"
	"time"
)

var defaultProvider = struct {
	sync.RWMutex
	location func() *time.Location
}{location: func() *time.Location { return time.Local }}

// Default returns the host default time zone snapshot used by DateTimeFormat
// construction when no explicit timeZone option is provided.
func Default() (string, *time.Location, error) {
	defaultProvider.RLock()
	location := defaultProvider.location
	defaultProvider.RUnlock()
	return defaultLocation(location())
}

// OverrideDefaultForTest replaces the default time-zone provider for tests and
// returns a restore function.
func OverrideDefaultForTest(name string) (func(), error) {
	loc, err := Resolve(name)
	if err != nil {
		return nil, err
	}
	defaultProvider.Lock()
	previous := defaultProvider.location
	defaultProvider.location = func() *time.Location { return loc }
	defaultProvider.Unlock()
	return func() {
		defaultProvider.Lock()
		defaultProvider.location = previous
		defaultProvider.Unlock()
	}, nil
}

func defaultLocation(local *time.Location) (string, *time.Location, error) {
	if local == nil {
		return "UTC", time.UTC, nil
	}
	if name := local.String(); name != "" && name != "Local" {
		return canonicalLocation(name)
	}
	if name := localtimeLinkName(); name != "" {
		return canonicalLocation(name)
	}
	return "UTC", time.UTC, nil
}

func canonicalLocation(name string) (string, *time.Location, error) {
	loc, err := Resolve(name)
	if err != nil {
		return name, nil, err
	}
	canonical := loc.String()
	return canonical, loc, nil
}

func localtimeLinkName() string {
	link, err := os.Readlink("/etc/localtime")
	if err != nil {
		return ""
	}
	const marker = "zoneinfo/"
	idx := strings.LastIndex(link, marker)
	if idx < 0 {
		return ""
	}
	return strings.TrimPrefix(link[idx+len(marker):], "/")
}
