# localized

[![CI](https://github.com/faustbrian/go-localized/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-localized/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-localized/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-risk_selected-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-optional_local%2Fselected_scheduled-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-localized/v6.svg)](https://pkg.go.dev/github.com/faustbrian/go-localized/v6)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-localized?sort=semver)](https://github.com/faustbrian/go-localized/releases)
[![Go](https://img.shields.io/badge/go-1.27.0-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`localized` provides immutable UTF-8 text keyed by canonical BCP 47 language
tags. Exact lookup, language matching, and application fallback are separate
operations. The package does not provide catalogs, formatting, pluralization,
translation loading, language detection, or global locale policy.

## Install

The root `/v6` module composes API Query v5, International v5, HTTP Client v3,
Validation v3 and Wire v3 nominal identities. Select its versioned tag when adopting this composition;
published tags, not the source on main, establish release availability.
See the [migration guide](docs/migration.md) for coordinated import changes.

```sh
go get github.com/faustbrian/go-localized/v6@v6.0.0
```

Go 1.27.0 or later is required. Development and CI use Go 1.27.2
to include standard-library security fixes. Rebuild applications with a
patched compiler to receive those fixes.

The published v4 module selects International v3 and API Query v4 nominal types while
retaining Validation v2. Migrate locale and query imports together as described
in the [migration guide](docs/migration.md).
Existing v1–v3 releases remain independently available; publication is established
by the repository's public tags and releases.

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

Optional Golib integrations use target-oriented imports under `adapters/`:
`adapters/config`, `adapters/httpclient`, `adapters/query`,
`adapters/validation`, and `adapters/wire`. Version 6 removed the former
`localizedconfig`, `localizedhttpclient`, `localizedquery`,
`localizedvalidation`, and `localizedwire` paths; the v5 line keeps them.

## Guarantees

- caller maps, entry slices, rows, iterators, and encoded bytes do not alias
  retained state;
- locale keys use canonical `international/locale.Tag` identity;
- missing and present-empty are distinguished by every lookup result;
- iteration and canonical encoding are lexically deterministic;
- fallback never inserts an invented translation;
- parser bytes, entry content and candidate/graph counts have explicit limits;
  these are not total heap, callback-time or cancellation guarantees;
- production code has no mutable globals, unsafe, cgo, `go:linkname`, cache,
  goroutine, registry refresh, or process-global locale;
- package-generated errors and events never include localized content.

## Documentation

Start at the [documentation index](docs/README.md). The normative behavior is
in [semantics](docs/semantics.md), the complete public surface in the
[API reference](docs/api.md), and operational constraints in
[security](docs/security.md) and [performance](docs/performance.md).
Observable standards interpretations are governed by the
[specification decision register](docs/specification-decisions.md).
Shared construction, ownership, lifecycle, and composition expectations are in
the versioned [Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.6.1/docs/ecosystem/README.md)
and its [Foundations family](https://github.com/faustbrian/go-library-tools/blob/v1.6.1/docs/ecosystem/design-language.md#package-families-and-selection).

## Development

`make check` and `make ci` run bounded local checks through the public pinned
`go-library-tools` CLI and retain the manifest's Docker-backed PostgreSQL service
and integration tag. Use direct untagged Go tests for database-free development,
and `make check-full` for explicit full assurance.
Pull requests run the bounded local CI tier; scheduled and release workflows
retain the relevant aggregate and release checks. Local
development does not depend on a remote branch or CI run. Direct PostgreSQL
tests are a separate route with an explicitly supplied disposable database:

```sh
POSTGRES_URL='postgres://postgres:postgres@127.0.0.1:5432/localized?sslmode=disable' \
  go test -count=1 -tags=integration ./postgres
```

The explicit `make -f verification/package.mk postgres-matrix` diagnostic uses
isolated ephemeral containers and verifies PostgreSQL 14 through 18 without
using ambient services. Ordinary local checks do not launch that matrix.

## License

MIT. See [LICENSE](LICENSE), [NOTICE](NOTICE), and
[THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).

## Documentation navigation

[Complete documentation map](docs/README.md#complete-documentation-map) · [Current source package map](docs/package-map.md)

Use the [versioned ecosystem entry](https://github.com/faustbrian/go-library-tools/blob/v1.6.1/docs/ecosystem/README.md) for package
selection. The [frozen compatibility cohort](https://github.com/faustbrian/go-library-tools/blob/v1.8.4/docs/ecosystem/compatibility-sets.md)
is separate from current package releases and unpublished source changes.
Adopting them does not silently update that cohort.
