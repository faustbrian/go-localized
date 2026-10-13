# Source documentation audit

[Documentation index](README.md) · [Repository entry](../README.md)

## Pre-edit scope and plan

Base `a231c5c60db231626b15094e58f5dbf5ce75c503`; comments/docs-only Tier A.
One root `github.com/faustbrian/go-localized/v5` module, 16 package directories,
75 tracked Go files. Original unrelated chore checkout is preserved; this is
an isolated detached main-source worktree. Runtime, API/dependency projections,
fixtures and frozen specification authorities stay unchanged.

Reproduce with `git ls-files '*.go'`, catalogs and native shared
`source-docs inventory --json` / `source-docs check --json`. Initial 91 useful
export/member/interface-comment gaps; counts locate omissions, not conformance.

## Complete tracked inventory

- `.`: `benchmark_test.go`, `boundary_test.go`, `error.go`, `hardening_test.go`, `json.go`, `json_fuzz_test.go`, `json_test.go`, `locale_helper_test.go`, `major_compatibility_test.go`, `properties_fuzz_test.go`, `security_bounds_test.go`, `successor_composition_test.go`, `text.go`, `text_advanced_test.go`, `text_operations_test.go`, `text_test.go`.
- `adapters/config`: `config.go`, `config_test.go`, `doc.go`.
- `adapters/httpclient`: `doc.go`, `httpclient.go`, `httpclient_test.go`.
- `adapters/query`: `doc.go`, `query.go`, `query_test.go`.
- `adapters/validation`: `doc.go`, `validation.go`, `validation_test.go`.
- `adapters/wire`: `doc.go`, `wire.go`, `wire_test.go`.
- `encoding`: `entries.go`, `entries_fuzz_test.go`, `entries_test.go`, `locale_helper_test.go`, `security_test.go`.
- `http`: `accept_language.go`, `accept_language_fuzz_test.go`, `accept_language_test.go`, `locale_helper_test.go`.
- `localizedconfig`: `config.go`, `config_test.go`, `locale_helper_test.go`, `security_test.go`.
- `localizedhttpclient`: `httpclient.go`, `httpclient_test.go`.
- `localizedquery`: `query.go`, `query_test.go`.
- `localizedtest`: `locale_helper_test.go`, `localizedtest.go`, `localizedtest_test.go`.
- `localizedvalidation`: `locale_helper_test.go`, `validation.go`, `validation_test.go`.
- `localizedwire`: `locale_helper_test.go`, `wire.go`, `wire_fuzz_test.go`, `wire_test.go`.
- `match`: `error.go`, `fallback_fuzz_test.go`, `locale_helper_test.go`, `match.go`, `match_test.go`, `observe.go`, `observe_test.go`, `plan.go`, `plan_test.go`, `security_test.go`.
- `postgres`: `integration_test.go`, `legacy_test.go`, `locale_helper_test.go`, `postgres.go`, `postgres_fuzz_test.go`, `postgres_test.go`, `security_test.go`.

## Source-contract coverage and private rationale

- Root canonical locale indexes, sorted immutable entries, presence/empty/null,
  cumulative pre-dedup text budgets, strict string boundaries, construction versus
  operation policies, framed hashing and builder accumulation. Locale identity,
  registry aliases and parent policy belong to the pinned International v4 dependency.
  Core does not translate, pluralize, interpolate, format, infer languages or own CLDR.
- JSON token-level canonical duplicate detection, entry-array streaming admission,
  receiver transactionality and encoding ownership. UTF-8 source admission is not
  a strict surrogate-escape or complete nested object-shape guarantee.
- Matching stable weight ordering, per-preference exact-then-x/text selection and
  HTTP wildcard special handling. Graph edge/source caps, longest-path validation,
  cycle states and shared-subgraph visited set; ParentRange starts at the strict
  parent, not the candidate itself. No fallback is implicitly stored as a new entry.
- Synchronous transform/resolver/rule callbacks, panic-isolated observers, copied
  rule/candidate slices and retained borrowed observer/rule implementations. No
  callback time/cancellation guarantee, ambient I/O, worker or shutdown is owned.
- Five canonical target-oriented adapter facades and five retained implementations:
  aliases preserve named type/error identity, all wrapper paths share real behavior.
  Wire/config limits and error details are dependency-owned; query helpers select
  exact presence only and do not validate a schema or execute queries. HTTP client
  preferences are trusted outbound inputs with dynamically derived admission caps.
- PostgreSQL nullable wrappers, transactional Scan, generic versus localized pgx
  JSONB branches, row copying and already-fetched byte ownership. No connection,
  transaction, SQL execution, migration, live database or pool lifecycle is owned.
