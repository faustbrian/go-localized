# Architecture

The root package owns immutable text entries, exact lookup, construction,
persistent transforms, merge policy, limits, equality, hashing, and canonical
JSON. It does not import HTTP, database, configuration, validation, or wire
packages.

`match` owns standards matching and application fallback. Flat exact chains and
validated fallback graphs are separate types. Graph parent traversal delegates
to the locale layer instead of truncating strings.

Adapters point inward:

```text
http ─────────────────────┐
postgres ─────────────────┤
adapters/wire ────────────┼──> root Text
adapters/config ──────────┤
adapters/httpclient ──────┤
adapters/query ───────────┤
adapters/validation ──────┘
                   match ───> root Text
```

There is no reverse dependency from core to an adapter. Observers receive only
bounded outcome metadata after resolution and are panic-isolated. No operation
starts a goroutine or uses a cache.

Released `localized*` adapter paths own the existing implementation and public
named identities. Target-oriented `adapters/*` packages delegate to those
paths, which avoids duplicate behavior while keeping old reflection and error
identity stable.

## Locale dependency

`international/locale` owns the public locale primitive, bounded parsing,
explicit canonicalization, parent fallback, and registry provenance. This
module canonicalizes that type at every retained-value boundary. `x/text` is a
private implementation dependency used only for language matching and Unicode
normalization; it does not appear in the public API.

## Generic values

The package intentionally omits `Value[T]` in v1. Arbitrary mutable `T` values
cannot satisfy the same understandable immutable ownership guarantee without a
required clone contract. Text receives the strongest supported API.
