package localematcher

import (
	"strconv"
	"testing"

	cldrlocale "github.com/agentable/go-intl/internal/cldr/locale"
)

func BenchmarkMatcherRequests(b *testing.B) {
	for _, mode := range []string{"exact", "repeated", "distinct"} {
		b.Run(mode, func(b *testing.B) {
			m := NewMatcher([]string{"en", "nb", "de"}, cldrlocale.Maximize)
			b.ReportAllocs()
			i := 0
			for b.Loop() {
				request := "en"
				if mode == "repeated" {
					request = "nn"
				}
				if mode == "distinct" {
					request = "nn-x-" + strconv.Itoa(i)
					i++
				}
				m.Match([]string{request}, "en", AlgorithmBestFit)
			}
		})
	}
}
