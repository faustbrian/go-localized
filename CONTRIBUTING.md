# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected package's documentation.
2. Run ordinary Go tests for the affected package before editing.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update
the affected [stable decision entries](docs/specification-decisions.md), and
complete the Specification Decisions
section of the pull request template. An unresolved interpretation or stale
source pin is release-blocking; peer behavior cannot silently select policy.

Selected mutation gates must finish with zero surviving viable mutants.

Do not add package-local workflows, permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use Go 1.27.2 for development and CI; the public minimum remains
Go 1.27.0. An ordinary checkout is sufficient for package tests and a focused
active fuzz smoke; no sibling repository or goal ledger is required:

```bash
go test ./...
go test . -run '^$' -fuzz '^FuzzDecodeJSON$' -fuzztime=2s
```

Select the native target and package affected by your change. Ordinary Go tests
already exercise deterministic fuzz seeds. For shared checks, install the
public `golib` v2.1.0 release declared in `.golib.yaml` and verify its checksums;
CI uses that same release, not a sibling source build:

```bash
make check
make ci
```

Both commands use bounded local checks and retain the manifest's PostgreSQL
service and integration tag, so the verifier routes require Docker. Untagged
direct Go tests are the database-free route; a supplied `POSTGRES_URL` for direct
tagged tests is distinct from verifier-managed services. Optional focused
mutation analysis retains those service prerequisites and can select one package:

```bash
golib mutation --module . --package ./encoding
```

Use `make check-full` for explicit full assurance with the configured native fuzz
budgets, standards, API compatibility and other selected full gates. For
PostgreSQL JSONB changes, supply a disposable database and explicitly run
`go test -count=1 -tags=integration ./postgres` with `POSTGRES_URL` set; the
manual `make -f verification/package.mk postgres-matrix` checks PostgreSQL
14 through 18. Run affected adapters and direct consumer tests for public-contract
changes. Scheduled and release qualification retain their full selected gates;
follow [the release guide](docs/releasing.md), and report unavailable or failing
checks without claiming release readiness. CI/configuration-only changes do not
require a library release.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
