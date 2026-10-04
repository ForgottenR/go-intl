package timezone

import (
	"reflect"
	"testing"
	"time"

	cldrlocale "github.com/agentable/go-intl/internal/cldr/locale"
)

func TestDisplayIdentityFacts(t *testing.T) {
	t.Parallel()
	loc, ok := cldrlocale.ResolveLocale("en")
	if !ok {
		t.Fatal("missing en")
	}
	for _, pair := range [][2]string{{"Asia/Kolkata", "Asia/Calcutta"}, {"Europe/Kyiv", "Europe/Kiev"}, {"UTC", "Etc/UTC"}} {
		for _, instant := range []int64{time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli(), time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC).UnixMilli()} {
			if a, b := TimeZoneMetazone(pair[0], instant), TimeZoneMetazone(pair[1], instant); a != b {
				t.Errorf("metazone %v = %q/%q", pair, a, b)
			}
		}
		namesOnce.Do(loadNames)
		if a, b := timeZoneNamesByLocale[loc][displayKey(pair[0])], timeZoneNamesByLocale[loc][displayKey(pair[1])]; a != b {
			t.Errorf("zone names %v = %v/%v", pair, a, b)
		}
		if a, b := exemplarCity(loc, displayKey(pair[0])), exemplarCity(loc, displayKey(pair[1])); a != b {
			t.Errorf("cities %v = %q/%q", pair, a, b)
		}
		metazonePeriodOnce.Do(loadMetazonePeriods)
		if !reflect.DeepEqual(metazonePeriodsForZone(displayKey(pair[0])), metazonePeriodsForZone(displayKey(pair[1]))) {
			t.Errorf("historical periods differ for %v", pair)
		}
	}
}
