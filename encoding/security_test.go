package encoding_test

import (
	"errors"
	"testing"

	localized "github.com/faustbrian/go-localized/v4"
	localizedencoding "github.com/faustbrian/go-localized/v4/encoding"
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
	for _, dimension := range []string{"locales", "tag", "text", "total"} {
		t.Run(dimension, func(t *testing.T) {
			limits := localized.DefaultLimits()
			var limit *int
			switch dimension {
			case "locales":
				limit = &limits.MaxLocales
			case "tag":
				limit = &limits.MaxTagBytes
			case "text":
				limit = &limits.MaxTextBytes
			case "total":
				limit = &limits.MaxTotalBytes
			}
			*limit = -1
			_, err := localizedencoding.UnmarshalEntries([]byte(`true`), localizedencoding.DecodeOptions{Limits: limits})
			if !errors.Is(err, localized.ErrLimitExceeded) {
				t.Fatalf("negative %s limit before parsing = %v", dimension, err)
			}
			*limit = 0
			value, err := localizedencoding.UnmarshalEntries([]byte(`[]`), localizedencoding.DecodeOptions{Limits: limits})
			if err != nil || !value.IsEmpty() {
				t.Fatalf("zero %s limit for empty input = %v, %v", dimension, value.Entries(), err)
			}
		})
	}
}

func TestSecurityEntriesExactAndCumulativeTextBudgets(t *testing.T) {
	limits := localized.DefaultLimits()
	limits.MaxLocales = 4
	limits.MaxTagBytes = 2
	limits.MaxTextBytes = 2
	limits.MaxTotalBytes = 4
	value, err := localizedencoding.UnmarshalEntries([]byte(`[{"locale":"en","text":"ok"},{"locale":"fi","text":"ok"}]`), localizedencoding.DecodeOptions{Limits: limits})
	if err != nil || value.Len() != 2 {
		t.Fatalf("exact entry tag/text/total budgets = %v, %v", value.Entries(), err)
	}
	limits.MaxTotalBytes = 5
	_, err = localizedencoding.UnmarshalEntries([]byte(`[{"locale":"en","text":"ok"},{"locale":"fi","text":"ok"},{"locale":"sv","text":"ok"},{"unknown":0}]`), localizedencoding.DecodeOptions{Limits: limits})
	if !errors.Is(err, localized.ErrLimitExceeded) {
		t.Fatalf("cumulative entry bytes before malformed input = %v", err)
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
