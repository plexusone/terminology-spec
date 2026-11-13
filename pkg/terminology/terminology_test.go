package terminology

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTermFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestLoadDir_SingleObjectAndArrayFiles(t *testing.T) {
	dir := t.TempDir()
	writeTermFile(t, dir, "aauth.json", `{"id":"aauth","canonicalForm":"AAuth"}`)
	writeTermFile(t, dir, "bulk.json", `[{"id":"api","canonicalForm":"API"},{"id":"json","canonicalForm":"JSON"}]`)

	terms, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("LoadDir() error = %v", err)
	}
	if len(terms) != 3 {
		t.Fatalf("LoadDir() returned %d terms, want 3 (got: %+v)", len(terms), terms)
	}
	// sorted by ID: aauth, api, json
	if terms[0].ID != "aauth" || terms[1].ID != "api" || terms[2].ID != "json" {
		t.Errorf("LoadDir() ids = [%s, %s, %s], want [aauth, api, json]", terms[0].ID, terms[1].ID, terms[2].ID)
	}
}

func TestLoadDir_DuplicateIDError(t *testing.T) {
	dir := t.TempDir()
	writeTermFile(t, dir, "a.json", `{"id":"dup","canonicalForm":"Dup"}`)
	writeTermFile(t, dir, "b.json", `{"id":"dup","canonicalForm":"Dup2"}`)

	if _, err := LoadDir(dir); err == nil {
		t.Fatal("LoadDir() expected error for duplicate id, got nil")
	}
}

func TestLoadDir_MissingIDError(t *testing.T) {
	dir := t.TempDir()
	writeTermFile(t, dir, "a.json", `{"canonicalForm":"NoID"}`)

	if _, err := LoadDir(dir); err == nil {
		t.Fatal("LoadDir() expected error for missing id, got nil")
	}
}

func TestMerge_LaterLayerOverridesEarlier(t *testing.T) {
	base := []Term{{ID: "oauth", CanonicalForm: "OAuth"}}
	override := []Term{{ID: "oauth", CanonicalForm: "OAuth", Pronunciations: map[string]string{"en-US": "oh auth"}}}

	merged := Merge(base, override)
	if len(merged) != 1 {
		t.Fatalf("Merge() returned %d terms, want 1", len(merged))
	}
	if merged[0].Pronunciations["en-US"] != "oh auth" {
		t.Errorf("Merge() did not let the later layer win: %+v", merged[0])
	}
}

func TestMerge_UnionOfDistinctIDs(t *testing.T) {
	a := []Term{{ID: "a", CanonicalForm: "A"}}
	b := []Term{{ID: "b", CanonicalForm: "B"}}
	merged := Merge(a, b)
	if len(merged) != 2 {
		t.Fatalf("Merge() returned %d terms, want 2", len(merged))
	}
}

func TestBuiltin_LoadsEmbeddedTerms(t *testing.T) {
	terms, err := Builtin()
	if err != nil {
		t.Fatalf("Builtin() error = %v", err)
	}
	if len(terms) < 200 {
		t.Errorf("Builtin() returned %d terms, want at least 200 (migrated from videoascode's builtInCorrections)", len(terms))
	}

	found := false
	for _, tm := range terms {
		if tm.ID == "openai" && tm.CanonicalForm == "OpenAI" {
			found = true
			break
		}
	}
	if !found {
		t.Error(`Builtin() missing expected term id "openai" -> "OpenAI"`)
	}
}

func TestExportPronunciationProfile(t *testing.T) {
	terms := []Term{
		{ID: "aauth", CanonicalForm: "AAuth", Aliases: []string{"A-Auth"}, Pronunciations: map[string]string{"en-US": "ay auth"}},
		{ID: "nopron", CanonicalForm: "NoPron"},
	}
	profile := ExportPronunciationProfile(terms)

	if profile["AAuth"]["en-US"] != "ay auth" {
		t.Errorf(`profile["AAuth"]["en-US"] = %q, want "ay auth"`, profile["AAuth"]["en-US"])
	}
	if profile["A-Auth"]["en-US"] != "ay auth" {
		t.Errorf("alias should map to the same pronunciation table")
	}
	if _, ok := profile["NoPron"]; ok {
		t.Error("terms with no pronunciations should not appear in the profile")
	}
}

