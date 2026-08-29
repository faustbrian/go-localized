# localized

[![CI](https://github.com/faustbrian/go-localized/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-localized/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-localized/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-localized.svg)](https://pkg.go.dev/github.com/faustbrian/go-localized)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-localized?sort=semver)](https://github.com/faustbrian/go-localized/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`localized` provides immutable UTF-8 text keyed by canonical BCP 47 language
tags. Exact lookup, language matching, and application fallback are separate
operations. The package does not provide catalogs, formatting, pluralization,
translation loading, language detection, or global locale policy.

## Install

```sh
go get github.com/faustbrian/go-localized
```

Go 1.26.6 or later is required.

## Five-minute tour

```go
englishTag, err := locale.Parse("en")
if err != nil {
    return err
}
finnishTag, _ := locale.Parse("fi")
canadianEnglish, _ := locale.Parse("en-CA")
swedishTag, _ := locale.Parse("sv")

text, err := localized.NewText(
    localized.Entry{Locale: englishTag, Text: "Hello"},
    localized.Entry{Locale: finnishTag, Text: "Hei"},
)
if err != nil {
    return err
}

english, present := text.Get(englishTag) // exact only
_ = english
_ = present

matched, err := localizedmatch.Best(text,
    localizedmatch.Preference{Locale: canadianEnglish, Weight: 1},
)
if err != nil {
    return err
}

plan, err := localizedmatch.NewFallbackPlan(
    []locale.Tag{swedishTag, englishTag}, nil, 4,
)
if err != nil {
    return err
}
fallback := plan.Resolve(text)

overlay, _ := localized.TextFromMap(map[string]string{"en": "Hi"})
merged, err := text.Merge(overlay, localized.RightWins)
if err != nil {
    return err
}

canonicalJSON, err := localized.EncodeJSON(merged)
_ = matched
_ = fallback
_ = canonicalJSON
```

For SQL and pgx, use `postgres.NewText(value)` and
`postgres.JSONBCodec()`. See the [quickstart](docs/quickstart.md) for complete
construction, fallback, merge, JSON, and PostgreSQL examples.

## Guarantees

- caller maps, entry slices, rows, iterators, and encoded bytes do not alias
  retained state;
- locale keys use canonical `international/locale.Tag` identity;
- missing and present-empty are distinguished by every lookup result;
- iteration and canonical encoding are lexically deterministic;
- fallback never inserts an invented translation;
- parser, locale, text, matching, fallback, merge, and telemetry work is
  bounded;
- production code has no mutable globals, unsafe, cgo, `go:linkname`, cache,
  goroutine, registry refresh, or process-global locale;
- package-generated errors and events never include localized content.

## Documentation

Start at the [documentation index](docs/README.md). The normative behavior is
in [semantics](docs/semantics.md), the complete public surface in the
[API reference](docs/api.md), and operational constraints in
[security](docs/security.md) and [performance](docs/performance.md).

## Development

`make check` runs the complete local gate stack through the pinned
`go-library-tools` CLI. Hosted workflows mirror these commands, but local
development does not depend on a remote branch or CI run. PostgreSQL-backed
checks require an explicitly supplied disposable database:

```sh
POSTGRES_URL='postgres://postgres:postgres@127.0.0.1:5432/localized?sslmode=disable' \
  make check
```

The PostgreSQL version matrix is a CI operation. It uses isolated ephemeral
containers and verifies PostgreSQL 14 through 18 without using ambient
services.

## License

MIT. See [LICENSE](LICENSE), [NOTICE](NOTICE), and
[THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).
