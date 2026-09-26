package locale

import (
	"fmt"

	"golang.org/x/text/language"

	cldrlocale "github.com/agentable/go-intl/internal/cldr/locale"
	"github.com/agentable/go-intl/internal/localeid"
)

func (l Locale) Maximize() Locale {
	lang, script, region := localeid.Parts(l.tag)
	if maxLang, maxScript, maxRegion, ok := cldrlocale.MaximizeSubtags(lang, script, region); ok {
		if tag, err := localeid.ReplaceLanguageSubtags(l.tag, maxLang, maxScript, maxRegion); err == nil {
			l.tag = tag
		}
	}
	l.freeze()
	return l
}

// Minimize prefers a region over a script when both preserve the maximal LSR.
func (l Locale) Minimize() Locale {
	maximal := l.Maximize()
	lang, script, region := localeid.Parts(maximal.tag)
	for _, candidate := range [][3]string{
		{lang, "", ""},
		{lang, "", region},
		{lang, script, ""},
	} {
		trialLang, trialScript, trialRegion, ok := cldrlocale.MaximizeSubtags(candidate[0], candidate[1], candidate[2])
		if !ok || trialLang != lang || trialScript != script || trialRegion != region {
			continue
		}
		tag, err := localeid.ReplaceLanguageSubtags(l.tag, candidate[0], candidate[1], candidate[2])
		if err != nil {
			return l
		}
		l.tag = tag
		l.freeze()
		return l
	}
	return maximal
}

func mustLanguageTag(s string) language.Tag {
	tag, err := language.Parse(s)
	if err != nil {
		panic(fmt.Errorf("locale: invalid internally constructed language tag: %w", err))
	}
	return tag
}
