package relativetimeformat

import (
	"math"
	"strings"
	"testing"

	"github.com/agentable/go-intl/internal/intltest"
)

func TestFormatProjectsPartition(t *testing.T) {
	t.Parallel()
	for _, locale := range []string{"ar", "ru", "pl", "en"} {
		for _, numeric := range []string{"always", "auto"} {
			t.Run(locale+"/"+numeric, func(t *testing.T) {
				t.Parallel()
				f, err := New(intltest.LocaleList(t, locale), Options{Numeric: new(numeric)})
				if err != nil {
					t.Fatal(err)
				}
				for _, value := range []float64{math.Copysign(0, -1), 0, -1, 1, 2, 5, 1.5, -1234.5} {
					for _, unit := range []Unit{Day, Unit("days"), Month, Second} {
						got, err := f.Format(Float(value), unit)
						if err != nil {
							t.Fatal(err)
						}
						parts, err := f.FormatToParts(Float(value), unit)
						if err != nil {
							t.Fatal(err)
						}
						var text strings.Builder
						for _, part := range parts {
							text.WriteString(part.Value)
						}
						if got != text.String() {
							t.Fatalf("Format(%v,%s)=%q, parts=%#v", value, unit, got, parts)
						}
					}
				}
			})
		}
	}
}
