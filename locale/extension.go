package locale

import (
	"maps"
	"slices"

	"github.com/agentable/go-intl/internal/localeid"
)

func (l *Locale) readExtensions(ext localeid.UnicodeExtension) {
	l.ext.attributes = ext.Attributes()
	for _, keyword := range ext.Keywords() {
		l.ext.setKeyword(keyword.Key, keyword.Value)
	}
}

func (e extensions) empty() bool {
	return len(e.attributes) == 0 && len(e.keywords) == 0
}

func (l Locale) unicodeExtensionKeywords() []localeid.UnicodeKeyword {
	keys := slices.Sorted(maps.Keys(l.ext.keywords))
	out := make([]localeid.UnicodeKeyword, len(keys))
	for i, key := range keys {
		out[i] = localeid.UnicodeKeyword{Key: key, Value: l.ext.keywords[key]}
	}
	return out
}

func (e *extensions) setKeyword(key, value string) {
	if e.keywords == nil {
		e.keywords = make(map[string]string)
	}
	if value == "true" {
		value = ""
	}
	e.keywords[key] = value
}
