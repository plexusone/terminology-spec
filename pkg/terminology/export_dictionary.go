package terminology

import "strings"

// CaseDictionary mirrors videoascode's tts.Dictionary JSON shape exactly
// (see grokify/videoascode pkg/tts/dictionary.go), so a file written by
// ExportCaseDictionary can be dropped straight into a consuming presentation
// repo's ./dictionaries/ directory and picked up by vac's existing
// DictionaryLoader with zero videoascode code changes.
type CaseDictionary struct {
	Name        string            `json:"name,omitempty"`
	Version     string            `json:"version,omitempty"`
	Corrections map[string]string `json:"corrections"`
}

// ExportCaseDictionary projects terms into a case-correction dictionary:
// lowercased CanonicalForm/Aliases -> CanonicalForm. videoascode applies
// this to STT-derived subtitle text via CaseCorrector, independently of
// (and complementary to) ExportPronunciationProfile — this fixes what's
// DISPLAYED, pronunciation fixes what's SPOKEN.
func ExportCaseDictionary(name, version string, terms []Term) CaseDictionary {
	dict := CaseDictionary{
		Name:        name,
		Version:     version,
		Corrections: map[string]string{},
	}
	for _, t := range terms {
		if t.CanonicalForm == "" {
			continue
		}
		dict.Corrections[strings.ToLower(t.CanonicalForm)] = t.CanonicalForm
		for _, alias := range t.Aliases {
			dict.Corrections[strings.ToLower(alias)] = t.CanonicalForm
		}
	}
	return dict
}
