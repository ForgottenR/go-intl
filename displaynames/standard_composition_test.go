package displaynames_test

import (
	"testing"

	"github.com/agentable/go-intl/displaynames"
	"github.com/agentable/go-intl/internal/intltest"
)

func TestStandardLanguageComposition(t *testing.T) {
	t.Parallel()
	for _, style := range []string{"long", "short", "narrow"} {
		for _, dataLocale := range []string{"en", "en-US"} {
			t.Run(dataLocale+"/"+style, func(t *testing.T) {
				t.Parallel()
				d, err := displaynames.New(intltest.LocaleList(t, dataLocale), displaynames.Options{Type: new("language"), LanguageDisplay: new("standard"), Style: new(style)})
				if err != nil {
					t.Fatal(err)
				}
				region := "United States"
				if style != "long" {
					region = "US"
				}
				for _, tc := range []struct{ code, want string }{
					{"en", "English"}, {"zh-Hans", "Chinese (Simplified)"}, {"en-Cyrl", "English (Cyrillic)"},
					{"en-US", "English (" + region + ")"}, {"en-Cyrl-US", "English (Cyrillic, " + region + ")"},
					{"sl-Latn-SI-rozaj", "Slovenian (Latin, Slovenia, Resian)"},
				} {
					if got, ok, err := d.Of(tc.code); err != nil || !ok || got != tc.want {
						t.Errorf("Of(%s) = %q,%t,%v, want %q", tc.code, got, ok, err, tc.want)
					}
				}
			})
		}
	}
	for _, tc := range []struct {
		loc, display, code, fallback, want string
		ok                                 bool
	}{
		{"zh", "standard", "en-Cyrl-US", "code", "英语（西里尔文，美国）", true},
		{"en", "standard", "en-Qaaa-US", "code", "English (Qaaa, United States)", true},
		{"en", "standard", "en-Qaaa-US", "none", "", false},
		{"en", "dialect", "zh-Hans-CN", "code", "Simplified Chinese (China)", true},
		{"en", "dialect", "en-US", "code", "American English", true},
	} {
		d, err := displaynames.New(intltest.LocaleList(t, tc.loc), displaynames.Options{Type: new("language"), LanguageDisplay: new(tc.display), Fallback: new(tc.fallback)})
		if err != nil {
			t.Fatal(err)
		}
		if got, ok, err := d.Of(tc.code); err != nil || ok != tc.ok || got != tc.want {
			t.Errorf("%s/%s Of(%s) = %q,%t,%v, want %q,%t", tc.loc, tc.display, tc.code, got, ok, err, tc.want, tc.ok)
		}
	}
}
