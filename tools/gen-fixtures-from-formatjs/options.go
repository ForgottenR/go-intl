package main

import (
	"regexp"
	"strings"
)

var (
	optionNameRE    = regexp.MustCompile(`^[A-Za-z_$][A-Za-z0-9_$]*`)
	optionLiteralRE = regexp.MustCompile(`^(?:'(?:\\.|[^\\'])*'|"(?:\\.|[^\\"])*"|-?\d+(?:_\d+)*(?:\.\d+)?(?:[eE][+-]?\d+)?|true|false)`)
)

// parseOptionsObject accepts only a fully consumed object of scalar literals.
// An omitted options argument is the empty object; unsupported expressions are
// left to the source-owned extraction audit.
func parseOptionsObject(raw string) (map[string]any, bool) {
	options := map[string]any{}
	raw, ok := skipJSTrivia(raw)
	if !ok {
		return nil, false
	}
	if raw == "" {
		return options, true
	}
	raw, ok = strings.CutPrefix(raw, "{")
	if !ok {
		return nil, false
	}
	for {
		raw, ok = skipJSTrivia(raw)
		if !ok {
			return nil, false
		}
		if rest, closed := strings.CutPrefix(raw, "}"); closed {
			rest, ok = skipJSTrivia(rest)
			if !ok || rest != "" {
				return nil, false
			}
			return options, true
		}
		name := optionNameRE.FindString(raw)
		if name == "" {
			return nil, false
		}
		raw, ok = skipJSTrivia(raw[len(name):])
		if !ok {
			return nil, false
		}
		raw, ok = strings.CutPrefix(raw, ":")
		if !ok {
			return nil, false
		}
		raw, ok = skipJSTrivia(raw)
		if !ok {
			return nil, false
		}
		literal := optionLiteralRE.FindString(raw)
		value, ok := parseOptionLiteral(literal)
		if !ok {
			return nil, false
		}
		options[name] = value
		raw, ok = skipJSTrivia(raw[len(literal):])
		if !ok {
			return nil, false
		}
		if rest, comma := strings.CutPrefix(raw, ","); comma {
			raw = rest
		} else if !strings.HasPrefix(raw, "}") {
			return nil, false
		}
	}
}

func parseOptionLiteral(literal string) (any, bool) {
	if literal == "" {
		return nil, false
	}
	switch literal[0] {
	case '\'', '"':
		return decodeJSString(literal[1 : len(literal)-1])
	}
	switch literal {
	case "true":
		return true, true
	case "false":
		return false, true
	default:
		return parseNumberLiteral(literal)
	}
}

func skipJSTrivia(raw string) (string, bool) {
	for {
		raw = strings.TrimSpace(raw)
		switch {
		case strings.HasPrefix(raw, "//"):
			end := strings.IndexAny(raw, "\r\n")
			if end < 0 {
				return "", true
			}
			raw = raw[end:]
		case strings.HasPrefix(raw, "/*"):
			_, rest, ok := strings.Cut(raw[2:], "*/")
			if !ok {
				return "", false
			}
			raw = rest
		default:
			return raw, true
		}
	}
}
