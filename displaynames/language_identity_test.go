package displaynames

import (
	"testing"

	"github.com/agentable/go-intl/internal/ecma402"
	"github.com/agentable/go-intl/internal/intltest"
	"github.com/agentable/go-intl/locale"
)

func TestLanguageCanonicalIdentity(t *testing.T) {
	t.Parallel()

	f, err := New(intltest.LocaleList(t, "en"), Options{Type: new("language")})
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		input string
		want  string
	}{
		{"twi", "ak"},
		{"TWI", "ak"},
		{"twi-Latn-GH", "ak-Latn-GH"},
		{"TWI-LATN-GH-EMODENG", "ak-Latn-GH-emodeng"},
		{"ak-Latn-GH-emodeng", "ak-Latn-GH-emodeng"},
		{"en-Latn", "en-Latn"},
		{"no", "no"},
		{"nb", "nb"},
		{"no-Latn-NO", "no-Latn-NO"},
		{"nb-Latn-NO", "nb-Latn-NO"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()

			loc, err := locale.Parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if got := loc.String(); got != tc.want {
				t.Fatalf("Locale identity = %q, want %q", got, tc.want)
			}
			canonical, err := ecma402.CanonicalCodeForDisplayNames("language", tc.input)
			if err != nil || canonical != tc.want {
				t.Fatalf("DisplayNames identity = %q, %v, want %q", canonical, err, tc.want)
			}
			got, ok, err := f.Of(tc.input)
			want, wantOK, wantErr := f.Of(tc.want)
			if err != nil || wantErr != nil || got != want || ok != wantOK {
				t.Fatalf("Of(%q) = %q, %v, %v; Of(%q) = %q, %v, %v", tc.input, got, ok, err, tc.want, want, wantOK, wantErr)
			}
		})
	}
}

func TestNorwegianLanguageNames(t *testing.T) {
	t.Parallel()
	f, err := New(intltest.LocaleList(t, "en"), Options{Type: new("language")})
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ code, want string }{{"no", "Norwegian"}, {"nb", "Norwegian Bokmål"}} {
		got, ok, err := f.Of(tc.code)
		if err != nil || !ok || got != tc.want {
			t.Errorf("Of(%q) = %q, %v, %v; want %q", tc.code, got, ok, err, tc.want)
		}
	}
}