func TestExportIPAProfile(t *testing.T) {
	terms := []Term{
		{ID: "lamborghini", CanonicalForm: "Lamborghini", Aliases: []string{"Lambo"}, IPA: map[string]string{"en-US": "lamborˈɡiːni"}},
		{ID: "aauth", CanonicalForm: "AAuth", Pronunciations: map[string]string{"en-US": "ay auth"}},
	}
	profile := ExportIPAProfile(terms)

	if profile["Lamborghini"]["en-US"] != "lamborˈɡiːni" {
		t.Errorf(`profile["Lamborghini"]["en-US"] = %q`, profile["Lamborghini"]["en-US"])
	}
	if profile["Lambo"]["en-US"] != "lamborˈɡiːni" {
		t.Error("alias should map to the same IPA table")
	}
	if _, ok := profile["AAuth"]; ok {
		t.Error("terms with no IPA should not appear in the profile")
	}
}

func TestExportCaseDictionary(t *testing.T) {
	terms := []Term{{ID: "aauth", CanonicalForm: "AAuth", Aliases: []string{"A-Auth"}}}
	dict := ExportCaseDictionary("test", "1.0", terms)

	if dict.Corrections["aauth"] != "AAuth" {
		t.Errorf(`Corrections["aauth"] = %q, want "AAuth"`, dict.Corrections["aauth"])
	}
	if dict.Corrections["a-auth"] != "AAuth" {
		t.Errorf(`Corrections["a-auth"] (from alias) = %q, want "AAuth"`, dict.Corrections["a-auth"])
	}
}

func TestFilterByScope(t *testing.T) {
	terms := []Term{
		{ID: "api", CanonicalForm: "API"},                               // unscoped, generic
		{ID: "aauth", CanonicalForm: "AAuth", Scope: "aistandardsio"},   // matches
		{ID: "internal", CanonicalForm: "Internal", Scope: "plexusone"}, // does not match
	}

	got := FilterByScope(terms, "aistandardsio")
	if len(got) != 2 {
		t.Fatalf("FilterByScope() returned %d terms, want 2 (got: %+v)", len(got), got)
	}
	ids := map[string]bool{}
	for _, tm := range got {
		ids[tm.ID] = true
	}
	if !ids["api"] || !ids["aauth"] {
		t.Errorf("FilterByScope() = %+v, want unscoped \"api\" and matching-scope \"aauth\"", got)
	}
	if ids["internal"] {
		t.Error("FilterByScope() should not include a term scoped to a different org")
	}
}

func TestExportTranslationStrings_PreserveAndFallback(t *testing.T) {
	terms := []Term{
		{
			ID: "ai", CanonicalForm: "AI",
			Translations: map[string]TranslationEntry{
				"fr-FR": {Policy: PolicyTranslate, Text: "IA"},
			},
		},
	}

	out := ExportTranslationStrings(terms, []string{"en-US", "fr-FR", "fr-CA"})

	if out["en-US"]["ai"] != "AI" {
		t.Errorf(`en-US "ai" = %q, want "AI" (implicit preserve)`, out["en-US"]["ai"])
	}
	if out["fr-FR"]["ai"] != "IA" {
		t.Errorf(`fr-FR "ai" = %q, want "IA"`, out["fr-FR"]["ai"])
	}
	if out["fr-CA"]["ai"] != "IA" {
		t.Errorf(`fr-CA "ai" = %q, want "IA" (base-language fallback to fr-FR)`, out["fr-CA"]["ai"])
	}
}
