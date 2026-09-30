package encoding_test

import (
	"errors"
	"testing"

	localized "github.com/faustbrian/go-localized"
	localizedencoding "github.com/faustbrian/go-localized/encoding"
)

func TestSecurityEntriesRejectCountBeforeFurtherDecoding(t *testing.T) {
	limits := localized.DefaultLimits()
	limits.MaxLocales = 1
	t.Run("count", func(t *testing.T) {
		_, err := localizedencoding.UnmarshalEntries([]byte(`[{"locale":"en","text":"ok"},{"unknown":0}]`), localizedencoding.DecodeOptions{Limits: limits})
		if !errors.Is(err, localized.ErrLimitExceeded) {
			t.Fatalf("entry count error = %v", err)
		}
	})
	t.Run("bytes", func(t *testing.T) {
		_, err := localizedencoding.UnmarshalEntries([]byte{0xff, 0xff}, localizedencoding.DecodeOptions{MaxInputBytes: 1})
		if !errors.Is(err, localized.ErrLimitExceeded) {
			t.Fatalf("input bytes error = %v", err)
		}
	})
}
