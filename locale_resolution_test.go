package gointl

import (
	"testing"

	"github.com/agentable/go-intl/datetimeformat"
	"github.com/agentable/go-intl/internal/intltest"
	"github.com/agentable/go-intl/numberformat"
)

func TestNumberFormatLocaleExtensionOption(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name   string
		option *string
		want   string
	}{
		{"absent", nil, "en-u-nu-arab"},
		{"equal", String("arab"), "en-u-nu-arab"},
		{"different", String("latn"), "en"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f, err := numberformat.New(intltest.LocaleList(t, "en-u-nu-arab"), numberformat.Options{NumberingSystem: tt.option})
			if err != nil {
				t.Fatal(err)
			}
			if got := f.ResolvedOptions().Locale.String(); got != tt.want {
				t.Fatalf("locale = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDateTimeFormatLocaleExtensionOptions(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name    string
		options datetimeformat.Options
		want    string
	}{
		{"absent", datetimeformat.Options{}, "en-u-ca-gregory-hc-h23-nu-arab"},
		{"equal", datetimeformat.Options{Calendar: String("gregory"), HourCycle: String("h23"), NumberingSystem: String("arab")}, "en-u-ca-gregory-hc-h23-nu-arab"},
		{"different numbering system", datetimeformat.Options{Calendar: String("gregory"), HourCycle: String("h23"), NumberingSystem: String("latn")}, "en-u-ca-gregory-hc-h23"},
		{"different hour cycle", datetimeformat.Options{Calendar: String("gregory"), HourCycle: String("h12"), NumberingSystem: String("arab")}, "en-u-ca-gregory-nu-arab"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.options.Hour = String("numeric")
			tt.options.TimeZone = String("UTC")
			f, err := datetimeformat.New(intltest.LocaleList(t, "en-u-ca-gregory-hc-h23-nu-arab"), tt.options)
			if err != nil {
				t.Fatal(err)
			}
			if got := f.ResolvedOptions().Locale.String(); got != tt.want {
				t.Fatalf("locale = %q, want %q", got, tt.want)
			}
		})
	}
}
