# Package map

[Documentation index](README.md) · [Repository entry](../README.md)

This map follows the current source modules, not a synchronized release.
Use the repository README for published installation versions; an unpublished
source major is not installable merely because its package appears here.
Internal packages and non-releasable fixtures are intentionally excluded.

## github.com/faustbrian/go-localized/v6

`localized` provides immutable UTF-8 text keyed by canonical BCP 47 language tags. Exact lookup, language matching, and application fallback are separate operations. The package does not provide catalogs, formatting, pluralization, translation loading, language detection, or global locale policy.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-localized/v6`](../error.go) | Use for immutable localized text values, exact locale lookup, bounded construction, persistent transforms, merge policy, and canonical JSON. |
| [`github.com/faustbrian/go-localized/v6/encoding`](../encoding/entries.go) | Use for the ordered entry-array encoding that preserves locale and localized-text boundaries. |
| [`github.com/faustbrian/go-localized/v6/http`](../http/accept_language.go) | Use for bounded Accept-Language parsing and selection with explicit matching semantics. |
| [`github.com/faustbrian/go-localized/v6/adapters/config`](../adapters/config/config.go) | Use to decode localized text through the Golib configuration contract. |
| [`github.com/faustbrian/go-localized/v6/adapters/httpclient`](../adapters/httpclient/httpclient.go) | Use to add Accept-Language preferences to immutable HTTP-client request specifications and select localized responses. |
| [`github.com/faustbrian/go-localized/v6/adapters/query`](../adapters/query/query.go) | Use for exact locale-aware values and predicates in API-query plans. |
| [`github.com/faustbrian/go-localized/v6/adapters/validation`](../adapters/validation/validation.go) | Use for reusable localized-text validation and explicit Unicode normalization. |
| [`github.com/faustbrian/go-localized/v6/adapters/wire`](../adapters/wire/wire.go) | Use to encode and decode localized text through bounded Golib wire formats. |
| [`github.com/faustbrian/go-localized/v6/localizedtest`](../localizedtest/localizedtest.go) | Use in tests for deterministic localized fixtures and assertions. |
| [`github.com/faustbrian/go-localized/v6/match`](../match/error.go) | Use for standards-aligned language matching and explicit application fallback plans. |
| [`github.com/faustbrian/go-localized/v6/postgres`](../postgres/postgres.go) | Use for SQL values, rows, and pgx JSONB integration backed by caller-owned database resources. |

Choose the next guide in the [documentation index](README.md)
or return to the [repository entry](../README.md) for installation and ownership.
