# Compatibility and provenance

[Documentation index](README.md) · [Repository entry](../README.md)

## Toolchain and dependencies

Localized v6 selects API Query v5, International v5, HTTP Client v3,
Validation v3 and Wire v3, with Config v3 in tests and x/text v0.42.0; v5
remains separately available. Clean-consumer qualification and maintained
Tools adoption require separate evidence; main's module identity does not
establish them. The v1–v5 API projections remain historical records, not proof
of v6 compatibility; v6 uses `api/v6.txt`.

| Component | Pinned baseline | Role |
|---|---|---|
| Go | 1.27.0 | minimum language and iterator contract |
| `international/v5/locale` | v5.0.0 | public BCP 47 identity and provenance |
| `api-query/v5` | v5.0.0 | public query values, operators and predicates |
| `validation/v3` | v3.0.1 | public validator identity |
| `golang.org/x/text` | v0.42.0 | private CLDR matching and Unicode normalization |
| pgx | v5.11.0 | JSONB and PostgreSQL integration |
| `wire/v3` | v3.0.0 | bounded format adapters |
| `http-client/v3` | v3.0.0 | immutable request specs and layers |
| `config/v3` | v3.0.0 | configuration hook conformance tests |
| PostgreSQL | 14–18 | JSONB integration matrix |

`locale.DatasetProvenance` records the IANA Language Subtag Registry retrieved
2026-07-16, upstream registry version 2026-06-14, x/text v0.42.0, and SHA-256
`be1fad86a99e3a932d07b80c9b3c271ec2381a5909ce22420144e5077ab0a43a`.
Releases MUST state changes to either locale dependency because preferred-value
canonicalization or matching can change.

## Locale classes

Language, script, region, variant, extension, private-use, grandfathered, and
deprecated tags follow the pinned locale dependency. Default construction
accepts valid `und`, `mul`, and private-use tags; `LocalePolicy` can reject each
class. Strict string boundaries reject underscores and whitespace.

Unknown-but-well-formed tags follow the pinned locale registry contract. The
package does not expose registry enumeration or mutable registry state.

## Stability

Canonical JSON, exact presence, missing/present-empty distinction, merge policy,
and result kinds preserve the earlier-major value commitments in v6. Nominal
locale, query, HTTP, validation and wire identities adopt their current
majors. Types cannot be mixed across module majors. The original v1–v5 API
baselines are preserved; v6 uses a separate projection. Matcher choices may change
only with a documented locale-data dependency update and compatibility vectors.

Choose the next guide in the [documentation index](README.md)
or return to the [repository entry](../README.md) for installation and ownership.
