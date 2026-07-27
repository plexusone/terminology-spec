# terminology-spec

Canonical term IR for consistent handling of proper nouns, acronyms, and
standard names across TTS pronunciation, subtitle case-correction,
multi-language narration, and UI strings.

## Why

TTS engines mispronounce brand names ("AAuth" as "auth", "kubectl" as
"kub-CTL"). STT transcribers mis-case them ("oauth" instead of "OAuth").
Translators need to know whether "AI" stays "AI" in French copy or becomes
"IA". Every consumer reimplementing this logic diverges.

This package is the **single source of truth**: one Term is authored once, by
hand, per real-world usage observed. The spec records facts — it doesn't guess
pronunciation or translation from spelling rules.

## Installation

```bash
go get github.com/plexusone/terminology-spec
```

## Usage

```go
import "github.com/plexusone/terminology-spec/pkg/terminology"

// Embedded generic-industry terms (~247: AI/ML, companies, dev tools,
// languages, frameworks, cloud/infra) — no filesystem access required.
terms, _ := terminology.Builtin()

// Hand-authored canonical terms with real pronunciation/translation data
// (protocol/standard names like "AAuth", "OAuth", "AI").
curated, _ := terminology.Curated()

// Builtin merged with Curated (Curated wins on ID collision).
all, _ := terminology.All()

// Or load a project's own canonical terms/*.json directory.
custom, _ := terminology.LoadDir("terms")
```

## Exports

| Export | Use case |
|--------|----------|
| `ExportPronunciationProfile(terms)` | videoascode `VideoConfig.Pronunciations` — term → locale → spoken form, applied before TTS |
| `ExportCaseDictionary(terms)` | Subtitle case-correction — lowercase alias → canonical form |
| `ExportTranslationMap(terms, locale)` | Localization — canonical form → translated form for a given locale |

## Term structure

```json
{
  "id": "aauth",
  "canonicalForm": "AAuth",
  "category": "protocol-name",
  "scope": "aistandardsio",
  "context": "AAuth is an agent-identity protocol. Hard first 'A' — not a typo of OAuth.",
  "pronunciations": {
    "en-US": "ay auth"
  },
  "translations": {
    "zh-Hans": { "policy": "preserve" }
  }
}
```

| Field | Description |
|-------|-------------|
| `id` | Stable lowercase-kebab slug |
| `canonicalForm` | Correct written form (used for case-correction) |
| `aliases` | Alternate spellings that resolve to this term |
| `category` | Grouping: `protocol-name`, `acronym`, `common-word`, etc. |
| `scope` | Org/product owner: `aistandardsio`, `plexusone`, etc. |
| `pronunciations` | BCP-47 locale → how TTS should say it |
| `translations` | BCP-47 locale → translation policy + localized text |

## Translation policies

| Policy | Meaning |
|--------|---------|
| `preserve` | Keep canonical form as-is (default when absent) |
| `translate` | Replace with a genuinely different word ("AI" → "IA" in French) |
| `transliterate` | Re-render phonetically in the target script |
| `localized-name` | Substitute a market-specific proper name |
| `context-dependent` | No single rule; see `notes` for guidance |

## Embedded term sets

- **Builtin** (~237 terms): generic industry tech terms migrated from
  videoascode's former hardcoded case-correction map — AI/ML, companies, dev
  tools, programming languages, frameworks, cloud/infra. Case-correction only;
  no pronunciation or translation data.

- **Curated** (~10 terms): hand-authored, multi-org canonical terms with real
  pronunciation and translation data — "AAuth", "OAuth", "AI", "app", IETF
  standards. The rich counterpart to Builtin.

## Consumers

- [omnivoice-core](https://github.com/plexusone/omnivoice-core) — shared
  `Pronouncer` and `CaseCorrector` engine that wraps this spec
- [videoascode](https://github.com/grokify/videoascode) — TTS pronunciation
  dictionary and subtitle case-correction via omnivoice-core

## Schema

The JSON schema for term files is at `schema/terminology.schema.json`,
generated from the Go types via `go generate ./pkg/terminology`.

## License

MIT
