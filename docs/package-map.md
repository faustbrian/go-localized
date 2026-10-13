# Package map

[Documentation index](README.md) · [Repository entry](../README.md)

This map follows the current source modules, not a synchronized release.
Use the repository README for published installation versions; an unpublished
source major is not installable merely because its package appears here.
Internal packages and non-releasable fixtures are intentionally excluded.

## github.com/faustbrian/go-localized/v5

`localized` provides immutable UTF-8 text keyed by canonical BCP 47 language tags. Exact lookup, language matching, and application fallback are separate operations. The package does not provide catalogs, formatting, pluralization, translation loading, language detection, or global locale policy.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-localized/v5`](../error.go) | Use for immutable localized text values, exact locale lookup, bounded construction, persistent transforms, merge policy, and canonical JSON. |
| [`github.com/faustbrian/go-localized/v5/encoding`](../encoding/entries.go) | Use for the ordered entry-array encoding that preserves locale and localized-text boundaries. |
| [`github.com/faustbrian/go-localized/v5/http`](../http/accept_language.go) | Use for bounded Accept-Language parsing and selection with explicit matching semantics. |
| [`github.com/faustbrian/go-localized/v5/adapters/config`](../adapters/config/doc.go) | Use to decode localized text through the Golib configuration contract. |
| [`github.com/faustbrian/go-localized/v5/adapters/httpclient`](../adapters/httpclient/doc.go) | Use to add Accept-Language preferences to immutable HTTP-client request specifications and select localized responses. |
| [`github.com/faustbrian/go-localized/v5/adapters/query`](../adapters/query/doc.go) | Use for exact locale-aware values and predicates in API-query plans. |
| [`github.com/faustbrian/go-localized/v5/adapters/validation`](../adapters/validation/doc.go) | Use for reusable localized-text validation and explicit Unicode normalization. |
| [`github.com/faustbrian/go-localized/v5/adapters/wire`](../adapters/wire/doc.go) | Use to encode and decode localized text through bounded Golib wire formats. |
| [`github.com/faustbrian/go-localized/v5/localizedconfig`](../localizedconfig/config.go) | Compatibility import for adapters/config. |
| [`github.com/faustbrian/go-localized/v5/localizedhttpclient`](../localizedhttpclient/httpclient.go) | Compatibility import for adapters/httpclient. |
| [`github.com/faustbrian/go-localized/v5/localizedquery`](../localizedquery/query.go) | Compatibility import for adapters/query. |
| [`github.com/faustbrian/go-localized/v5/localizedtest`](../localizedtest/localizedtest.go) | Use in tests for deterministic localized fixtures and assertions. |
| [`github.com/faustbrian/go-localized/v5/localizedvalidation`](../localizedvalidation/validation.go) | Compatibility import for adapters/validation. |
| [`github.com/faustbrian/go-localized/v5/localizedwire`](../localizedwire/wire.go) | Compatibility import for adapters/wire. |
| [`github.com/faustbrian/go-localized/v5/match`](../match/error.go) | Use for standards-aligned language matching and explicit application fallback plans. |
| [`github.com/faustbrian/go-localized/v5/postgres`](../postgres/postgres.go) | Use for SQL values, rows, and pgx JSONB integration backed by caller-owned database resources. |

Choose the next guide in the [documentation index](README.md)
or return to the [repository entry](../README.md) for installation and ownership.
