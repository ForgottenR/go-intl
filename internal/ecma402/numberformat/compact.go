package ecma402nf

import (
	"github.com/agentable/go-intl/internal/decimal"
)

const (
	MinCompactMagnitude = 3
	MaxCompactMagnitude = 14
)

type CompactExponentLookup func(magnitude int) (exponent int, ok bool)

// ResolveCompactMagnitude scales d into its compact mantissa via the shared
// computeExponent engine (which owns the post-rounding carry recheck) and returns
// the scaled value, its magnitude, the compact exponent, and whether compact
// notation applies. Rounding can promote a value below MinCompactMagnitude
// into the first available compact magnitude.
func ResolveCompactMagnitude(d decimal.Decimal, digitOptions ResolvedDigitOptions, lookup CompactExponentLookup) (decimal.Decimal, int, int, bool) {
	exponent, magnitude, ok := computeExponent(d, digitOptions, compactExponentForMagnitude(lookup))
	if !ok {
		return d, 0, 0, false
	}
	return decimal.Scale10(d, -int32(exponent)), magnitude, exponent, true // #nosec G115 -- compact exponents are small generated data keys.
}

func compactExponentForMagnitude(lookup CompactExponentLookup) exponentForMagnitude {
	return func(magnitude int) (int, bool) {
		if magnitude < MinCompactMagnitude || lookup == nil {
			return 0, false
		}
		return lookup(magnitude)
	}
}
