package terminology

import "sort"

// Merge composes term layers in precedence order: later layers override
// earlier ones by ID. This is the generic layering primitive behind, e.g.,
// videoascode's dictionary composition — embedded Builtin() terms as the
// base layer, with an org's or project's own terms/ directory layered on
// top.
func Merge(layers ...[]Term) []Term {
	byID := map[string]Term{}
	var order []string
	for _, layer := range layers {
		for _, t := range layer {
			if _, exists := byID[t.ID]; !exists {
				order = append(order, t.ID)
			}
			byID[t.ID] = t
		}
	}
	merged := make([]Term, len(order))
	for i, id := range order {
		merged[i] = byID[id]
	}
	sort.Slice(merged, func(i, j int) bool { return merged[i].ID < merged[j].ID })
	return merged
}
