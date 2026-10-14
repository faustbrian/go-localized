package localized_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/faustbrian/go-international/v4/locale"
	localized "github.com/faustbrian/go-localized/v5"
)

func TestSubsetTransformsPreserveCustomLimitValues(t *testing.T) {
	defaults := localized.DefaultLimits()
	removed := mustLocale(t, "fi")
	missing := mustLocale(t, "sv")
	large := strings.Repeat("x", defaults.MaxTextBytes)
	cases := []struct {
		name    string
		entries []localized.Entry
	}{
		{"text bytes", []localized.Entry{{Locale: mustLocale(t, "en"), Text: large + "x"}}},
		{"locale count", nil},
		{"total bytes", nil},
	}
	for i := 0; i <= defaults.MaxLocales; i++ {
		cases[1].entries = append(cases[1].entries, localized.Entry{
			Locale: mustLocale(t, fmt.Sprintf("x-test%03d", i)), Text: "value",
		})
	}
	for i := 0; i <= defaults.MaxTotalBytes/defaults.MaxTextBytes; i++ {
		cases[2].entries = append(cases[2].entries, localized.Entry{
			Locale: mustLocale(t, fmt.Sprintf("x-test%03d", i)), Text: large,
		})
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			// Retain present-empty content as well as the entry selected for removal.
			entries := append(test.entries, localized.Entry{Locale: removed, Text: ""})
			limits := defaults
			limits.MaxLocales = len(entries)
			limits.MaxTextBytes++
			limits.MaxTotalBytes = 0
			for _, entry := range entries {
				limits.MaxTotalBytes += len(entry.Text)
			}
			value, err := localized.NewTextWithLimits(limits, entries...)
			if err != nil {
				t.Fatalf("construct custom-limit value: %v", err)
			}
			if _, err := localized.NewText(entries...); !errors.Is(err, localized.ErrLimitExceeded) {
				t.Fatalf("default construction error = %v, want ErrLimitExceeded", err)
			}
			original := value.Entries()
			retained := make([]localized.Entry, 0, len(original)-1)
			for _, entry := range original {
				if entry.Locale != removed {
					retained = append(retained, entry)
				}
			}
			operations := []struct {
				name string
				run  func() localized.Text
				want []localized.Entry
			}{
				{"remove", func() localized.Text { return value.Remove(removed) }, retained},
				{"filter all", func() localized.Text {
					return value.Filter(func(locale.Tag, string) bool { return true })
				}, original},
				{"filter subset", func() localized.Text {
					return value.Filter(func(tag locale.Tag, _ string) bool { return tag != removed })
				}, retained},
				{"remove absent", func() localized.Text { return value.Remove(missing) }, original},
				{"remove invalid", func() localized.Text { return value.Remove(locale.Tag{}) }, original},
				{"filter nil", func() localized.Text { return value.Filter(nil) }, original},
				{"filter none", func() localized.Text {
					return value.Filter(func(locale.Tag, string) bool { return false })
				}, nil},
			}
			for _, operation := range operations {
				t.Run(operation.name, func(t *testing.T) {
					got := operation.run()
					if !reflect.DeepEqual(got.Entries(), operation.want) {
						t.Fatalf("retained entries differ: got %d, want %d", got.Len(), len(operation.want))
					}
					for _, entry := range operation.want {
						if text, present := got.Get(entry.Locale); !present || text != entry.Text {
							t.Fatalf("retained lookup %s: present=%t, bytes=%d, want %d", entry.Locale, present, len(text), len(entry.Text))
						}
					}
					if !reflect.DeepEqual(value.Entries(), original) {
						t.Fatal("subset transform mutated its receiver")
					}
					for _, entry := range original {
						if text, present := value.Get(entry.Locale); !present || text != entry.Text {
							t.Fatalf("subset transform changed original lookup %s", entry.Locale)
						}
					}
				})
			}
		})
	}
}
