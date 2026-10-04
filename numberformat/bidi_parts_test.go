package numberformat

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/agentable/go-intl/internal/intltest"
)

func TestBidiSignPartitions(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ tag, mark, minus, digits, nan string }{{"ar-EG", "\u061c", "-", "١٢٣", "ليس رقمًا"}, {"fa", "\u200e", "−", "۱۲۳", "ناعدد"}} {
		t.Run(tc.tag, func(t *testing.T) {
			t.Parallel()
			for _, v := range []struct {
				x        float64
				typ      PartType
				value    string
				negative bool
			}{{-123, PartInteger, tc.digits, true}, {123, PartInteger, tc.digits, false}, {math.Copysign(0, -1), PartInteger, map[string]string{"ar-EG": "٠", "fa": "۰"}[tc.tag], true}, {0, PartInteger, map[string]string{"ar-EG": "٠", "fa": "۰"}[tc.tag], false}, {math.NaN(), PartNaN, tc.nan, false}, {math.Inf(1), PartInfinity, "∞", false}, {math.Inf(-1), PartInfinity, "∞", true}} {
				f, err := New(intltest.LocaleList(t, tc.tag), Options{SignDisplay: new("always")})
				if err != nil {
					t.Fatal(err)
				}
				sign := Part{Type: PartPlusSign, Value: "+"}
				if v.negative {
					sign = Part{Type: PartMinusSign, Value: tc.minus}
				}
				want := []Part{{Type: PartLiteral, Value: tc.mark}, sign, {Type: v.typ, Value: v.value}}
				got := f.FormatToParts(Float(v.x))
				if !reflect.DeepEqual(got, want) {
					t.Errorf("parts(%v) = %#v, want %#v", v.x, got, want)
				}
				if f.Format(Float(v.x)) != partsText(want) {
					t.Errorf("text(%v) = %q, want %q", v.x, f.Format(Float(v.x)), partsText(want))
				}
			}
			for _, notation := range []string{"scientific", "engineering"} {
				f, err := New(intltest.LocaleList(t, tc.tag), Options{Notation: new(notation)})
				if err != nil {
					t.Fatal(err)
				}
				for _, x := range []float64{-0.0123, 1, 12345} {
					got := f.FormatToParts(Float(x))
					if partsText(got) != f.Format(Float(x)) {
						t.Fatal("text differs from parts")
					}
					found := false
					for i, p := range got {
						if p.Type == PartExponentMinusSign {
							found = true
							if p.Value != tc.minus || i == 0 || got[i-1] != (Part{Type: PartLiteral, Value: tc.mark}) {
								t.Errorf("exponent boundary = %#v", got)
							}
						}
					}
					if found != (x == -0.0123) {
						t.Errorf("exponent sign = %#v", got)
					}
				}
			}
			for _, opts := range []Options{{Style: new("unit"), Unit: new("meter")}, {Style: new("currency"), Currency: new("USD")}, {Style: new("currency"), Currency: new("USD"), CurrencyDisplay: new("code")}, {Style: new("currency"), Currency: new("USD"), CurrencyDisplay: new("name")}, {Style: new("percent")}} {
				f, err := New(intltest.LocaleList(t, tc.tag), opts)
				if err != nil {
					t.Fatal(err)
				}
				got := f.FormatToParts(Int(-123))
				if partsText(got) != f.Format(Int(-123)) {
					t.Fatal("style text differs from parts")
				}
				signs := 0
				for _, p := range got {
					if p.Type == PartMinusSign {
						signs++
						if p.Value != tc.minus {
							t.Errorf("style sign = %#v", got)
						}
					}
				}
				if signs != 1 || strings.Count(partsText(got), tc.mark) > 2 {
					t.Errorf("duplicate sign/literal = %#v", got)
				}
			}
			f, err := New(intltest.LocaleList(t, tc.tag), Options{})
			if err != nil {
				t.Fatal(err)
			}
			got, err := f.FormatRangeToParts(Int(-123), Int(-456))
			if err != nil {
				t.Fatal(err)
			}
			wantRange := []RangePart{
				{Type: PartLiteral, Value: tc.mark, Source: SourceShared},
				{Type: PartMinusSign, Value: tc.minus, Source: SourceShared},
				{Type: PartInteger, Value: tc.digits, Source: SourceStartRange},
				{Type: PartLiteral, Value: "–", Source: SourceShared},
				{Type: PartInteger, Value: map[string]string{"ar-EG": "٤٥٦", "fa": "۴۵۶"}[tc.tag], Source: SourceEndRange},
			}
			if !reflect.DeepEqual(got, wantRange) {
				t.Errorf("range = %#v, want %#v", got, wantRange)
			}
			text, err := f.FormatRange(Int(-123), Int(-456))
			if err != nil {
				t.Fatal(err)
			}
			joined := ""
			shared := false
			for _, p := range got {
				joined += p.Value
				shared = shared || p.Source == SourceShared
			}
			if joined != text || !shared {
				t.Errorf("range join/source = %#v", got)
			}
		})
	}
}

// Node 26.10.0 / ICU 78.3 keeps the interval literal shared, including the
// second endpoint's adjacent bidi mark.
func TestBidiMixedSignRanges(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ tag, mark, minus, first, second string }{
		{"ar-EG", "\u061c", "-", "١٢٣", "٤٥٦"},
		{"fa", "\u200e", "−", "۱۲۳", "۴۵۶"},
	} {
		for _, signDisplay := range []string{"auto", "always"} {
			for _, negativeStart := range []bool{true, false} {
				t.Run(tc.tag+"/"+signDisplay+"/"+map[bool]string{true: "negative start", false: "negative end"}[negativeStart], func(t *testing.T) {
					t.Parallel()
					f, err := New(intltest.LocaleList(t, tc.tag), Options{SignDisplay: new(signDisplay)})
					if err != nil {
						t.Fatal(err)
					}
					start, end := int64(123), int64(-456)
					startSign, endSign := "+", tc.minus
					if negativeStart {
						start, end = -123, 456
						startSign, endSign = tc.minus, "+"
					}
					want := []RangePart{}
					separator := "–"
					if negativeStart || signDisplay == "always" {
						separator = " – "
						typ := PartPlusSign
						if negativeStart {
							typ = PartMinusSign
						}
						want = append(want, RangePart{Type: PartLiteral, Value: tc.mark, Source: SourceStartRange}, RangePart{Type: typ, Value: startSign, Source: SourceStartRange})
					}
					want = append(want, RangePart{Type: PartInteger, Value: tc.first, Source: SourceStartRange})
					if !negativeStart || signDisplay == "always" {
						separator += tc.mark
					}
					want = append(want, RangePart{Type: PartLiteral, Value: separator, Source: SourceShared})
					if !negativeStart || signDisplay == "always" {
						typ := PartPlusSign
						if !negativeStart {
							typ = PartMinusSign
						}
						want = append(want, RangePart{Type: typ, Value: endSign, Source: SourceEndRange})
					}
					want = append(want, RangePart{Type: PartInteger, Value: tc.second, Source: SourceEndRange})
					parts, err := f.FormatRangeToParts(Int(start), Int(end))
					if err != nil || !reflect.DeepEqual(parts, want) {
						t.Errorf("range parts = %#v, %v; want %#v", parts, err, want)
					}
					text, err := f.FormatRange(Int(start), Int(end))
					if err != nil || text != rangePartsText(want) {
						t.Errorf("range text = %q, %v; want %q", text, err, rangePartsText(want))
					}
				})
			}
		}
	}
}
