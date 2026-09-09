package config_test

import (
	"reflect"
	"testing"

	localized "github.com/faustbrian/go-localized"
	config "github.com/faustbrian/go-localized/adapters/config"
	legacy "github.com/faustbrian/go-localized/localizedconfig"
)

func TestCanonicalConfigAdapterPreservesCompatibilityIdentity(t *testing.T) {
	t.Parallel()

	value, err := localized.TextFromMap(map[string]string{"en": "Hello"})
	if err != nil {
		t.Fatal(err)
	}
	got := config.NewText(value)
	if !got.Valid || !got.Localized.Equal(value) {
		t.Fatalf("NewText() = %+v", got)
	}
	if reflect.TypeFor[config.Text]() != reflect.TypeFor[legacy.Text]() {
		t.Fatal("compatibility Text identity changed")
	}
	if reflect.TypeOf(config.ErrInvalidValue) != reflect.TypeOf(legacy.ErrInvalidValue) {
		t.Fatal("compatibility Error identity changed")
	}
}
