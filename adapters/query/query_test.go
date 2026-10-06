package query_test

import (
	"testing"

	apiquery "github.com/faustbrian/go-api-query/v4"
	"github.com/faustbrian/go-international/v4/locale"
	localized "github.com/faustbrian/go-localized/v5"
	query "github.com/faustbrian/go-localized/v5/adapters/query"
)

func TestCanonicalQueryAdapterPreservesExactLookup(t *testing.T) {
	t.Parallel()

	value, _ := localized.TextFromMap(map[string]string{"en": "Hello"})
	english, _ := locale.Parse("en")
	got, present := query.ExactValue(value, english)
	if !present || got.String() != "Hello" {
		t.Fatalf("ExactValue() = %v, %v", got, present)
	}
	predicate, err := query.ExactPredicate("title", apiquery.OpEqual, value, english)
	if err != nil || predicate.Name != "title" || predicate.Values[0].String() != "Hello" {
		t.Fatalf("ExactPredicate() = %+v, %v", predicate, err)
	}
}
