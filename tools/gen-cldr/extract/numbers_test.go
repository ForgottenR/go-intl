package extract

import (
	"testing"

	"github.com/agentable/go-intl/tools/gen-cldr/cldr"
)

func TestExtractCurrenciesProfilesLocaleNamesAndKeepsFractions(t *testing.T) {
	t.Parallel()

	fractions := map[string]cldr.CurrencyFraction{
		"DEFAULT": {Digits: 2},
		"USD":     {Digits: 2},
		"XXX":     {Digits: 0},
	}
	currencies := map[string]cldr.Currencies{
		"en": {
			"USD": {Canonical: "US Dollar", Symbol: "$"},
			"XXX": {Canonical: "Test Currency", Symbol: "X"},
		},
		"fr": {
			"EUR": {Canonical: "euro", Symbol: "EUR"},
		},
	}

	got := ExtractCurrencies(fractions, currencies, []string{"en"})

	if len(got.Fractions) != len(fractions) {
		t.Fatalf("ExtractCurrencies fractions length = %d, want %d", len(got.Fractions), len(fractions))
	}
	if _, ok := got.Fractions["XXX"]; !ok {
		t.Fatalf("ExtractCurrencies fractions dropped non-display currency XXX")
	}
	got.Fractions["USD"] = cldr.CurrencyFraction{Digits: 9}
	if fractions["USD"].Digits != 2 {
		t.Fatalf("ExtractCurrencies fractions share input map; input USD digits = %d, want 2", fractions["USD"].Digits)
	}

	filtered, ok := got.Currencies["en"]
	if !ok {
		t.Fatalf(`ExtractCurrencies currencies missing selected locale "en"`)
	}
	if _, ok := got.Currencies["fr"]; ok {
		t.Fatalf(`ExtractCurrencies currencies kept unselected locale "fr"`)
	}
	if _, ok := filtered["USD"]; !ok {
		t.Fatalf("ExtractCurrencies currencies dropped allowlisted USD")
	}
	if _, ok := filtered["XXX"]; !ok {
		t.Fatal("ExtractCurrencies dropped pinned XXX name")
	}
}
