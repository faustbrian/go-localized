package query

import (
	apiquery "github.com/faustbrian/go-api-query/v4"
	"github.com/faustbrian/go-international/v4/locale"
	localized "github.com/faustbrian/go-localized/v5"
	legacy "github.com/faustbrian/go-localized/v5/localizedquery"
)

// ExactValue returns an api-query string value only when tag is present.
// It performs no language matching or application fallback.
func ExactValue(value localized.Text, tag locale.Tag) (apiquery.Value, bool) {
	return legacy.ExactValue(value, tag)
}

// ExactPredicate creates one api-query predicate from an exact localized
// value. Missing returns localized.ErrMissingLocale.
func ExactPredicate(
	name string,
	operator apiquery.Operator,
	value localized.Text,
	tag locale.Tag,
) (*apiquery.Predicate, error) {
	return legacy.ExactPredicate(name, operator, value, tag)
}
