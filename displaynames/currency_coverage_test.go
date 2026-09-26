package displaynames_test

import (
	"slices"
	"testing"

	gointl "github.com/agentable/go-intl"
	"github.com/agentable/go-intl/displaynames"
	"github.com/agentable/go-intl/internal/intltest"
	"github.com/agentable/go-intl/locale"
	"github.com/agentable/go-intl/numberformat"
	"github.com/agentable/go-intl/option"
)

func TestPinnedCurrencyNamesAndMembership(t *testing.T) {
	t.Parallel()
	codes := gointl.SupportedCurrencies()
	if !slices.IsSorted(codes) {
		t.Fatal("currencies not sorted")
	}
	for _, code := range []string{"USD", "JPY", "AED", "SAR", "BDT", "AFN", "SLL", "DEM", "XXX"} {
		if !slices.Contains(codes, code) {
			t.Errorf("missing %s", code)
		}
	}
	if slices.Contains(codes, "ZZZ") {
		t.Error("unknown ZZZ listed")
	}
	for _, tag := range []string{"en", "fr"} {
		locales := locale.List{intltest.Locale(t, tag)}
		dn, err := displaynames.New(locales, displaynames.Options{Type: option.String("currency")})
		if err != nil {
			t.Fatal(err)
		}
		for _, code := range []string{"AED", "SAR", "BDT"} {
			name, ok, err := dn.Of(code)
			if err != nil || !ok || name == code {
				t.Errorf("%s %s name=%q ok=%t err=%v", tag, code, name, ok, err)
			}
			nf, err := numberformat.New(locales, numberformat.Options{Style: option.String("currency"), Currency: option.String(code), CurrencyDisplay: option.String("name")})
			if err != nil {
				t.Fatal(err)
			}
			text := nf.Format(numberformat.Int(2))
			if text == "" || text == "2 "+code {
				t.Errorf("%s %s format=%q", tag, code, text)
			}
		}
	}
}
