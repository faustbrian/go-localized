GO ?= go

.PHONY: docs postgres-matrix standards

docs:
	./verification/check-docs.sh

postgres-matrix:
	./verification/check-postgres-matrix.sh

standards:
	$(GO) test ./... -run 'Standards|Canonicalization|MatchingMatrix|LegacyCompatibility'
