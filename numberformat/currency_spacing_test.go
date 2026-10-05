package numberformat_test

import (
	"math"
	"strings"
	"testing"

	gointl "github.com/agentable/go-intl"
	"github.com/agentable/go-intl/internal/intltest"
	"github.com/agentable/go-intl/numberformat"
)

func TestCurrencySpacingReusesConstructor(t *testing.T) {
	t.Parallel()

	f, err := numberformat.New(intltest.LocaleList(t, "en-US"), numberformat.Options{
		Style:           gointl.String(numberformat.CurrencyStyle),
		Currency:        gointl.String("ZWD"),
		CurrencyDisplay: gointl.String(numberformat.CurrencyDisplayCode),
	})
	if err != nil {
		t.Fatal(err)
	}
	// Keep this sequence on one instance to expose retained per-call state.
	for _, tc := range []struct {
		value numberformat.Value
		want  string
	}{
		{numberformat.Int(10000), "ZWD\u00a010,000"},
		{numberformat.Float(math.Inf(1)), "ZWD∞"},
		{numberformat.Float(math.NaN()), "ZWDNaN"},
		{numberformat.Int(10000), "ZWD\u00a010,000"},
	} {
		if got := f.Format(tc.value); got != tc.want {
			t.Errorf("Format() = %q, want %q", got, tc.want)
		}
		var parts strings.Builder
		for _, part := range f.FormatToParts(tc.value) {
			parts.WriteString(part.Value)
		}
		if got := parts.String(); got != tc.want {
			t.Errorf("joined FormatToParts() = %q, want %q", got, tc.want)
		}
	}
}
