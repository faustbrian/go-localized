// Package localizedquery adapts exact localized values to go-api-query.
package localizedquery

import (
	apiquery "github.com/faustbrian/go-api-query"
	"github.com/faustbrian/go-international/locale"
	localized "github.com/faustbrian/go-localized"
)

// ExactValue returns a go-api-query string value only when tag is present.
// It performs no language matching or application fallback, and preserves a
// present-empty string as a present query value.
func ExactValue(value localized.Text, tag locale.Tag) (apiquery.Value, bool) {
	text, present := value.Get(tag)
	if !present {
		return apiquery.Value{}, false
	}
	return apiquery.StringValue(text), true
}

// ExactPredicate creates one go-api-query predicate from an exact localized
// value. Missing returns localized.ErrMissingLocale; present-empty remains an
// explicit empty string for the schema to accept or reject.
func ExactPredicate(
	name string,
	operator apiquery.Operator,
	value localized.Text,
	tag locale.Tag,
) (*apiquery.Predicate, error) {
	queryValue, present := ExactValue(value, tag)
	if !present {
		return nil, localized.ErrMissingLocale
	}
	return &apiquery.Predicate{
		Name: name, Operator: operator, Values: []apiquery.Value{queryValue},
	}, nil
}
