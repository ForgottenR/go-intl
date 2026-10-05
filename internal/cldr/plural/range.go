package plural

import pluralop "github.com/agentable/go-intl/internal/plural"

// ResolveRange returns the explicit CLDR range category, or other when
// the sparse range data has no row for the category pair.
func ResolveRange(loc string, start, end pluralop.Category) pluralop.Category {
	if category, ok := Range(loc, start, end); ok {
		return category
	}
	return pluralop.Other
}
