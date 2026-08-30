# Specification conformance matrix

The [specification decision register](../docs/specification-decisions.md) owns
the interpretations behind these bindings. `monitoring.json` pins each
normative or registry authority separately from its errata or update surface.

| Decision | Observable boundary | Peer boundary |
|---|---|---|
| LOCALIZED-DEC-001 | BCP 47 and IANA canonical locale identity | Maintained behavior is inherited from pinned `go-international`; no independent local differential harness |
| LOCALIZED-DEC-002 | Special, private-use, and registry-valid tag acceptance | Not assessed independently |
| LOCALIZED-DEC-003 | Duplicate identity after canonicalization | Not assessed independently |
| LOCALIZED-DEC-004 | Exact lookup versus RFC-style matching | Pinned `x/text` matcher behavior is executable dependency evidence, not an independent peer |
| LOCALIZED-DEC-005 | Explicit fallback graph and parent traversal | Not assessed independently |
| LOCALIZED-DEC-006 | `Accept-Language` parsing, quality ordering, and wildcard | Not assessed independently |
| LOCALIZED-DEC-007 | RFC 8259 object mapping and duplicate-member policy | Not assessed independently |
| LOCALIZED-DEC-008 | Legacy JSON underscore and null acceptance | Not assessed independently |

Empty differential lanes are deliberate: the package does not claim an
independent maintained-peer comparison where its implementation directly
delegates to the pinned locale or matcher dependency.
