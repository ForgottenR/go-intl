package main

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/agentable/go-intl/tools/gen-cldr/cldr"
	"github.com/agentable/go-intl/tools/gen-cldr/tzdb"
)

func TestTimeZoneDisplayKeysPreserveRows(t *testing.T) {
	t.Parallel()
	data := cldr.Metazones{
		ZoneToMetazones: map[string][]cldr.MetazonePeriod{"Old/City": {{Metazone: "Before", Start: -100, End: 0}, {Metazone: "After", Start: 0, End: 100}}},
		ZoneNames:       map[string]map[string]cldr.MetazoneNames{"en": {"Old/City": {LongStandard: "City Time"}}},
		ExemplarCities:  map[string]map[string]string{"en": {"Old/City": "City"}},
	}
	history := slices.Clone(data.ZoneToMetazones["Old/City"])
	registry := tzdb.Registry{Records: []tzdb.Record{{Identifier: "Old/City", Primary: "New/City"}, {Identifier: "New/City", Primary: "New/City"}}}
	got, err := timeZoneDisplayKeys(data, registry, map[string]string{"New/City": "Old/City"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, map[string]string{"New/City": "Old/City"}) || !reflect.DeepEqual(data.ZoneToMetazones["Old/City"], history) || data.ExemplarCities["en"]["Old/City"] != "City" {
		t.Fatalf("relation/source = %v / %v", got, data)
	}
	for _, field := range []string{"periods", "names", "cities"} {
		t.Run(field, func(t *testing.T) {
			// This subtest mutates its parent-owned fixture; do not run in parallel.
			switch field {
			case "periods":
				data.ZoneToMetazones["New/City"] = []cldr.MetazonePeriod{{Metazone: "Conflict"}}
			case "names":
				data.ZoneNames["en"]["New/City"] = cldr.MetazoneNames{LongStandard: "Conflict"}
			case "cities":
				data.ExemplarCities["en"]["New/City"] = "Conflict"
			}
			_, err := timeZoneDisplayKeys(data, registry, map[string]string{"New/City": "Old/City"})
			if err == nil || !strings.Contains(err.Error(), "Old/City") || !strings.Contains(err.Error(), "New/City") {
				t.Fatalf("conflict error = %v", err)
			}
			delete(data.ZoneToMetazones, "New/City")
			delete(data.ZoneNames["en"], "New/City")
			delete(data.ExemplarCities["en"], "New/City")
		})
	}
}
