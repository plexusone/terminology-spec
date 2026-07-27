package terminology

import "embed"

//go:embed embedded/builtin
var builtinFS embed.FS

//go:embed embedded/curated
var curatedFS embed.FS

// Builtin returns the generic industry-technology terms shipped with
// terminology-spec (proper-noun/acronym case corrections migrated from
// videoascode's former hardcoded builtInCorrections map: AI/ML terms,
// companies, dev tools, languages, frameworks, cloud/infra, etc.).
//
// These carry no pronunciation or translation data — they exist so any
// consumer can get sensible default case-correction behavior as a Go
// dependency, without maintaining its own copy of the same ~250 terms.
// Org/product-specific terms with real pronunciation/translation data
// belong in Curated.
func Builtin() ([]Term, error) {
	return LoadFS(builtinFS, "embedded/builtin")
}

// Curated returns the hand-authored, multi-org canonical terms (protocol
// and standard names like "AAuth"/"OAuth", plus cross-cutting demonstrative
// entries like "AI"/"app") — the rich, individually-curated counterpart to
// Builtin. This is the single source consumers like videoascode and
// presentation decks should pull real pronunciation/translation data from,
// rather than each maintaining their own local copy.
func Curated() ([]Term, error) {
	return LoadFS(curatedFS, "embedded/curated")
}

// All returns the full canonical term set: Builtin merged with Curated,
// where Curated wins on any ID collision (e.g. "oauth" gets Builtin's bare
// case-correction plus Curated's real pronunciation).
func All() ([]Term, error) {
	builtin, err := Builtin()
	if err != nil {
		return nil, err
	}
	curated, err := Curated()
	if err != nil {
		return nil, err
	}
	return Merge(builtin, curated), nil
}
