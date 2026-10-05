package pattern

import (
	"fmt"
	"strings"
)

type CompactPattern struct {
	Prefix, Suffix string
	Exponent       int
	HasNumber      bool
}

// ParseCompact reads the positive LDML subpattern. Quotes protect literal
// digits and semicolons; doubled apostrophes represent one literal apostrophe.
// The remaining subpattern is scanned for unclosed quotes but does not supply
// the sign: NumberFormat's sign partition owns that separately.
func ParseCompact(magnitude int, raw string) (CompactPattern, error) {
	var prefix, suffix strings.Builder
	text := &prefix
	quoted, positive := false, true
	number := ""
	for i := 0; i < len(raw); {
		c := raw[i]
		if c == '\'' {
			if i+1 < len(raw) && raw[i+1] == '\'' {
				if positive {
					text.WriteByte('\'')
				}
				i += 2
			} else {
				quoted = !quoted
				i++
			}
			continue
		}
		if !quoted && c == ';' {
			positive = false
			i++
			continue
		}
		if positive {
			if !quoted && number == "" && (c == '#' || c == '0') {
				start := i
				for i < len(raw) && strings.IndexByte("#0,.", raw[i]) >= 0 {
					i++
				}
				number = raw[start:i]
				text = &suffix
				continue
			}
			text.WriteByte(c)
		}
		i++
	}
	if quoted {
		return CompactPattern{}, fmt.Errorf("compact pattern %q has an unclosed quote: %w", raw, ErrInvalid)
	}
	out := CompactPattern{Prefix: prefix.String(), Suffix: suffix.String(), Exponent: magnitude, HasNumber: number != ""}
	if number == "0" && out.Prefix == "" && out.Suffix == "" {
		out.Exponent = 0
	} else if zeros := strings.Count(number, "0"); zeros > 0 {
		out.Exponent = magnitude + 1 - zeros
	}
	return out, nil
}
