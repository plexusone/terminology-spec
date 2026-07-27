package schema

import (
	"encoding/json"
	"testing"
)

func TestTermJSON_IsValidJSONWithExpectedShape(t *testing.T) {
	var doc map[string]any
	if err := json.Unmarshal(TermJSON, &doc); err != nil {
		t.Fatalf("embedded schema is not valid JSON: %v", err)
	}

	if doc["title"] != "Term" {
		t.Errorf(`schema title = %v, want "Term"`, doc["title"])
	}

	props, ok := doc["properties"].(map[string]any)
	if !ok {
		t.Fatal("schema missing \"properties\" object")
	}
	for _, field := range []string{"id", "canonicalForm", "pronunciations", "translations"} {
		if _, ok := props[field]; !ok {
			t.Errorf("schema properties missing %q", field)
		}
	}

	required, ok := doc["required"].([]any)
	if !ok || len(required) == 0 {
		t.Fatal("schema missing non-empty \"required\" array")
	}
}
