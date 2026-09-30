package localizedconfig_test

import (
	"errors"
	"fmt"
	"testing"

	localized "github.com/faustbrian/go-localized"
	"github.com/faustbrian/go-localized/localizedconfig"
)

func TestSecurityConfigRejectsCountBeforeCopy(t *testing.T) {
	values := make(map[string]any)
	for i := range 129 {
		values[fmt.Sprintf("invalid_%d", i)] = 0
	}
	var target localizedconfig.Text
	if err := target.UnmarshalConfigValue(values); !errors.Is(err, localized.ErrLimitExceeded) {
		t.Fatalf("config count error = %v", err)
	}
}
