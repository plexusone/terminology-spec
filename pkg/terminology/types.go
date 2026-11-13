// Package terminology defines the canonical term IR shared across
// videoascode (TTS pronunciation, subtitle case-correction, multi-language
// narration), presentation decks, and product surfaces (websites, desktop
// apps) that need consistent handling of the same proper nouns, acronyms,
// and standard names.
//
// A Term is authored once, by hand, per real-world usage observed — this
// package never derives pronunciation or translation from spelling rules.
// The "app" -> "A P P" (Chinese) and "AI" -> "IA" (French) cases only exist
// because someone said the word out loud or read a real translation; the
// spec records that fact, it doesn't guess it.
package terminology

//go:generate sh -c "cd ../.. && go run ./cmd/schema-gen"

// TranslationPolicy describes how a term's written form should be handled
// when localizing text into a given locale.
type TranslationPolicy string

const (
	// PolicyPreserve keeps the term in its canonical (usually English/Latin)
	// written form untouched, e.g. "AI" as used in-line in Chinese tech copy.
	PolicyPreserve TranslationPolicy = "preserve"
	// PolicyTranslate replaces the term with a genuinely different word in
	// the target locale, e.g. "AI" -> "IA" (French), "AI" -> "KI" (German).
	PolicyTranslate TranslationPolicy = "translate"
	// PolicyTransliterate re-renders the term phonetically in the target
	// script without translating its meaning.
	PolicyTransliterate TranslationPolicy = "transliterate"
	// PolicyLocalizedName substitutes a locale-specific proper name that
	// isn't a direct translation (e.g. a market-specific product name).
	PolicyLocalizedName TranslationPolicy = "localized-name"
	// PolicyContextDependent means no single rule applies; see Notes on the
	// TranslationEntry for the actual guidance a human translator needs.
	PolicyContextDependent TranslationPolicy = "context-dependent"
)

// TranslationEntry is the per-locale override of a term's translation
// behavior. Locale keys are BCP-47 (e.g. "fr-FR", "zh-Hans"), matching the
// locale tags already used by videoascode's transcript IR.
type TranslationEntry struct {
	Policy TranslationPolicy `json:"policy" jsonschema:"enum=preserve,enum=translate,enum=transliterate,enum=localized-name,enum=context-dependent"`
	// Text is the localized written form to use when Policy != preserve.
	// Left empty when Policy == preserve (canonicalForm is used as-is).
	Text string `json:"text,omitempty"`
	// Notes captures the human rationale/source for this override — e.g.
	// "matches the fr-FR videoascode voiceover translation" — so later
	// authors can tell a deliberate call from a placeholder.
	Notes string `json:"notes,omitempty"`
}

// Term is one canonical concept/word/name whose spoken form, written form,
// and translation behavior needs to stay consistent across every consumer
// (videoascode TTS + subtitles, presentation slide text, website content,
// desktop app UI strings).
type Term struct {
	// ID is a stable, lowercase-kebab slug, e.g. "aauth", "oauth", "ai".
	ID string `json:"id"`
	// CanonicalForm is the correct written form, e.g. "AAuth", "OAuth",
	// "AI". This is also what feeds videoascode's subtitle case-correction
	// dictionary (see ExportCaseDictionary): Aliases (lowercased) -> this.
	CanonicalForm string `json:"canonicalForm"`
	// Aliases are alternate spellings/casings that should resolve to this
	// term, e.g. ["A-Auth"] for "AAuth". CanonicalForm's own lowercase form
	// is always included implicitly — Aliases only needs the extras.
	Aliases []string `json:"aliases,omitempty"`
	// Category groups terms for filtering/scoping exports, e.g.
	// "standard-name", "protocol-name", "acronym", "common-word".
	Category string `json:"category,omitempty"`
	// Scope names the org/product this term belongs to, e.g.
	// "aistandardsio", "plexusone", "productbuildershq" — useful once the
	// corpus spans many orgs and a consumer wants only its own terms.
	Scope string `json:"scope,omitempty"`
	// Tags are freeform, e.g. ["aauth-deck"], for narrower export filters
	// than Category/Scope allow.
	Tags []string `json:"tags,omitempty"`
	// Context is a human/LLM-facing gloss: what the term means, why it's
	// distinct from lookalikes (e.g. "Not a typo of OAuth"). Feeds any
	// future LLM-context export and gives translators the disambiguation
	// they need before choosing a Translations override.
	Context string `json:"context,omitempty"`
	// Pronunciations maps BCP-47 locale -> how the term should be spoken by
	// TTS in that locale, e.g. {"en-US": "ay auth", "zh-Hans": "A P P"}.
	// Applied only to text sent to a TTS provider; CanonicalForm is what
	// appears in subtitles/manifests/UI regardless of Pronunciations.
	Pronunciations map[string]string `json:"pronunciations,omitempty"`
	// IPA maps BCP-47 locale -> the term's pronunciation in the
	// International Phonetic Alphabet, without slashes or brackets, e.g.
	// {"en-US": "lamborˈɡiːni"} for "Lamborghini". Providers that accept
	// phonemes (SSML <phoneme alphabet="ipa">, PLS lexicons) use it in
	// preference to the respelling in Pronunciations, which remains the
	// fallback for providers and models that read text only.
	IPA map[string]string `json:"ipa,omitempty"`
	// Translations maps BCP-47 locale -> the per-locale translation policy
	// and, when applicable, the actual localized text. Absence of a locale
	// key means PolicyPreserve (use CanonicalForm as-is).
	Translations map[string]TranslationEntry `json:"translations,omitempty"`
}
