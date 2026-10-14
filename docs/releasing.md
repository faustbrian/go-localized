# Releasing

[Documentation index](README.md) · [Repository entry](../README.md)

1. Finish implementation and documentation locally with a clean worktree.
2. Run `go mod tidy -diff`, `make check-full`, and the selected shared release
   rehearsal. Ordinary `make check` and `make ci` are bounded local routes, not
   release qualification. PostgreSQL changes also require the explicit tagged
   integration test and PostgreSQL 14–18 matrix described in `CONTRIBUTING.md`.
   The CI workflow input `root_release_qualification` runs these full checks
   and rehearsal on the selected source with private hosted resources. It is
   separate from the rehearsal-only `release_dry_run` input.
3. Review mutation survivors, fuzz crashes, vulnerability findings, benchmark
   changes, dependency licenses, and locale-data changes.
4. Update `CHANGELOG.md`, compatibility provenance, API baseline, and evidence.
5. Create a signed semantic-version tag only after the user verifies final
   hosted CI. Git state is never an implementation blocker.
6. CI does not publish tags or releases. Build a deterministic source archive
   and checksums from the verified tag, then publish the changelog notes
   separately. Verify the remote tag, release assets and public module
   resolution before reporting publication complete.

No release may claim successful hosted checks from local workflow syntax alone.

Choose the next guide in the [documentation index](README.md)
or return to the [repository entry](../README.md) for installation and ownership.
