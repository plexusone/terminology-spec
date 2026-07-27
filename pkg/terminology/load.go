package terminology

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"
)

// LoadDir loads every terms/*.json file in dir from the real filesystem. See
// LoadFS for the file-format rules.
func LoadDir(dir string) ([]Term, error) {
	return LoadFS(os.DirFS(dir), ".")
}

// LoadFS loads every *.json file under dir in fsys and returns the combined
// terms sorted by ID for deterministic export output.
//
// Each file may hold either a single Term object (the one-file-per-entity
// convention used for rich, individually-curated terms like "AAuth") or a
// JSON array of Term objects (for bulk collections like the embedded
// builtin generic-industry-term layer, where hundreds of terms carry little
// more than an ID and CanonicalForm).
func LoadFS(fsys fs.FS, dir string) ([]Term, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("read terms dir %s: %w", dir, err)
	}

	var terms []Term
	seen := map[string]string{} // id -> source file, to catch duplicate IDs
	for _, entry := range entries {
		if entry.IsDir() || path.Ext(entry.Name()) != ".json" {
			continue
		}
		filePath := path.Join(dir, entry.Name())
		b, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", filePath, err)
		}

		fileTerms, err := parseTermsFile(b)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", filePath, err)
		}
		for _, t := range fileTerms {
			if t.ID == "" {
				return nil, fmt.Errorf("%s: term missing required \"id\" field", filePath)
			}
			if prior, ok := seen[t.ID]; ok {
				return nil, fmt.Errorf("duplicate term id %q in %s (already defined in %s)", t.ID, filePath, prior)
			}
			seen[t.ID] = filePath
			terms = append(terms, t)
		}
	}

	sort.Slice(terms, func(i, j int) bool { return terms[i].ID < terms[j].ID })
	return terms, nil
}

func parseTermsFile(b []byte) ([]Term, error) {
	trimmed := jsonFirstToken(b)
	if trimmed == '[' {
		var terms []Term
		if err := json.Unmarshal(b, &terms); err != nil {
			return nil, err
		}
		return terms, nil
	}
	var t Term
	if err := json.Unmarshal(b, &t); err != nil {
		return nil, err
	}
	return []Term{t}, nil
}

func jsonFirstToken(b []byte) byte {
	for _, c := range b {
		switch c {
		case ' ', '\t', '\n', '\r':
			continue
		default:
			return c
		}
	}
	return 0
}
