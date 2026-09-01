# Changelog

All notable changes follow Keep a Changelog and Semantic Versioning.

## [Unreleased]

### Changed

- Use the released `go-library-tools` v1.2.0 CLI and immutable merged
  workflow at `1f9629e5f27418600460b55a50a5b2fc81697fab` while preserving
  package-owned verification evidence and source-specific checks.

### Documentation

- Replace historical repository links and completed execution artifacts with a
  standalone, human-oriented documentation structure.
- Govern BCP 47 identity, IANA registry use, language matching, fallback,
  `Accept-Language`, and JSON choices through the
  [specification decision register](docs/specification-decisions.md).
- Current decision records: `LOCALIZED-DEC-001 sha256:37ebb82077f351332005c6004546474bfa67d618cf5bbeb855e406dd559de7c1`;
  `LOCALIZED-DEC-002 sha256:b7a30860890d67cfea8ce5c0bd2f6e4fc90bea708bb5b5c3b49b60b4f227f515`;
  `LOCALIZED-DEC-003 sha256:ad97d7ae459d27f9103e951c7cb020671d2f2d55deeb100fe114ba93b505af41`;
  `LOCALIZED-DEC-004 sha256:58f9434796ee9b5a9cc70aa43761d771168de04641a0105fde7b3e8fa087627a`;
  `LOCALIZED-DEC-005 sha256:d24678855c787b97e7567d1076b5cb590afc4046397e9f7d438b5f7bdc6bc891`;
  `LOCALIZED-DEC-006 sha256:b7ed14e5ee8d3817bc0c22c6f811c8a5139e757c9e78a56f2b60b1d1270ad82c`;
  `LOCALIZED-DEC-007 sha256:34bd3a418f23297185f8aab87bfae1c424ecd22817449e79f4cdaa0a1d032419`;
  `LOCALIZED-DEC-008 sha256:9b486e60b10afa18beae1d9f047bf10c9e282cbbca88bfbd801bce1df5d45e7b`.

## [1.0.0] - 2026-08-25

### Changed

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Documentation

- Replace obsolete repository links and workflow claims with current release
  guidance.
- Document the package's initial stable `v1.0.0` scope and security policy.

- Link the package README to package-owned documentation.

### Fixed

- Reject non-numeric matching weights instead of silently treating them as
  absent preferences.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-localized` identity while preserving its documented API and behavior.
- Refresh local `v0.0.0` owned-module checksums after dependency manifests and
  release notes were normalized; runtime behavior and public APIs are
  unchanged.
- Upgrade `golang.org/x/text` to v0.41.0 so the isolated module graph remains
  aligned with the repository security baseline.
- Pin owned dependencies to published source revisions so clean consumers can
  resolve the module before the first tags.
- Snapshot localized validation rules in an explicit immutable adapter instead
  of relying on analyzer-ambiguous closure capture.
- Refresh owned-module checksums against the final consolidated archives.
- Normalized standalone module metadata against the canonical owned dependency
  graph, including complete checksums for clean consumer resolution.
- Refreshed the canonical HTTP client checksum after its boundary tests
  changed, preserving isolated module verification.
- Refreshed the canonical API query checksum after its API compatibility
  tooling was standardized.

### Added

- Immutable localized text, explicit matching and fallback, deterministic
  encoding, validation, persistence, HTTP, config, wire, and test adapters.
- Bounded hostile-input, race, fuzz, mutation, PostgreSQL, benchmark, and
  compatibility gates.
- Public locale identity and registry provenance through
  `international/locale`.
- Typed `validation` findings with canonical locale paths and content-free
  diagnostic codes.
- Exact, presence-aware `api-query` values and predicates without implicit
  matching or persistence policy.
- Canonical bounded `Accept-Language` integration for immutable http-client
  request specs and standard responses.
- Strict ordered-pair construction with explicit duplicate and limit options.
- Enforced allocation ceilings for construction, lookup, matching, fallback,
  merge, and canonical JSON operations.
- Property fuzzing for canonicalization, merge identities, deterministic order,
  equality, hashes, and canonical round trips.
- Reproducible dependency pins matching the exact locally verified sibling
  revisions.

### v1.0.0 scope

The following initial scope is included in `v1.0.0`.

#### Added

- Initial production contract for localized domain values.

[Unreleased]: https://github.com/faustbrian/go-localized/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/faustbrian/go-localized/releases/tag/v1.0.0
