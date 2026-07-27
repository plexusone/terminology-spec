package terminology

// FilterByScope keeps only terms that apply to a consumer requesting the
// given scopes (e.g. an org name like "aistandardsio"): unscoped terms
// (Scope == "", generic/cross-org — acronyms, common words) always pass;
// scoped terms pass only if their Scope matches one of the requested
// scopes. This is how a project pulls in the shared generic layer plus its
// own org's terms without also picking up every other org's — an
// org-specific abbreviation like "PS" (this deck's shorthand for "Person
// Server") should never leak into an unrelated deck's TTS input just
// because both happened to call All().
func FilterByScope(terms []Term, scopes ...string) []Term {
	wanted := make(map[string]bool, len(scopes))
	for _, s := range scopes {
		wanted[s] = true
	}
	var filtered []Term
	for _, t := range terms {
		if t.Scope == "" || wanted[t.Scope] {
			filtered = append(filtered, t)
		}
	}
	return filtered
}
