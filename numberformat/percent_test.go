package numberformat

import (
	"math"
	"reflect"
	"testing"

	"github.com/agentable/go-intl/internal/intltest"
	"github.com/agentable/go-intl/option"
)

func TestPercentPatterns(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		locale string
		value  float64
		text   string
		parts  []Part
	}{
		{"fr", .12, "12\u00a0%", []Part{{PartInteger, "12"}, {PartLiteral, "\u00a0"}, {PartPercentSign, "%"}}},
		{"fr", -.12, "-12\u00a0%", []Part{{PartMinusSign, "-"}, {PartInteger, "12"}, {PartLiteral, "\u00a0"}, {PartPercentSign, "%"}}},
		{"tr", .12, "%12", []Part{{PartPercentSign, "%"}, {PartInteger, "12"}}},
		{"tr", -.12, "-%12", []Part{{PartMinusSign, "-"}, {PartPercentSign, "%"}, {PartInteger, "12"}}},
		{"en", .12, "12%", []Part{{PartInteger, "12"}, {PartPercentSign, "%"}}},
		{"ar-EG-u-nu-arab", -.12, "\u061c-١٢٪\u061c", []Part{{PartLiteral, "\u061c"}, {PartMinusSign, "-"}, {PartInteger, "١٢"}, {PartPercentSign, "٪"}, {PartLiteral, "\u061c"}}},
		{"ar-u-nu-latn", -.12, "\u200e-12\u200e%\u200e", []Part{{PartLiteral, "\u200e"}, {PartMinusSign, "-"}, {PartInteger, "12"}, {PartLiteral, "\u200e"}, {PartPercentSign, "%"}, {PartLiteral, "\u200e"}}},
	} {
		t.Run(tc.locale+tc.text, func(t *testing.T) {
			t.Parallel()
			f, err := New(intltest.LocaleList(t, tc.locale), Options{Style: option.String("percent")})
			if err != nil {
				t.Fatal(err)
			}
			got := f.FormatToParts(Float(tc.value))
			if !reflect.DeepEqual(got, tc.parts) {
				t.Errorf("parts = %#v, want %#v", got, tc.parts)
			}
			if text := f.Format(Float(tc.value)); text != tc.text || text != partsText(got) {
				t.Errorf("Format = %q, want %q", text, tc.text)
			}
		})
	}
}

func TestPercentPatternNotationAndSpecials(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		locale, notation string
		value            float64
		text             string
	}{
		{"fr", "standard", math.NaN(), "NaN\u00a0%"}, {"tr", "standard", math.Inf(-1), "-%∞"},
		{"fr", "scientific", 12, "1E3\u00a0%"}, {"tr", "scientific", 12, "%1E3"},
		{"fr", "compact", 120, "12\u00a0k\u00a0%"}, {"tr", "compact", 120, "%12\u00a0B"},
	} {
		f, err := New(intltest.LocaleList(t, tc.locale), Options{Style: option.String("percent"), Notation: option.String(tc.notation)})
		if err != nil {
			t.Fatal(err)
		}
		if got := f.Format(Float(tc.value)); got != tc.text {
			t.Errorf("%s %s = %q, want %q", tc.locale, tc.notation, got, tc.text)
		}
	}
}

func TestPercentRangePatternSources(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		locale, text string
		parts        []RangePart
	}{
		{"fr", "12–34\u00a0%", []RangePart{{PartInteger, "12", SourceStartRange}, {PartLiteral, "–", SourceShared}, {PartInteger, "34", SourceEndRange}, {PartLiteral, "\u00a0", SourceShared}, {PartPercentSign, "%", SourceShared}}},
		{"tr", "%12 – %34", []RangePart{{PartPercentSign, "%", SourceStartRange}, {PartInteger, "12", SourceStartRange}, {PartLiteral, " – ", SourceShared}, {PartPercentSign, "%", SourceEndRange}, {PartInteger, "34", SourceEndRange}}},
	} {
		f, err := New(intltest.LocaleList(t, tc.locale), Options{Style: option.String("percent")})
		if err != nil {
			t.Fatal(err)
		}
		parts, err := f.FormatRangeToParts(Float(.12), Float(.34))
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(parts, tc.parts) {
			t.Errorf("%s range parts = %#v, want %#v", tc.locale, parts, tc.parts)
		}
		text, err := f.FormatRange(Float(.12), Float(.34))
		if err != nil || text != tc.text || rangePartsText(parts) != text {
			t.Errorf("%s range = %q, %v; want %q", tc.locale, text, err, tc.text)
		}
	}
}
