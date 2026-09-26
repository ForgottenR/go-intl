package ecma402

import (
	"math"
	"math/big"
	"strings"

	"github.com/agentable/go-intl/internal/decimal"
)

var (
	// RoundMVResult ties choose infinity at 2^1024 - 2^970 and zero at 2^-1075.
	stringOverflow  = decimal.FromBigInt(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 1024), new(big.Int).Lsh(big.NewInt(1), 970)))
	stringUnderflow = decimal.New(false, new(big.Int).Exp(big.NewInt(5), big.NewInt(1075), nil), -1075)
)

// parseFiniteIntlDecimal receives a validated decimal literal. Classify huge
// exponents before constructing a backend value; work follows input length.
func parseFiniteIntlDecimal(value string) (decimal.Decimal, error) {
	negative := value[0] == '-'
	value = strings.TrimLeft(value, "+-")
	mantissa := value
	exponent := new(big.Int)
	if i := strings.IndexAny(value, "eE"); i >= 0 {
		mantissa = value[:i]
		exponent.SetString(value[i+1:], 10)
	}
	if i := strings.IndexByte(mantissa, '.'); i >= 0 {
		exponent.Sub(exponent, big.NewInt(int64(len(mantissa)-i-1)))
		mantissa = mantissa[:i] + mantissa[i+1:]
	}
	digits := strings.TrimLeft(mantissa, "0")
	if digits == "" {
		return decimal.WithSign(decimal.Zero, negative), nil
	}
	trimmed := strings.TrimRight(digits, "0")
	exponent.Add(exponent, big.NewInt(int64(len(digits)-len(trimmed))))
	digits = trimmed
	magnitude := new(big.Int).Add(exponent, big.NewInt(int64(len(digits)-1)))
	if magnitude.Cmp(big.NewInt(308)) > 0 {
		return decimal.WithSign(decimal.PosInfinity, negative), nil
	}
	if magnitude.Cmp(big.NewInt(-324)) < 0 {
		return decimal.WithSign(decimal.Zero, negative), nil
	}
	// A finite result may retain all significant input digits; the exponent
	// storage limit is independent of the mathematical Number range.
	if !exponent.IsInt64() || exponent.Int64() < math.MinInt32 || exponent.Int64() > math.MaxInt32 {
		return decimal.Decimal{}, decimal.ErrInvalidDecimal
	}
	coefficient, _ := new(big.Int).SetString(digits, 10)
	bounded := exponent.Int64()                             // checked against int32 limits above
	d := decimal.New(negative, coefficient, int32(bounded)) // #nosec G115 -- bounded above
	return normalizeStringValue(d), nil
}

func normalizeStringValue(value decimal.Decimal) decimal.Decimal {
	magnitude := decimal.Abs(value)
	if magnitude.Cmp(stringOverflow) >= 0 {
		return decimal.WithSign(decimal.PosInfinity, value.Negative())
	}
	if magnitude.Cmp(stringUnderflow) <= 0 {
		return decimal.WithSign(decimal.Zero, value.Negative())
	}
	return value
}
