package datetimeformat

import (
	"testing"
	"time"

	"github.com/agentable/go-intl/internal/intltest"
)

func TestPrimaryTimeZoneDisplayIdentity(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ locale, zone, primary, want string }{
		{"en-US", "Asia/Kolkata", "Asia/Kolkata", "India Standard Time"}, {"en-US", "Asia/Calcutta", "Asia/Kolkata", "India Standard Time"},
		{"en-US", "Europe/Kyiv", "Europe/Kyiv", "Eastern European Standard Time"}, {"en-US", "Europe/Kiev", "Europe/Kyiv", "Eastern European Standard Time"},
		{"en-US", "UTC", "UTC", "Coordinated Universal Time"}, {"en-US", "Etc/UTC", "UTC", "Coordinated Universal Time"},
		{"fr", "Asia/Kolkata", "Asia/Kolkata", "heure de l’Inde"},
	} {
		t.Run(tc.locale+tc.zone, func(t *testing.T) {
			t.Parallel()
			f, err := New(intltest.LocaleList(t, tc.locale), Options{TimeZone: new(tc.zone), TimeZoneName: new("long"), Hour: new("numeric"), Minute: new("2-digit")})
			if err != nil {
				t.Fatal(err)
			}
			if got := f.ResolvedOptions().TimeZone; got != tc.primary {
				t.Errorf("primary = %q, want %q", got, tc.primary)
			}
			assertZoneNameMethods(t, f, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), tc.want)
		})
	}
}
