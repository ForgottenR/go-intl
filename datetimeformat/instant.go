package datetimeformat

import (
	"time"

	"github.com/agentable/go-intl/internal/ecma402"
)

const maxEpochSeconds int64 = 8_640_000_000_000

// normalizeInstant applies the ECMA-262 TimeClip boundary before truncating
// toward whole milliseconds. UnixNano cannot represent the full Date domain.
func normalizeInstant(t time.Time, name string) (time.Time, error) {
	t = t.Round(0)
	seconds := t.Unix()
	nanos := t.Nanosecond()
	if seconds > maxEpochSeconds || seconds < -maxEpochSeconds ||
		(seconds == maxEpochSeconds && nanos != 0) {
		return time.Time{}, ecma402.InvalidValueErrorExpected("DateTimeFormat", name, t.String(), "", "an instant within ±8,640,000,000,000,000 epoch milliseconds", nil)
	}
	// Unix uses floor division for negative instants; reconstruct a truncation
	// toward zero without losing the nanosecond remainder at the negative bound.
	millis := seconds*1000 + int64(nanos)/1_000_000
	if seconds < 0 && nanos%1_000_000 != 0 {
		millis++
	}
	return time.UnixMilli(millis).UTC(), nil
}
