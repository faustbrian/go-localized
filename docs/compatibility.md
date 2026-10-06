# Compatibility and provenance

## Toolchain and dependencies

Published Localized v4 remains separately available. Prepared Localized v5
selects International v4.0.0, HTTP Client v2.0.0 and Wire v3.0.0; its own
release and consumer proof remain pending. Query v4.0.0, Validation v2.0.0,
Config v2.0.0 and x/text v0.42.0 remain selected. The v1–v4 API projections
remain historical records, not proof of v5 compatibility.

| Component | Pinned baseline | Role |
|---|---|---|
| Go | 1.27.0 | minimum language and iterator contract |
| `international/v4/locale` | v4.0.0 | public BCP 47 identity and provenance |
| `api-query/v4` | v4.0.0 | public query values, operators and predicates |
| `validation/v2` | v2.0.0 | public validator identity |
| `golang.org/x/text` | v0.42.0 | private CLDR matching and Unicode normalization |
| pgx | v5.11.0 | JSONB and PostgreSQL integration |
| `wire/v3` | v3.0.0 | bounded format adapters |
| `http-client/v2` | v2.0.0 | immutable request specs and layers |
| `config/v2` | v2.0.0 | configuration hook conformance |
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
and result kinds preserve the earlier-major value commitments in v5. Nominal
locale, HTTP and wire identities adopt their supported successor majors;
query remains v4 and validator remains v2. Types cannot be mixed across module
majors. The original v1–v4 API baselines are preserved; v5 uses a separate
projection. Matcher choices may change
only with a documented locale-data dependency update and compatibility vectors.
