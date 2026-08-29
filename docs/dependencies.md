# Dependencies

Core depends on `international/locale` for bounded BCP 47 identity,
canonicalization, parent fallback, and registry provenance. It uses `x/text`
privately for language matching and Unicode normalization. PostgreSQL and
optional wire/config/query packages introduce pgx, wire, config, and
api-query. The HTTP client adapter adds http-client. Versions and
checksums are pinned in `go.mod` and `go.sum`.

Tool commands use explicit versions for govulncheck, NilAway, and Gremlins.
GitHub Actions use full commit SHAs with human-readable version comments.

Owned dependencies are consumed as released modules. `go.mod` contains no
checkout-relative replacements or workspace-only resolution. The
checksum-pinned `go-library-tools` contract validates the module graph and
prevents dependency or action drift. See [compatibility](compatibility.md).
