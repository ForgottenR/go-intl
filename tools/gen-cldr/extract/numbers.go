package extract

import (
	"maps"

	"github.com/agentable/go-intl/tools/gen-cldr/cldr"
)

type Numbers = map[string]cldr.Numbers

type CurrencyData struct {
	Fractions  map[string]cldr.CurrencyFraction
	Currencies map[string]cldr.Currencies
}

func ExtractNumbers(raw map[string]cldr.Numbers, locales []string) Numbers {
	selected := localeSet(locales)
	out := make(Numbers, len(selected))
	for locale, numbers := range raw {
		if selected[locale] {
			out[locale] = numbers
		}
	}
	return out
}

func ExtractCurrencies(fractions map[string]cldr.CurrencyFraction, currencies map[string]cldr.Currencies, locales []string) CurrencyData {
	selected := localeSet(locales)
	filteredFractions := maps.Clone(fractions)
	filteredCurrencies := make(map[string]cldr.Currencies, len(selected))
	for locale, byCurrency := range currencies {
		if !selected[locale] {
			continue
		}
		filteredCurrencies[locale] = maps.Clone(byCurrency)
	}
	return CurrencyData{Fractions: filteredFractions, Currencies: filteredCurrencies}
}

func localeSet(locales []string) map[string]bool {
	out := make(map[string]bool, len(locales))
	for _, locale := range locales {
		out[locale] = true
	}
	return out
}
