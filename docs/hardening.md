# Hardening report

## Scope

The hardening goal covers locale identity, exact lookup, matching, fallback,
merge, ownership, deterministic encoding, persistence, legacy fixtures,
Unicode, concurrent use, hooks, resource bounds, and privacy.

## Current evidence categories

- table-driven canonicalization, presence, merge, matching, fallback, HTTP,
  persistence, and legacy tests;
- fuzz targets for JSON, entry arrays, HTTP weights, fallback tags, SQL, pgx,
  and JSON/YAML/TOML/MessagePack wire decoders;
- race tests for all packages plus concurrent lookup, iteration, matching,
  fallback, merge, encoding, and observation;
- conditional-decision mutation gate configured for 100% mutant coverage and
  100% efficacy across runtime packages;
- statement coverage gate configured for exactly 100.0% across runtime
  packages, with fatal-only `localizedtest` support paths excluded;
- PostgreSQL 14–18 workflow and disposable local container matrix;
- production-source checks for unsafe, cgo, linkname, and mutable globals;
- content-free typed error constants and bounded observer events.

Completion requires the commands in [evidence](evidence.md) to pass from the
final tree. This document records the contract and must not be read as a claim
that future or hosted commands have already run.
