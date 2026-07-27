// Command schema-gen generates schema/terminology.schema.json from the
// canonical Go types in pkg/terminology, per the Go-first JSON Schema
// workflow: Go structs are the source of truth, the schema is generated
// output, never hand-edited.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/plexusone/terminology-spec/pkg/terminology"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "schema-gen:", err)
		os.Exit(1)
	}
}

func run() error {
	reflector := &jsonschema.Reflector{
		DoNotReference: false,
		ExpandedStruct: true,
	}
	if err := reflector.AddGoComments("github.com/plexusone/terminology-spec", "./pkg/terminology"); err != nil {
		return fmt.Errorf("load go doc comments: %w", err)
	}
	schema := reflector.Reflect(&terminology.Term{})
	schema.Title = "Term"
	schema.Description = "A canonical terminology-spec term: pronunciation, translation policy, and case-correction data for one concept/word/name, shared across TTS, subtitles, and localized text consumers."

	out, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal schema: %w", err)
	}
	out = append(out, '\n')

	const path = "schema/terminology.schema.json"
	if err := os.WriteFile(path, out, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	fmt.Println("wrote", path)
	return nil
}