- Test support fixtures/assertions (synthetic data may be printed on failures);
  all tests reviewed by contract families: constructor/presence/transform/merge,
  strict/legacy wire and security boundaries, matching/fallback/observer, canonical
  adapter identity/composition, major dependency composition and PostgreSQL fixtures.
  Eight fuzz targets, eight benchmarks, allocation checks and test helpers inventoried.
  No Examples currently present. One integration build constraint gates live PostgreSQL.
  No generated Go or go:generate directives are present; frozen fixture/specification
  and all five major API projections remain byte-exact, not regenerated or requalified.

## Missing comments versus existing limitations/questions

Repair 91 useful member/error/interface gaps and package/API ownership and rationale.
Limits normalize only the wholly zero struct; individual zeros in a partially
specified struct remain zero allowances. MaxTotalBytes counts input text before
FirstWins/LastWins reduction, not tags or total Go heap. Raw pair/JSON keys retain
an independent 255-byte parser ceiling even when callers raise value limits.
Builders and outbound headers allocate before final validation; merge callbacks
execute before output validation, and external callback effects are not rolled back.

Existing transform limitation: Set/Map/Normalize revalidate with default construction
limits rather than remembering the originating custom policy. Remove/Filter also
reconstruct with defaults, but discard the error and can silently return empty Text
for an otherwise valid custom-limit value. This is not an intentional policy or a
behavioral repair. A tiny disposable observation is appropriate; no broad campaign.

Entry-array structs follow encoding/json matching: duplicate object fields overwrite,
case aliases are accepted, and omitted/null text becomes present-empty. Raw malformed
UTF-8 is rejected, but unpaired JSON surrogate escapes can repair to U+FFFD. These
existing wire limitations are not normative interpretations or conformance success.

Unknown rationale/decisions: whether transform methods should retain input limits,
whether Remove/Filter should expose construction failures, whether entry-array shape
should reject null/missing/case-alias/duplicate fields and lone surrogate escapes.
No runtime, dependency or fixture change is authorized in this source pass.

## Selected checks

One explicitly granted serialized Go 1.27.2 bounded session: formatting, exact all 75
comment-free AST/signature/raw-token/literal/implicit-semicolon and directive identity,
all 16 package GoDocs and native docs script. `.golib.yaml` selects apidiff api/v5.txt;
no full-prose snapshot generation applies. No native Examples exist; script's Example
selector is structural compilation/no-examples evidence, not behavioral suite proof.

Use task-owned GOCACHE/GOMODCACHE/GOTMPDIR with cleanup on all exits. Tiny disposable
transform and decoder probes substantiate existing limits only. No live PostgreSQL,
full standards matrix, race/fuzz/mutation/benchmark or public-consumer campaign.
Also choose native source check, affected links/markers, offline specification check,
pinned added/staged-content secrets, complete ordinary Tier A review and explicit-path
Intent/ACK safe-wrapper commit. Parent owns remotes/CI/publication/tooling/ledger.

## Results

- Go 1.27.2: changed-file gofmt and all 75-file formatting passed. Disposable
  parser/scanner proof passed exact executable AST/signatures, raw tokens,
  literals, implicit semicolons and directives; 316 comments parsed/rendered.
- Native `go doc -all` rendered all 16 catalog packages; package and contract
  output inspected. Native `make -f verification/package.mk docs` passed required
  files/all Markdown links and its Example selector across 16 packages; all
  reported no tests to run. No runtime suite or live database qualification.
- Read-only disposable imported-source probes confirmed valid custom-limit
  Remove/Filter silently empty; entry-array omitted/null text and case aliases
  present-empty, repeated text last-wins, and lone-surrogate repair to U+FFFD
  in both entry-array and core JSON. These are observed existing limitations,
  not behavioral repairs or specification interpretation changes.
- Current shared source checker: zero findings; offline specification check:
  eight decisions valid. Affected links: 15 resolve; no unfinished markers.
  All 19 protected API/dependency/catalog/specification/verification/history
  paths remain byte-identical. Apidiff is configured; no prose snapshot applies
  and native API analyzer/full conformance campaigns were not executed.
- Exact owned GOCACHE/GOMODCACHE/GOTMPDIR and imported-source probe were removed
  after the bounded session. Disposable parser/render/inventory resources are
  removed after this evidence capture; parent-owned CLI/scanner are retained.

Final added/staged-content secret and complete ordinary Tier A review results
are reported in the worker handoff, separately from native source proof.

## Objective Enforcement

The checksum-pinned public Tools v2.4.0 verifier now supplies
`make source-docs-check` as a prerequisite of the existing local and CI Make
routes. Ordinary shared CI selects strict source-documentation checking and
requires its result. Native prerequisites, operation budgets, security and
release controls are retained. Hosted exact-head evidence is recorded by the
coordination ledger; configuration alone does not establish green delivery.

Choose the next guide in the [documentation index](README.md)
or return to the [repository entry](../README.md) for installation and ownership.
