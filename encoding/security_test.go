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

func TestSecurityEntriesValidateLimitsBeforeMalformedInput(t *testing.T) {
	limits := localized.DefaultLimits()
	limits.MaxLocales = -1
	_, err := localizedencoding.UnmarshalEntries([]byte(`[{"unknown":0}]`), localizedencoding.DecodeOptions{Limits: limits})
	if !errors.Is(err, localized.ErrLimitExceeded) {
		t.Fatalf("negative entry limit error = %v", err)
	}
}

func TestSecurityEntriesRejectBudgetsBeforeLaterMalformedInput(t *testing.T) {
	for _, dimension := range []string{"tag", "text", "total"} {
		t.Run(dimension, func(t *testing.T) {
			limits := localized.DefaultLimits()
			switch dimension {
			case "tag":
				limits.MaxTagBytes = 1
			case "text":
				limits.MaxTextBytes = 1
			case "total":
				limits.MaxTotalBytes = 3
			}
			input := `[{"locale":"en","text":"ok"},{"locale":"fi","text":"ok"},{"unknown":0}]`
			_, err := localizedencoding.UnmarshalEntries([]byte(input), localizedencoding.DecodeOptions{Limits: limits})
			if !errors.Is(err, localized.ErrLimitExceeded) {
				t.Fatalf("entry %s budget error = %v", dimension, err)
			}
		})
	}
}

func TestSecurityEntriesRejectNonArrayAndMalformedClose(t *testing.T) {
	for _, input := range []string{`true`, `{}`, `[{"locale":"en","text":"ok"}`, `[{"locale":"en","text":"ok"}}`} {
		t.Run(input, func(t *testing.T) {
			value, err := localizedencoding.UnmarshalEntries([]byte(input), localizedencoding.DecodeOptions{})
			if !errors.Is(err, localized.ErrInvalidEncoding) || !value.IsEmpty() {
				t.Fatalf("malformed entry input returned %v, %v", value.Entries(), err)
			}
		})
	}
}
