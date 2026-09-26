package relativetimeformat

import (
	"strings"
	"testing"

	cldrrelative "github.com/agentable/go-intl/internal/cldr/relativetime"
)

func TestRelativeTimeConstructorRejectsIncompleteData(t *testing.T) {
	t.Parallel()
	valid := cldrrelative.RelativeTimeField{
		Future: map[string]string{"other": "in {0} units"},
		Past:   map[string]string{"other": "{0} units ago"},
	}
	for _, tc := range []struct {
		name   string
		mutate func(cldrrelative.RelativeTimeFields)
	}{
		{"missing unit", func(raw cldrrelative.RelativeTimeFields) { delete(raw, string(Year)) }},
		{"missing future other", func(raw cldrrelative.RelativeTimeFields) {
			raw[string(Day)][string(LongStyle)] = cldrrelative.RelativeTimeField{Future: map[string]string{}, Past: valid.Past}
		}},
		{"missing past other", func(raw cldrrelative.RelativeTimeFields) {
			raw[string(Day)][string(LongStyle)] = cldrrelative.RelativeTimeField{Future: valid.Future, Past: map[string]string{}}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			raw := cldrrelative.RelativeTimeFields{}
			for _, unit := range relativeTimeUnits {
				raw[string(unit)] = map[string]cldrrelative.RelativeTimeField{string(LongStyle): valid}
			}
			tc.mutate(raw)
			_, err := compileRelativeTimeFields(raw, LongStyle)
			if err == nil || !strings.Contains(err.Error(), "day") && !strings.Contains(err.Error(), "year") {
				t.Fatalf("compile error = %v", err)
			}
		})
	}
	for _, tag := range cldrrelative.SupportedLocales() {
		loc, ok := cldrrelative.ResolveLocale(tag)
		if !ok {
			t.Fatalf("%s cannot resolve", tag)
		}
		for _, style := range []Style{LongStyle, ShortStyle, NarrowStyle} {
			if _, err := compileRelativeTimeFields(cldrrelative.FieldsFor(loc), style); err != nil {
				t.Fatalf("%s %s: %v", tag, style, err)
			}
		}
	}
}
