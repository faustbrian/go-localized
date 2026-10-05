# Security model

Model version: 2. This document describes the v3.0.0 source published from
`95d141cd9d0316bed21c96b3752ce219770e550d`. The nominal v2 and v3 migrations
change module and dependency identities, not the localized algorithms or
limits described here. See [migration.md](migration.md) for migration details
and [../SECURITY.md](../SECURITY.md) for release-line and reporting policy.

## Threats and controls

| Threat | Control |
|---|---|
| Invalid UTF-8 or repair drift | reject before decoding or construction |
| Canonical duplicate smuggling | token-level/object or entry-level detection |
| Locale enumeration | no global registry API; observers omit tags |
| Private-use leakage | content-free errors/events; opt-in rejection policy |
| Fallback cycles/amplification | longest-path depth validation and one visit per configured chain |
| Oversized text/parser input | `Limits`, `MaxInputBytes`, HTTP and wire bounds |
| Mutable aliasing | copy all caller containers and database byte input |
| Content exposure | package errors never include localized strings |
| Observer failure | recover observer panics after resolution |
| Supply chain drift | pinned modules, checksums, Dependabot, vulnerability gate |

## Default budgets

| Resource | Default |
|---|---:|
| Locales per value | 128 |
| Canonical tag bytes | 255 |
| Bytes per text | 1 MiB |
| Total text bytes | 8 MiB |
| Canonical JSON input | 9 MiB |
| HTTP header bytes | 8 KiB |
| HTTP/match preferences | 64 |
| Fallback graph | caller-required depth and candidate limits |

Explicit limits are part of the operation, not mutable globals. Negative limits
are rejected. Zero selects documented defaults where an options type says so.

## Unicode

Text is valid UTF-8 but is not normalized by default. NFC, NFD, NFKC, and NFKD
are explicit persistent transforms. The package does not claim grapheme limits,
visual equivalence, confusable detection, HTML safety, language identity, or
semantic equivalence. Output contexts remain responsible for escaping.

Production-source gates reject unsafe, cgo, `go:linkname`, and package-level
variables. Race tests cover shared values, plans, codecs, and observers.

## Trust boundaries and ownership

The root constructors accept caller-owned locale tags, strings, maps and pair
slices. JSON and entry-array decoders accept hostile bytes. HTTP accepts an
untrusted `Accept-Language` header; matching and fallback consume explicit
caller policy. PostgreSQL wrappers accept driver values and normalized rows;
they do not open a database, issue SQL, or own a transaction. Config and wire
adapters inherit their upstream parser limits, then apply localized limits.
Target-oriented `adapters/*` entry points delegate to the same implementations.

There is no implicit network, filesystem, process, environment, goroutine,
registry, retry, authentication, or authorization boundary in this module.
Core decoding and resolution are synchronous CPU work under explicit limits.
Network and database cancellation belongs to the caller's HTTP client or
database context; passing an already-fetched value here does not add blocking
I/O. Query adapters create values and predicates, not SQL. Output escaping
and authorization remain caller owned.

Bounded byte decoders and collection constructors reject oversized input
before further parsing or collection allocation. JSON counts and text budgets
are checked as entries arrive. PostgreSQL `Scan(string)` first copies the
already-fetched driver string to bytes; callers must cap row bytes before this
adapter boundary to bound that transient copy.
Entry-array decoding does not first materialize an unbounded slice. A fallback
plan bounds source-chain count independently from candidate-edge count, checks
the longest path regardless of map iteration order, and visits shared chains
once per resolution. Limits supplied by callers are explicit trust decisions;
increasing them increases permissible memory and CPU work.

Localized content and private-use identifiers are potentially sensitive.
Package-generated errors are fixed identities and observers include only
operation, outcome kind and count. Validation adapter findings include locale
paths but not text; custom rule errors and upstream codec errors are outside
the core redaction guarantee. Applications must not log content, private-use
tags or arbitrary dependency error strings without their own redaction policy.

## Residual-risk register

| Risk | Owner | Rationale and mitigation | Review condition |
|---|---|---|---|
| Caller raises resource limits | Application owner | Explicit options intentionally permit larger trusted workloads; cap limits before accepting hostile input. | Limits or ingress trust change. |
| Builder accumulation before `Build` | Application owner | `Add` collects caller-owned construction policy and reports validation on `Build`; enforce the configured count and text budget while feeding a builder, or use bounded constructors for ingress. | Builder begins receiving attacker-selected input. |
| Outbound HTTP preference construction | Application owner | `WithPreferences` creates a header from caller-selected preferences and derives its budget from that list; bound list and tag sizes before using external preferences. Inbound selection has independent parser budgets. | Preference source becomes untrusted. |
| PostgreSQL driver string-copy admission | Application DB owner | `Scan(string)` copies an already-fetched driver string before bounded JSON decoding; cap row bytes before calling `Scan`. | Database field or driver ingress constraints change. |
| Custom observer or validation rule blocks or leaks | Application owner | Synchronous caller code cannot be forcibly cancelled safely; keep callbacks bounded and content-free. Observer panics are contained, not logged. | New callback or rule implementation. |
| Wire/config parser allocation or error disclosure | Wire/config maintainer and application owner | Dependency parsers own byte, depth and node budgets; localized checks the resulting collection. Do not expose raw upstream errors to clients. | Dependency or parser-policy change. |
| Validation report locale paths expose private-use tags | Application owner | Paths identify the failing field by contract; reject private-use tags or redact report paths when they identify tenants. | Reports cross a tenant or logging boundary. |
| Visually similar Unicode characters and markup injection | Application owner | Valid UTF-8 is not semantic or rendering safety; normalize and escape explicitly for the destination. | New rendering or identity use. |
| Maintainer or build dependency compromise | Repository maintainer | Review pinned dependencies and immutable CI actions, require hosted checks, and use private coordinated disclosure. | Dependency update, advisory or release. |

## Verification and release boundary

Deterministic hostile-input regressions live in `security_bounds_test.go`,
`encoding/security_test.go`, `match/security_test.go` and
`postgres/security_test.go`. Existing codec, HTTP, persistence, adapter and
canonicalization suites cover malformed inputs, transactional receiver state,
privacy-safe errors and concurrent reads. The shared CI workflow owns static,
dependency, secret and workflow gates; a source review or passing unit tests
alone do not establish those gates.

The historical v1 resource-bound and nil-destination corrections preserve
exported APIs, canonical values, matching order and successful encodings.
Oversized or invalid inputs can now fail earlier, and graphs exceeding the
existing depth policy are rejected consistently. Those fixes were published
in v1.1.2 as a patch;
the later module/dependency identity migrations require major releases.
Release acceptance still requires exact-main CI and a
clean public consumer; this threat model is not a claim of publication.
