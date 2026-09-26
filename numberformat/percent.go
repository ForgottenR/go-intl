package numberformat

import (
	"strings"

	cldrnumber "github.com/agentable/go-intl/internal/cldr/number"
)

type percentPatternSet struct {
	unsigned, negative, positive numberAffixPattern
}

func percentPatternsForNumberFormat(loc cldrnumber.Locale, opts ResolvedOptions, symbols cldrnumber.NumberSymbols) percentPatternSet {
	if opts.Style != PercentStyle {
		return percentPatternSet{}
	}
	positive, negative, hasNegative := strings.Cut(loc.PercentPattern(opts.NumberingSystem), ";")
	if !hasNegative {
		negative = "-" + positive
	}
	plus := "+" + positive
	if strings.Contains(negative, "-") {
		plus = strings.ReplaceAll(negative, "-", "+")
	}
	return percentPatternSet{
		unsigned: compilePercentPattern(positive, symbols),
		negative: compilePercentPattern(negative, symbols),
		positive: compilePercentPattern(plus, symbols),
	}
}

func (p percentPatternSet) append(parts []Part) []Part {
	sign, number := splitLeadingSign(parts)
	switch sign.Type {
	case PartMinusSign:
		return p.negative.append(number)
	case PartPlusSign:
		return p.positive.append(number)
	default:
		return p.unsigned.append(number)
	}
}

func compilePercentPattern(pattern string, symbols cldrnumber.NumberSymbols) numberAffixPattern {
	start, end := numberPatternBounds(pattern)
	if start < 0 {
		return numberAffixPattern{prefix: compilePercentLiteral(pattern, symbols)}
	}
	return numberAffixPattern{
		prefix: compilePercentLiteral(pattern[:start], symbols),
		suffix: compilePercentLiteral(pattern[end:], symbols),
	}
}

func compilePercentLiteral(text string, symbols cldrnumber.NumberSymbols) []Part {
	var out []Part
	var literal strings.Builder
	flush := func() { out = appendLiteral(out, literal.String()); literal.Reset() }
	for _, r := range text {
		var part Part
		switch r {
		case '%':
			part = Part{Type: PartPercentSign, Value: symbols.Percent}
		case '-':
			part = Part{Type: PartMinusSign, Value: symbols.Minus}
		case '+':
			part = Part{Type: PartPlusSign, Value: symbols.Plus}
		default:
			literal.WriteRune(r)
			continue
		}
		flush()
		out = appendBidiSymbol(out, part)
	}
	flush()
	return out
}
