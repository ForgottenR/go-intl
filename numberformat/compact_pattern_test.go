package numberformat

import (
	"reflect"
	"testing"

	cldrpattern "github.com/agentable/go-intl/internal/pattern"
)

func TestCompactPatternQuotedLiterals(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		pattern string
		parts   []Part
	}{
		{"00\u00a0тис'.'", []Part{{Type: PartInteger, Value: "10"}, {Type: PartLiteral, Value: "\u00a0"}, {Type: PartCompact, Value: "тис."}}},
		{"00 o''clock", []Part{{Type: PartInteger, Value: "10"}, {Type: PartLiteral, Value: " "}, {Type: PartCompact, Value: "o'clock"}}},
		{"'0;'00K", []Part{{Type: PartCompact, Value: "0;"}, {Type: PartInteger, Value: "10"}, {Type: PartCompact, Value: "K"}}},
		{"00'a;0'", []Part{{Type: PartInteger, Value: "10"}, {Type: PartCompact, Value: "a;0"}}},
		{"00K;'-'00K", []Part{{Type: PartInteger, Value: "10"}, {Type: PartCompact, Value: "K"}}},
	} {
		t.Run(tc.pattern, func(t *testing.T) {
			t.Parallel()
			parsed, err := cldrpattern.ParseCompact(4, tc.pattern)
			if err != nil {
				t.Fatal(err)
			}
			got := compileCompactAffixPattern(parsed, Part{}).append([]Part{{Type: PartInteger, Value: "10"}})
			if !reflect.DeepEqual(got, tc.parts) {
				t.Errorf("parts = %#v, want %#v", got, tc.parts)
			}
			if got := parsed.Exponent; got != 3 {
				t.Errorf("exponent = %d, want 3", got)
			}
		})
	}
}
