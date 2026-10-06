# Migration guide

## Prepared v5 successor composition

The current source prepares `github.com/faustbrian/go-localized/v5` at the
repository root; it is not yet published. Move all owning Localized imports
to `/v5`, locale values to `github.com/faustbrian/go-international/v4/locale`
v4.0.0, HTTP request specs and layers to `github.com/faustbrian/go-http-client/v2`
v2.0.0, and codec options to `github.com/faustbrian/go-wire/v3` v3.0.0. Both
canonical `adapters/*` and retained `localized*` packages select these same
nominal types. Different majors cannot be mixed in a function signature.

API Query v4.0.0, Validation v2.0.0 and Config v2.0.0 remain selected. No
localized format, exact-presence, matching, fallback, ownership or alias
semantics change is intended. Existing public v1–v4 tags and API projections
remain available; frozen earlier-major consumers need not migrate. Select
published dependencies without sibling replacements. Localized v5 release,
public consumer and maintained Tools adoption are separate pending boundaries.

## Prepared v4 query identity

The current source prepares `github.com/faustbrian/go-localized/v4` at the
repository root. This baseline is published as v4.0.0. Move every Localized import
to `/v4`, including canonical and retained adapters, and move query values,
operators and predicates together to `github.com/faustbrian/go-api-query/v4`.
International stays on `/v3/locale` v3.0.0 and Validation stays on `/v2` v2.0.0.
Select the published API Query v4.0.0 module without sibling replacements.
Localized v4 remains available as a published consumer release.

Both `localizedquery` and `adapters/query` expose the same API Query v4 nominal
types. Exact lookup still distinguishes missing from present-empty and never
applies language matching or fallback. Text construction, copying, lifecycle,
sentinels and encoded formats are unchanged by this namespace migration.
Historical v1–v3 API projections and public tags remain available. Tools adoption
and future-major release qualification remain separate boundaries.

## Official v3 module and dependency identities

Use `github.com/faustbrian/go-localized/v3` and append `/v3` before package
suffixes in every owning import. Migrate locale values to
`github.com/faustbrian/go-international/v3/locale`, query values, operators and
predicates to `github.com/faustbrian/go-api-query/v3`, selecting v3.0.0 for both
dependencies. Keep validators on `github.com/faustbrian/go-validation/v2`
v2.0.0. Types from different majors are not interchangeable; use the public
tags and releases to establish dependency publication before adoption.

Construction, exact lookup, matching, normalization and encoded formats remain
unchanged. Canonical and retained adapter paths remain in this module, with
their existing `Error`, `Rule` and `Form` aliases preserved within v3. The v1
and v2 API projections and tags are retained; existing consumers can stay on
their selected earlier-major dependencies. The maintained Tools compatibility
consumer selects Localized v2.0.0 with API Query and International v2; it is not
evidence of v3 adoption.

## Target-oriented adapter imports

New consumers should import `adapters/config`, `adapters/httpclient`,
`adapters/query`, `adapters/validation`, and `adapters/wire`. Existing
`localizedconfig`, `localizedhttpclient`, `localizedquery`,
`localizedvalidation`, and `localizedwire` imports remain supported. Their
named values and errors are aliased by the canonical packages, so migration is
behavior-compatible with no serialization, ownership, default, or runtime
change. Either adopt the canonical package identifier and update selectors, or
retain the old identifier with an explicit import alias:

```go
import localizedvalidation "github.com/faustbrian/go-localized/v4/adapters/validation"
```

## Locale-keyed Go maps

1. Inventory whether keys contain underscores, aliases, invalid tags, or
   canonical duplicates.
2. Decide duplicate and `und`/`mul`/private-use policy explicitly.
3. Construct with `TextFromMap` or `NewTextWithOptions` at ingress.
4. Replace direct indexing with `Get` and preserve the presence boolean.
5. Add matching or fallback only at call sites that previously promised it.
6. Persist `EncodeJSON` output and compare semantic round trips before deleting
   legacy code.

## Spatie Translatable JSON

Spatie objects map locale strings to strings and commonly contain
present-empty values. The fixture `postgres/testdata/spatie-translatable.json`
is accepted without semantic drift. PHP locale aliases or underscore keys MUST
be audited; strict v1 decode rejects underscores, while `PermissiveJSON` can be
used only as an explicit migration bridge before canonical re-encoding.

## Track

Treat Track display text as a domain value, not a translation catalog. Decode
the existing JSON column strictly, report canonical duplicates, write canonical
objects in a shadow or audited migration, and switch readers only after
round-trip comparison. The representative `track.json` fixture includes a
regional Swedish tag.

## Postal

Postal names are content; postal-code and locality lookup behavior remains in
Postal. The `postal.json` fixture proves Finnish and Swedish text. Do not infer a
fallback from country or postal code in this package.

## Location

Location pickup-point names may contain regional tags. The `location.json`
fixture preserves `fi-FI`. Coordinate, carrier, and nearest-location selection
remain outside this package.

## Normalized rows

Use `postgres.Rows` and `postgres.FromRows` to bridge `(entity_id, locale,
text)` schemas. They do not create tables or run migrations. Order rows by
entity and canonical locale for reproducible comparisons. Example extraction
SQL is in `postgres/testdata/normalized-rows.sql`.

## Rollout checks

- compare locale count, canonical key set, present-empty set, and hash;
- classify every rejected key before changing data;
- dual-read without silently preferring one representation;
- never materialize fallback results into stored translations;
- retain rollback data until semantic counts match.
