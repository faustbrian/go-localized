package localized_test

import (
	"testing"

	apiquery "github.com/faustbrian/go-api-query/v5"
	"github.com/faustbrian/go-international/v5/locale"
	localized "github.com/faustbrian/go-localized/v6"
	query "github.com/faustbrian/go-localized/v6/adapters/query"
	validation "github.com/faustbrian/go-localized/v6/adapters/validation"
	validationcore "github.com/faustbrian/go-validation/v3"
)

func TestPublicQueryValidationComposition(t *testing.T) {
	t.Parallel()
	english, err := locale.Parse("en")
	if err != nil {
		t.Fatal(err)
	}
	value, err := localized.NewText(localized.Entry{Locale: english, Text: "Hello"})
	if err != nil {
		t.Fatal(err)
	}
	for _, exact := range []func(localized.Text, locale.Tag) (apiquery.Value, bool){
		query.ExactValue,
	} {
		got, present := exact(value, english)
		if !present || got.String() != "Hello" {
			t.Fatalf("exact value = %v, %v", got, present)
		}
	}
	for _, predicate := range []func(string, apiquery.Operator, localized.Text, locale.Tag) (*apiquery.Predicate, error){
		query.ExactPredicate,
	} {
		got, err := predicate("title", apiquery.OpEqual, value, english)
		if err != nil || got == nil || got.Name != "title" || got.Operator != apiquery.OpEqual || len(got.Values) != 1 || got.Values[0].String() != "Hello" {
			t.Fatalf("exact predicate = %+v, %v", got, err)
		}
	}
	ctx, err := validationcore.NewContext(validationcore.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	for _, validator := range []validationcore.Validator[localized.Text]{
		validation.Validator(validation.RequireNonEmpty()),
	} {
		report := validator.Validate(ctx, value)
		if report.Len() != 0 || report.Err() != nil {
			t.Fatalf("normal validator report = %s", report.String())
		}
	}
}
