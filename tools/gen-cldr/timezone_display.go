package main

import (
	"fmt"
	"maps"
	"reflect"
	"slices"

	"github.com/agentable/go-intl/tools/gen-cldr/cldr"
	"github.com/agentable/go-intl/tools/gen-cldr/tzdb"
)

// The relation joins the verified primary namespace to existing CLDR display
// keys. It leaves all source rows and historical periods intact.
func timeZoneDisplayKeys(data cldr.Metazones, registry tzdb.Registry, aliases map[string]string) (map[string]string, error) {
	out := map[string]string{}
	for _, rec := range registry.Records {
		key := aliases[rec.Primary]
		if rec.Identifier != rec.Primary || key == "" || key == rec.Primary {
			continue
		}
		if err := validateDisplayPair(data.ZoneToMetazones, rec.Primary, key, "metaZones.json"); err != nil {
			return nil, err
		}
		for _, tag := range slices.Sorted(maps.Keys(data.ZoneNames)) {
			if err := validateDisplayPair(data.ZoneNames[tag], rec.Primary, key, tag+"/timeZoneNames.json zone"); err != nil {
				return nil, err
			}
		}
		for _, tag := range slices.Sorted(maps.Keys(data.ExemplarCities)) {
			if err := validateDisplayPair(data.ExemplarCities[tag], rec.Primary, key, tag+"/timeZoneNames.json exemplarCity"); err != nil {
				return nil, err
			}
		}
		out[rec.Primary] = key
	}
	return out, nil
}

func validateDisplayPair[T any](rows map[string]T, primary, key, source string) error {
	a, hasA := rows[primary]
	b, hasB := rows[key]
	if hasA && hasB && !reflect.DeepEqual(a, b) {
		return fmt.Errorf("timezone display %s: conflicting keys %q and %q", source, primary, key)
	}
	return nil
}
