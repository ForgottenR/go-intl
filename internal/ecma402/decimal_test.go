package ecma402

import (
	"errors"
	"math/big"
	"testing"

	"github.com/agentable/go-intl/internal/decimal"
)

func TestParseDecimalInputAllowsSpecialValues(t *testing.T) {
	t.Parallel()

	for _, value := range []string{"NaN", "Infinity", "-Infinity", "1.25"} {
		t.Run(value, func(t *testing.T) {
			t.Parallel()

			if _, err := ParseDecimalInput(value); err != nil {
				t.Fatalf("ParseDecimalInput(%q) error = %v, want nil", value, err)
			}
		})
	}
}

func TestParseFiniteDecimalInputRejectsMalformedAndNonFinite(t *testing.T) {
	t.Parallel()

	for _, value := range []string{"bad", "NaN", "Infinity", "-Infinity"} {
		t.Run(value, func(t *testing.T) {
			t.Parallel()

			_, err := ParseFiniteDecimalInput(value)
			if !errors.Is(err, decimal.ErrInvalidDecimal) {
				t.Fatalf("ParseFiniteDecimalInput(%q) error = %v, want ErrInvalidDecimal", value, err)
			}
		})
	}
}

func TestRequireFiniteDecimalInput(t *testing.T) {
	t.Parallel()

	if err := RequireFiniteDecimalInput(decimal.FromInt64(1)); err != nil {
		t.Fatalf("RequireFiniteDecimalInput(finite) error = %v, want nil", err)
	}
	if err := RequireFiniteDecimalInput(decimal.NaNValue); !errors.Is(err, decimal.ErrInvalidDecimal) {
		t.Fatalf("RequireFiniteDecimalInput(NaN) error = %v, want ErrInvalidDecimal", err)
	}
}

func TestParseDecimalInputRoundMVResult(t *testing.T) {
	t.Parallel()
	// Overflow midpoint: 2^1024 - 2^970. Underflow midpoint: 2^-1075.
	overflow := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 1024), new(big.Int).Lsh(big.NewInt(1), 970))
	underflow := new(big.Int).Exp(big.NewInt(5), big.NewInt(1075), nil)
	cases := []struct {
		input, form string
		negative    bool
	}{
		{"1e309", "infinity", false}, {"-1e309", "infinity", true},
		{"1e-400", "zero", false}, {"-1e-400", "zero", true},
		{overflow.String(), "infinity", false}, {"-" + overflow.String(), "infinity", true},
		{new(big.Int).Sub(overflow, big.NewInt(1)).String(), "finite", false},
		{underflow.String() + "e-1075", "zero", false},
		{"-" + underflow.String() + "e-1075", "zero", true},
		{new(big.Int).Add(underflow, big.NewInt(1)).String() + "e-1075", "finite", false},
		{"1e9999999999999999999999999", "infinity", false},
		{"-1e-999999999999999999999999", "zero", true},
		{"-0e9999999999999999999999999", "zero", true},
		{"0x" + new(big.Int).Lsh(big.NewInt(1), 1024).Text(16), "infinity", false},
	}
	for _, tc := range cases {
		got, err := ParseDecimalInput(tc.input)
		if err != nil {
			t.Errorf("ParseDecimalInput(%s): %v", tc.input, err)
			continue
		}
		form := "finite"
		if got.IsInf() {
			form = "infinity"
		}
		if got.IsZero() {
			form = "zero"
		}
		if form != tc.form || got.Negative() != tc.negative {
			t.Errorf("ParseDecimalInput(%s): %s negative=%v, want %s negative=%v", tc.input, form, got.Negative(), tc.form, tc.negative)
		}
	}
	got, err := ParseDecimalInput("9007199254740993.1")
	if err != nil || got.String() != "9007199254740993.1" {
		t.Fatalf("finite precision = %s, %v", got.String(), err)
	}
}
