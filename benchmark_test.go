package localized_test

import (
	"fmt"
	"testing"

	localized "github.com/faustbrian/go-localized"
	localizedmatch "github.com/faustbrian/go-localized/match"
)

func benchmarkEntries(t testing.TB, size int) []localized.Entry {
	t.Helper()
	entries := make([]localized.Entry, size)
	for index := range entries {
		entries[index] = localized.Entry{
			Locale: mustLocale(t, fmt.Sprintf("en-x-%04d", index)),
			Text:   fmt.Sprintf("localized value %d", index),
		}
	}
	return entries
}

func BenchmarkConstruction(b *testing.B) {
	entries := benchmarkEntries(b, 16)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := localized.NewText(entries...); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExactLookup(b *testing.B) {
	value, _ := localized.NewText(benchmarkEntries(b, 16)...)
	locale := mustLocale(b, "en-x-0008")
	b.ReportAllocs()
	for b.Loop() {
		if _, ok := value.Get(locale); !ok {
			b.Fatal("missing")
		}
	}
}

func BenchmarkMatching(b *testing.B) {
	value, _ := localized.TextFromMap(map[string]string{"en-US": "Hello", "fi": "Hei", "sv": "Hej"})
	preference := localizedmatch.Preference{Locale: mustLocale(b, "en-CA"), Weight: 1}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := localizedmatch.Best(value, preference); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFallback(b *testing.B) {
	value, _ := localized.TextFromMap(map[string]string{"en": "Hello"})
	requested := mustLocale(b, "fi")
	plan, _ := localizedmatch.NewPlan([]localizedmatch.Chain{{From: requested, Candidates: []localizedmatch.Candidate{{Kind: localizedmatch.ExactLocale, Locale: mustLocale(b, "en")}}}}, localizedmatch.PlanOptions{MaxDepth: 4, MaxCandidates: 4})
	b.ReportAllocs()
	for b.Loop() {
		if result := plan.Resolve(value, requested); !result.Present {
			b.Fatal("missing")
		}
	}
}

func BenchmarkMerge(b *testing.B) {
	left, _ := localized.NewText(benchmarkEntries(b, 16)...)
	right, _ := localized.TextFromMap(map[string]string{"fi": "Hei", "sv": "Hej"})
	b.ReportAllocs()
	for b.Loop() {
		if _, err := left.Merge(right, localized.RightWins); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkLargeConstruction(b *testing.B) {
	entries := benchmarkEntries(b, 128)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := localized.NewText(entries...); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONEncoding(b *testing.B) {
	value, _ := localized.NewText(benchmarkEntries(b, 16)...)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := localized.EncodeJSON(value); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSONDecoding(b *testing.B) {
	value, _ := localized.NewText(benchmarkEntries(b, 16)...)
	encoded, _ := localized.EncodeJSON(value)
	b.ReportAllocs()
	for b.Loop() {
		if _, err := localized.DecodeJSON(encoded, localized.DecodeOptions{}); err != nil {
			b.Fatal(err)
		}
	}
}
