# Third-party notices

Direct runtime dependencies are pinned in `go.mod`:

- `github.com/faustbrian/go-international`, MIT, provides locale identity and
  pinned IANA registry provenance;
- `golang.org/x/text`, BSD-3-Clause, provides private matching data, Unicode
  normalization, and CLDR-derived behavior;
- `github.com/jackc/pgx/v5`, MIT, provides native PostgreSQL and JSONB codecs;
- `github.com/faustbrian/go-wire`, MIT, provides bounded YAML, TOML, JSON, and
  MessagePack adapters;
- `github.com/faustbrian/go-validation`, MIT, provides typed bounded reports;
- `github.com/faustbrian/go-api-query`, MIT, provides typed query values;
- `github.com/faustbrian/go-http-client`, MIT, provides immutable request
  specifications.

`github.com/faustbrian/go-config/v3`, MIT, is required only by tests of the
configuration hook contract.

Transitive dependencies and cryptographic checksums are recorded in `go.sum`.
Release review MUST re-run `go-licenses` or an equivalent license inventory
before distributing a binary that embeds this module.
