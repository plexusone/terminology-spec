package terminology

import "maps"

// ExportPronunciationProfile projects terms into the exact shape videoascode
// consumes: VideoConfig.Pronunciations (map[string]map[string]string), a
// term -> BCP-47 locale -> spoken-form dictionary applied only to text sent
// to TTS. See grokify/videoascode pkg/tts/pronunciation.go.
//
// Both CanonicalForm and every Alias become separate keys mapping to the
// same locale table, so the term is matched however it's actually spelled
// in narration text.
func ExportPronunciationProfile(terms []Term) map[string]map[string]string {
	profile := map[string]map[string]string{}
	for _, t := range terms {
		if len(t.Pronunciations) == 0 {
			continue
		}
		locales := make(map[string]string, len(t.Pronunciations))
		maps.Copy(locales, t.Pronunciations)
		profile[t.CanonicalForm] = locales
		for _, alias := range t.Aliases {
			profile[alias] = locales
		}
	}
	return profile
}
