package timezone

import "testing"

func TestOffsetPatternSeconds(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		pattern  string
		offset   int64
		long     bool
		nu, want string
	}{
		{rootHourFormat, 1000, false, "latn", "+0:00:01"},
		{rootHourFormat, -7000, true, "arab", "-٠٠:٠٠:٠٧"},
		{rootHourFormat, 561000, true, "latn", "+00:09:21"},
		{"+H.mm;-H.mm", -561000, false, "latn", "-0.09.21"},
		{rootHourFormat, 540000, true, "latn", "+00:09"},
	} {
		if got := offsetPattern(tc.pattern, tc.offset, tc.long, tc.nu); got != tc.want {
			t.Errorf("offsetPattern(%d) = %q, want %q", tc.offset, got, tc.want)
		}
	}
}
