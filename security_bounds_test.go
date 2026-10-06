package localized_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	localized "github.com/faustbrian/go-localized/v5"
)

func TestSecurityJSONLimitsBeforeFurtherInput(t *testing.T) {
	limits := localized.DefaultLimits()
	limits.MaxLocales = 1
	for _, test := range []struct {
		input   string
		maxText int
	}{{`{"en":"ok","not_a_locale":0}`, limits.MaxTextBytes}, {`{"en":"toolong","fi":0}`, 2}} {
		t.Run(test.input, func(t *testing.T) {
			options := localized.DecodeOptions{Limits: limits}
			options.Limits.MaxTextBytes = test.maxText
			if _, err := localized.DecodeJSON([]byte(test.input), options); !errors.Is(err, localized.ErrLimitExceeded) {
				t.Fatalf("bounded input error = %v, want limit", err)
			}
		})
	}
	t.Run("bytes", func(t *testing.T) {
		if _, err := localized.DecodeJSON([]byte{0xff, 0xff}, localized.DecodeOptions{MaxInputBytes: 1}); !errors.Is(err, localized.ErrLimitExceeded) {
			t.Fatalf("oversized invalid UTF-8 error = %v", err)
		}
	})
}

func TestSecurityNegativeLimitsBeforeLocaleParsing(t *testing.T) {
	limits := localized.DefaultLimits()
	limits.MaxLocales = -1
	t.Run("JSON", func(t *testing.T) {
		_, err := localized.DecodeJSON([]byte(`{"not_a_locale":0}`), localized.DecodeOptions{Limits: limits})
		if !errors.Is(err, localized.ErrLimitExceeded) {
			t.Fatalf("negative JSON limit error = %v", err)
		}
	})
	t.Run("pairs", func(t *testing.T) {
		_, err := localized.TextFromPairsWithOptions(localized.ConstructionOptions{Limits: limits}, localized.Pair{Locale: "not_a_locale"})
		if !errors.Is(err, localized.ErrLimitExceeded) {
			t.Fatalf("negative pair limit error = %v", err)
		}
	})
}

func TestSecurityCustomJSONBudgetRetainsParserTagCeiling(t *testing.T) {
	limits := localized.DefaultLimits()
	limits.MaxTagBytes = 512
	_, err := localized.DecodeJSON([]byte(`{"`+strings.Repeat("a", 256)+`":"text"}`), localized.DecodeOptions{Limits: limits})
	if !errors.Is(err, localized.ErrLimitExceeded) {
		t.Fatalf("custom tag budget error = %v", err)
	}
}

func TestSecurityJSONTextBudgetsWithoutCountMasking(t *testing.T) {
	limits := localized.DefaultLimits()
	limits.MaxLocales = 4
	limits.MaxTextBytes = 2
	limits.MaxTotalBytes = 5
	for _, input := range []string{
		`{"en":"big","invalid_locale":0}`,
		`{"en":"ok","fi":"ok","sv":"ok","invalid_locale":0}`,
	} {
		_, err := localized.DecodeJSON([]byte(input), localized.DecodeOptions{Limits: limits})
		if !errors.Is(err, localized.ErrLimitExceeded) {
			t.Fatalf("text budget before malformed input = %v", err)
		}
	}
	limits.MaxTotalBytes = 4
	value, err := localized.DecodeJSON([]byte(`{"en":"ok","fi":"ok"}`), localized.DecodeOptions{Limits: limits})
	if err != nil || value.Len() != 2 {
		t.Fatalf("exact text and total budgets = %v, %v", value.Entries(), err)
	}
}

func TestSecurityConstructorAccumulatesThreeTextValues(t *testing.T) {
	limits := localized.DefaultLimits()
	limits.MaxTextBytes = 2
	limits.MaxTotalBytes = 5
	_, err := localized.NewTextWithLimits(limits,
		localized.Entry{Locale: mustLocale(t, "en"), Text: "ok"},
		localized.Entry{Locale: mustLocale(t, "fi"), Text: "ok"},
		localized.Entry{Locale: mustLocale(t, "sv"), Text: "ok"},
	)
	if !errors.Is(err, localized.ErrLimitExceeded) {
		t.Fatalf("cumulative constructor bytes = %v", err)
	}
}

func TestSecurityConstructorsAcceptExactCollectionCount(t *testing.T) {
	values := make(map[string]string)
	for i := range localized.DefaultLimits().MaxLocales {
		values[fmt.Sprintf("x-entry-%d", i)] = "text"
	}
	value, err := localized.TextFromMap(values)
	if err != nil || value.Len() != len(values) {
		t.Fatalf("exact map count = %d, %v", value.Len(), err)
	}
	limits := localized.DefaultLimits()
	limits.MaxLocales = 1
	value, err = localized.TextFromPairsWithOptions(localized.ConstructionOptions{Limits: limits}, localized.Pair{Locale: "en", Text: "text"})
	if err != nil || value.Len() != 1 {
		t.Fatalf("exact pair count = %d, %v", value.Len(), err)
	}
}

func TestSecurityCollectionCountBeforeLocaleParsing(t *testing.T) {
	values := make(map[string]string)
	pairs := make([]localized.Pair, 129)
	for i := range pairs {
		key := fmt.Sprintf("invalid_%d", i)
		values[key] = "text"
		pairs[i] = localized.Pair{Locale: key, Text: "text"}
	}
	t.Run("map", func(t *testing.T) {
		if _, err := localized.TextFromMap(values); !errors.Is(err, localized.ErrLimitExceeded) {
			t.Fatalf("map error = %v", err)
		}
	})
	t.Run("pairs", func(t *testing.T) {
		if _, err := localized.TextFromPairs(pairs...); !errors.Is(err, localized.ErrLimitExceeded) {
			t.Fatalf("pairs error = %v", err)
		}
	})
}

func TestSecurityConstructionTextBytesBeforeUTF8(t *testing.T) {
	limits := localized.DefaultLimits()
	limits.MaxTextBytes = 1
	_, err := localized.NewTextWithLimits(limits, localized.Entry{Locale: mustLocale(t, "en"), Text: "\xff\xff"})
	if !errors.Is(err, localized.ErrLimitExceeded) {
		t.Fatalf("oversized text error = %v", err)
	}
}
