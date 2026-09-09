package validation_test

import (
	"errors"
	"reflect"
	"testing"

	localized "github.com/faustbrian/go-localized"
	validation "github.com/faustbrian/go-localized/adapters/validation"
	legacy "github.com/faustbrian/go-localized/localizedvalidation"
	validationcore "github.com/faustbrian/go-validation"
)

func TestCanonicalValidationAdapterPreservesRulesAndIdentities(t *testing.T) {
	t.Parallel()

	value, _ := localized.TextFromMap(map[string]string{"en": " Hello \n"})
	rules := []validation.Rule{
		validation.RequireNonEmpty(),
		validation.RequireNonWhitespace(),
		validation.MaxBytes(16),
		validation.MaxRunes(16),
		validation.MaxLines(2),
		validation.NoControlCharacters(),
	}
	if err := validation.Validate(value, rules...); err != nil {
		t.Fatal(err)
	}
	validator := validation.Validator(rules...)
	ctx, err := validationcore.NewContext(validationcore.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	if report := validator.Validate(ctx, value); report.HasErrors() {
		t.Fatalf("Validator() report = %+v", report)
	}
	for _, form := range []validation.Form{validation.NFC, validation.NFD, validation.NFKC, validation.NFKD} {
		if _, err := validation.Normalize(value, form); err != nil {
			t.Fatalf("Normalize(%v) error = %v", form, err)
		}
	}
	if reflect.TypeFor[validation.Error]() != reflect.TypeFor[legacy.Error]() ||
		reflect.TypeFor[validation.Form]() != reflect.TypeFor[legacy.Form]() {
		t.Fatal("compatibility named type identity changed")
	}
	if !errors.Is(validation.ErrInvalidRule, legacy.ErrInvalidRule) {
		t.Fatal("compatibility error identity changed")
	}
}
