# Releasing

1. Finish implementation and documentation locally with a clean worktree.
2. Run `go mod tidy -diff`, `make check`, and the shared release gates through
   `make ci`.
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
