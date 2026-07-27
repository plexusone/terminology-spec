package terminology

import "strings"

// ExportTranslationStrings projects terms into a flat locale -> term ID ->
// resolved text resource bundle. This targets no specific i18n library on
// purpose: none of the current consumers (plexusone.dev, aistandards.io,
// productbuildershq.com, the visionstudio Electron app) has adopted one yet,
// and a plain nested-JSON resource is trivially consumable by any of them
// once they do.
//
// locales is the set of BCP-47 locales to resolve output for. Resolution
// per term or locale, in order:
//  1. An exact Translations[locale] entry: PolicyPreserve -> CanonicalForm,
//     otherwise -> entry.Text (falling back to CanonicalForm if Text is
//     empty, which should not happen for well-formed data).
//  2. A Translations entry for the same base language (e.g. requested
//     "fr-CA" falls back to an authored "fr-FR" entry) — same rule as
//     videoascode's transcript locale fallback (specific region -> base
//     language -> default -> original term).
//  3. CanonicalForm, unchanged (implicit PolicyPreserve).
func ExportTranslationStrings(terms []Term, locales []string) map[string]map[string]string {
	out := make(map[string]map[string]string, len(locales))
	for _, locale := range locales {
		resolved := make(map[string]string, len(terms))
		for _, t := range terms {
			resolved[t.ID] = resolveTranslation(t, locale)
		}
		out[locale] = resolved
	}
	return out
}

func resolveTranslation(t Term, locale string) string {
	if entry, ok := t.Translations[locale]; ok {
		return applyPolicy(t, entry)
	}
	base := baseLanguage(locale)
	for candidate, entry := range t.Translations {
		if baseLanguage(candidate) == base {
			return applyPolicy(t, entry)
		}
	}
	return t.CanonicalForm
}

func applyPolicy(t Term, entry TranslationEntry) string {
	if entry.Policy == PolicyPreserve || entry.Text == "" {
		return t.CanonicalForm
	}
	return entry.Text
}

func baseLanguage(locale string) string {
	base, _, _ := strings.Cut(locale, "-")
	return base
}
