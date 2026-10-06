# Compatibility Policy

Each independently releasable Go module follows semantic versioning. The root
module uses root Git tags such as `v4.0.0`; independently released nested modules
use `<module-directory>/v<version>` tags.

Implementation stays on `main`. Necessary breaking releases use new major Git
tags from main and Go's required major module/import suffixes, such as `/v4`.
Those suffixes do not imply version-specific source directories or branches.

Before `v1`, minor releases MAY contain reviewed breaking changes, but every
break MUST be documented with migration guidance. Patch releases MUST remain
backward compatible. At and after `v1`, incompatible exported API or documented
behavior changes require a new major version.

Compatibility includes exported Go APIs, error classification, serialization,
protocol behavior, persistence schemas, environment variables, command output,
resource ownership, ordering, retry/idempotency semantics, and documented
defaults. A compile-compatible change can still be behaviorally breaking.

Specification-backed modules MUST NOT diverge from their declared standards.
Ambiguities require documented decisions and stable tests. Deprecated APIs
follow [`DEPRECATION.md`](DEPRECATION.md).

The [specification decision register](docs/specification-decisions.md) governs
BCP 47 identity, matching, fallback, HTTP preference, and JSON wire choices.
Any observable decision change requires compatibility and changelog review.
