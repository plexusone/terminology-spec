// Package schema embeds the generated JSON Schema for terminology-spec's
// canonical Term type, for consumers that need it at runtime (e.g. to
// validate a terms/*.json file without depending on pkg/terminology's Go
// types). The schema itself is generated — see cmd/schema-gen — never
// hand-edited.
package schema

import _ "embed"

//go:embed terminology.schema.json
var TermJSON []byte
