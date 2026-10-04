package listformat

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/agentable/go-intl/internal/ecma402"
)

type listContext uint8

const (
	listContextNone listContext = iota
	listContextSpanishE
	listContextSpanishU
	listContextHebrewVavDash
)

type contextualListTemplate struct {
	pattern   ecma402.Pattern
	alternate ecma402.Pattern
	context   listContext
}

func compileContextualListTemplate(text, dataLocale string) contextualListTemplate {
	template := contextualListTemplate{pattern: compileListTemplate(text)}
	lang, _, _ := strings.Cut(dataLocale, "-")
	if lang == "es" {
		switch text {
		case "{0} y {1}":
			template.context = listContextSpanishE
			template.alternate = compileListTemplate("{0} e {1}")
		case "{0} o {1}":
			template.context = listContextSpanishU
			template.alternate = compileListTemplate("{0} u {1}")
		}
	} else if lang == "he" && text == "{0} ו{1}" {
		template.context = listContextHebrewVavDash
		template.alternate = compileListTemplate("{0} ו-{1}")
	}
	return template
}

func (p contextualListTemplate) forElement(next string) ecma402.Pattern {
	if p.context.matches(next) {
		return p.alternate
	}
	return p.pattern
}

// These finite prefix rules mirror ICU listformatter.cpp's shouldChangeToE/U.
// Selection uses the original element and never edits caller-owned text.
func (c listContext) matches(text string) bool {
	if text == "" {
		return false
	}
	switch c {
	case listContextSpanishE:
		if text[0] == 'i' || text[0] == 'I' {
			return true
		}
		return (text[0] == 'h' || text[0] == 'H') && len(text) > 1 &&
			(text[1] == 'i' || text[1] == 'I') &&
			(len(text) == 2 || strings.IndexByte("aAeE", text[2]) < 0)
	case listContextSpanishU:
		if text[0] == 'o' || text[0] == 'O' || text[0] == '8' {
			return true
		}
		if (text[0] == 'h' || text[0] == 'H') && len(text) > 1 && (text[1] == 'o' || text[1] == 'O') {
			return true
		}
		return strings.HasPrefix(text, "11") && (len(text) == 2 || text[2] == ' ')
	case listContextHebrewVavDash:
		first, _ := utf8.DecodeRuneInString(text)
		return !unicode.Is(unicode.Hebrew, first)
	case listContextNone:
	}
	return false
}
