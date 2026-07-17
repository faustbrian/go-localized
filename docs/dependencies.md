# Dependencies

Core depends on `go-international/locale` for bounded BCP 47 identity,
canonicalization, parent fallback, and registry provenance. It uses `x/text`
privately for language matching and Unicode normalization. PostgreSQL and
optional wire/config/query packages introduce pgx, go-wire, go-config, and
go-api-query. The HTTP client adapter adds go-http-client. Versions and
checksums are pinned in `go.mod` and `go.sum`.

Tool commands use explicit versions for govulncheck, NilAway, and Gremlins.
GitHub Actions use full commit SHAs with human-readable version comments.

The pinned `go-international`, `go-validation`, and `go-api-query` commits may
be resolved locally through an ignored Go workspace while awaiting publication;
`go.mod` contains no checkout-relative replacements. Hosted verification runs
only after those pins are published. `make dependency-revisions` rejects pin
drift and tests archived clean commits rather than sibling working trees. See
[compatibility](compatibility.md).
