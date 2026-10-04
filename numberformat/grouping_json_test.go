package numberformat

import (
	"encoding/json/v2"
	"testing"

	"github.com/agentable/go-intl/internal/intltest"
)

func TestResolvedGroupingJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mode *string
		want any
		text string
	}{
		{"default", nil, "auto", "1,000"},
		{"disabled", new("false"), false, "1000"},
		{"auto", new("auto"), "auto", "1,000"},
		{"always", new("always"), "always", "1,000"},
		{"min2", new("min2"), "min2", "1000"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f, err := New(intltest.LocaleList(t, "en"), Options{UseGrouping: tc.mode})
			if err != nil {
				t.Fatal(err)
			}
			data, err := json.Marshal(f.ResolvedOptions())
			if err != nil {
				t.Fatal(err)
			}
			var record map[string]any
			if err := json.Unmarshal(data, &record); err != nil {
				t.Fatal(err)
			}
			if got, present := record["useGrouping"]; !present || got != tc.want {
				t.Fatalf("useGrouping JSON = %T(%v), present=%v, want %T(%v)", got, got, present, tc.want, tc.want)
			}
			if got := f.Format(Int(1000)); got != tc.text {
				t.Fatalf("Format(1000) = %q, want %q", got, tc.text)
			}
		})
	}
}
