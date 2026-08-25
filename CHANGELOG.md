# Changelog

All notable changes follow Keep a Changelog and Semantic Versioning.

## [Unreleased]

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

## [1.0.0] - 2026-08-25

### Documentation

- Replace obsolete standalone-repository links and workflow claims with
  monorepo-canonical targets and current release guidance.
- Document the package's initial stable `v1.0.0` scope and security policy.

- Link the package README to the repository-wide Golib documentation portal.

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
