package locale

import "testing"

func TestCanonicalAliasPreservesSuffix(t *testing.T) {
	t.Parallel()

	const suffix = "-Latn-GH-emodeng-t-en-us-u-ca-gregory-x-private"
	loc, err := Parse("twi" + suffix)
	if err != nil {
		t.Fatal(err)
	}
	if got := loc.String(); got != "ak"+suffix {
		t.Fatalf("Parse identity = %q, want %q", got, "ak"+suffix)
	}
	loc, err = New("en"+suffix, Options{Language: new("TWI")})
	if err != nil {
		t.Fatal(err)
	}
	if got := loc.String(); got != "ak"+suffix {
		t.Fatalf("New language identity = %q, want %q", got, "ak"+suffix)
	}
	if got := loc.Maximize().String(); got != "ak"+suffix {
		t.Fatalf("Maximize = %q, want retained suffix", got)
	}
	if got := loc.Minimize().String(); got != "ak-emodeng-t-en-us-u-ca-gregory-x-private" {
		t.Fatalf("Minimize = %q, want retained suffix", got)
	}
}

func TestNorwegianLanguageIdentity(t *testing.T) {
	t.Parallel()
	const suffix = "-Latn-NO-emodeng-t-en-us-u-ca-gregory-x-private"
	for _, base := range []string{"no", "nb"} {
		t.Run(base, func(t *testing.T) {
			t.Parallel()
			for _, input := range []string{base, base + "-NO", base + suffix} {
				loc, err := Parse(input)
				if err != nil {
					t.Fatal(err)
				}
				if got := loc.String(); got != input {
					t.Errorf("Parse(%q) = %q", input, got)
				}
			}
			loc, err := New("en"+suffix, Options{Language: new(base)})
			if err != nil {
				t.Fatal(err)
			}
			if got := loc.String(); got != base+suffix {
				t.Errorf("New language identity = %q, want %q", got, base+suffix)
			}
			if got := loc.Maximize().String(); got != base+suffix {
				t.Errorf("Maximize = %q, want %q", got, base+suffix)
			}
		})
	}
}
