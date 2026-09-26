package ecma402

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"unicode"

	"github.com/agentable/go-intl/internal/decimal"
)

const (
	// ExpectedDecimalValue describes ECMA-402 decimal-string bridge inputs that
	// allow NumberFormat's special numeric values.
	ExpectedDecimalValue = "an Intl numeric string (decimal, binary, octal, hexadecimal, or Infinity), or NaN"
	// ExpectedFiniteNumericValue describes ECMA-402 method inputs that reject
	// NaN and infinities.
	ExpectedFiniteNumericValue = "a finite numeric value"
)

var decimalLiteral = regexp.MustCompile(`^[+-]?(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[eE][+-]?[0-9]+)?$`)

// intlSpace is ECMAScript WhiteSpace plus LineTerminator, not Go's wider IsSpace.
func intlSpace(r rune) bool {
	return r == '\t' || r == '\v' || r == '\f' || r == '\n' || r == '\r' ||
		r == '\u2028' || r == '\u2029' || r == '\ufeff' || unicode.Is(unicode.Zs, r)
}

// ParseDecimalInput parses an ECMA-402 decimal-string bridge value.
func ParseDecimalInput(value string) (decimal.Decimal, error) {
	value = strings.TrimFunc(value, intlSpace)
	if value == "" {
		return decimal.Zero, nil
	}
	switch value {
	case "NaN", "Infinity", "+Infinity", "-Infinity":
		return decimal.ParseString(value)
	}
	if len(value) > 2 && value[0] == '0' {
		base := 0
		switch value[1] {
		case 'b', 'B':
			base = 2
		case 'o', 'O':
			base = 8
		case 'x', 'X':
			base = 16
		}
		if base != 0 {
			digits := value[2:]
			if digits[0] != '+' && digits[0] != '-' && !strings.ContainsRune(digits, '_') {
				if n, ok := new(big.Int).SetString(digits, base); ok {
					return normalizeStringValue(decimal.FromBigInt(n)), nil
				}
			}
			return decimal.Decimal{}, decimal.ErrInvalidDecimal
		}
	}
	if !decimalLiteral.MatchString(value) {
		return decimal.Decimal{}, decimal.ErrInvalidDecimal
	}
	return parseFiniteIntlDecimal(value)
}

// ParseFiniteDecimalInput parses a decimal-string bridge value that must be
// finite at the owning operation boundary.
func ParseFiniteDecimalInput(value string) (decimal.Decimal, error) {
	d, err := ParseDecimalInput(value)
	if err != nil {
		return decimal.Decimal{}, err
	}
	if err := RequireFiniteDecimalInput(d); err != nil {
		return decimal.Decimal{}, err
	}
	return d, nil
}

// RequireFiniteDecimalInput rejects NaN and infinities at ECMA-402 method
// boundaries where the native operation would throw on non-finite input.
func RequireFiniteDecimalInput(value decimal.Decimal) error {
	if value.IsFinite() {
		return nil
	}
	return fmt.Errorf("ecma402: non-finite decimal input %q: %w", value.String(), decimal.ErrInvalidDecimal)
}

// InvalidDecimalValueError records a decimal-string bridge value failure with
// the shared ECMA-402 decimal-input guidance.
func InvalidDecimalValueError(owner, name, value string, err error) error {
	return InvalidValueErrorExpected(owner, name, value, "", ExpectedDecimalValue, err)
}

// InvalidFiniteNumericValueError records a finite numeric value failure with
// the shared ECMA-402 finite-input guidance.
func InvalidFiniteNumericValueError(owner, name, value, loc string, err error) error {
	return InvalidValueErrorExpected(owner, name, value, loc, ExpectedFiniteNumericValue, err)
}
