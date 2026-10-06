package localizedconfig_test

import (
	"errors"
	"fmt"
	"testing"

	localized "github.com/faustbrian/go-localized/v5"
	"github.com/faustbrian/go-localized/v5/localizedconfig"
)

func TestSecurityConfigRejectsCountBeforeCopy(t *testing.T) {
	values := make(map[string]any)
	for i := range 129 {
		values[fmt.Sprintf("invalid_%d", i)] = 0
	}
	t.Run("any values", func(t *testing.T) {
		var target localizedconfig.Text
		if err := target.UnmarshalConfigValue(values); !errors.Is(err, localized.ErrLimitExceeded) {
			t.Fatalf("config count error = %v", err)
		}
	})
	stringValues := make(map[string]string)
	for key := range values {
		stringValues[key] = "text"
	}
	t.Run("string values", func(t *testing.T) {
		var target localizedconfig.Text
		if err := target.UnmarshalConfigValue(stringValues); !errors.Is(err, localized.ErrLimitExceeded) {
			t.Fatalf("string config count error = %v", err)
		}
	})
}

func TestSecurityConfigAcceptsExactLocaleCount(t *testing.T) {
	stringValues := make(map[string]string)
	values := make(map[string]any)
	for i := range localized.DefaultLimits().MaxLocales {
		key := fmt.Sprintf("x-entry-%d", i)
		stringValues[key] = "text"
		values[key] = "text"
	}
	for _, input := range []any{stringValues, values} {
		var target localizedconfig.Text
		if err := target.UnmarshalConfigValue(input); err != nil || !target.Valid || target.Localized.Len() != len(stringValues) {
			t.Fatalf("exact %T count = %d, present=%v, error=%v", input, target.Localized.Len(), target.Valid, err)
		}
	}
}
