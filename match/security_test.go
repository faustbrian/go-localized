package match_test

import (
	"errors"
	"fmt"
	"testing"

	localized "github.com/faustbrian/go-localized/v2"
	localizedmatch "github.com/faustbrian/go-localized/v2/match"
)

func TestSecurityPlanLongestDepthIndependentOfMapOrder(t *testing.T) {
	chains := []localizedmatch.Chain{
		{From: mustLocale(t, "en"), Candidates: []localizedmatch.Candidate{{Locale: mustLocale(t, "fi")}}},
		{From: mustLocale(t, "fi"), Candidates: []localizedmatch.Candidate{{Locale: mustLocale(t, "sv")}}},
		{From: mustLocale(t, "sv")},
	}
	for range 100 {
		if _, err := localizedmatch.NewPlan(chains, localizedmatch.PlanOptions{MaxDepth: 2, MaxCandidates: 3}); !errors.Is(err, localizedmatch.ErrDepthLimit) {
			t.Fatalf("longest path accepted: %v", err)
		}
	}
}

func TestSecurityPlanAcceptsExactSourceAndEdgeCaps(t *testing.T) {
	chains := []localizedmatch.Chain{
		{From: mustLocale(t, "en"), Candidates: []localizedmatch.Candidate{{Locale: mustLocale(t, "fi")}}},
		{From: mustLocale(t, "de"), Candidates: []localizedmatch.Candidate{{Locale: mustLocale(t, "sv")}}},
	}
	options := localizedmatch.PlanOptions{MaxDepth: 1, MaxCandidates: 2}
	plan, err := localizedmatch.NewPlan(chains, options)
	if err != nil {
		t.Fatal(err)
	}
	value, err := localized.NewText(localized.Entry{Locale: mustLocale(t, "sv"), Text: "Hej"})
	if err != nil {
		t.Fatal(err)
	}
	if result := plan.Resolve(value, chains[1].From); !result.Present || result.Kind != localizedmatch.Fallback || result.Text != "Hej" {
		t.Fatalf("exact graph budget resolution = %+v", result)
	}
	chains[1].Candidates = append(chains[1].Candidates, localizedmatch.Candidate{Locale: mustLocale(t, "it")})
	if _, err := localizedmatch.NewPlan(chains, options); !errors.Is(err, localizedmatch.ErrCandidateLimit) {
		t.Fatalf("aggregate edge overflow error = %v", err)
	}
}

func TestSecurityPlanEdgeBudgetBeforeCandidateValidation(t *testing.T) {
	chain := localizedmatch.Chain{From: mustLocale(t, "en"), Candidates: []localizedmatch.Candidate{{}, {}}}
	_, err := localizedmatch.NewPlan([]localizedmatch.Chain{chain}, localizedmatch.PlanOptions{MaxDepth: 1, MaxCandidates: 1})
	if !errors.Is(err, localizedmatch.ErrCandidateLimit) {
		t.Fatalf("edge budget error = %v", err)
	}
}

func TestSecurityPlanCumulativeBudgetAcrossThreeChains(t *testing.T) {
	chains := []localizedmatch.Chain{
		{From: mustLocale(t, "en"), Candidates: []localizedmatch.Candidate{{Locale: mustLocale(t, "fi")}, {Locale: mustLocale(t, "sv")}}},
		{From: mustLocale(t, "de"), Candidates: []localizedmatch.Candidate{{Locale: mustLocale(t, "it")}}},
		{From: mustLocale(t, "fr"), Candidates: []localizedmatch.Candidate{{Locale: mustLocale(t, "es")}}},
	}
	if _, err := localizedmatch.NewPlan(chains, localizedmatch.PlanOptions{MaxDepth: 1, MaxCandidates: 3}); !errors.Is(err, localizedmatch.ErrCandidateLimit) {
		t.Fatalf("four edges within three-edge budget: %v", err)
	}
}

func TestSecurityPlanDepthLimitBeforeCycleInspection(t *testing.T) {
	chains := []localizedmatch.Chain{
		{From: mustLocale(t, "en"), Candidates: []localizedmatch.Candidate{{Locale: mustLocale(t, "fi")}}},
		{From: mustLocale(t, "fi"), Candidates: []localizedmatch.Candidate{{Locale: mustLocale(t, "en")}}},
	}
	if _, err := localizedmatch.NewPlan(chains, localizedmatch.PlanOptions{MaxDepth: 1, MaxCandidates: 2}); !errors.Is(err, localizedmatch.ErrDepthLimit) {
		t.Fatalf("depth limit before traversing cycle: %v", err)
	}
}

func TestSecurityPlanBoundsSourceCountWithoutEdges(t *testing.T) {
	_, err := localizedmatch.NewPlan([]localizedmatch.Chain{{From: mustLocale(t, "en")}, {From: mustLocale(t, "fi")}}, localizedmatch.PlanOptions{MaxDepth: 1, MaxCandidates: 1})
	if !errors.Is(err, localizedmatch.ErrCandidateLimit) {
		t.Fatalf("source count error = %v", err)
	}
}

func TestSecurityPlanSharedGraphWorkIsBounded(t *testing.T) {
	const nodes = 20
	chains := make([]localizedmatch.Chain, nodes)
	for i := range chains {
		chains[i].From = mustLocale(t, fmt.Sprintf("x-node-%d", i))
		for j := i + 1; j <= i+2 && j < nodes; j++ {
			chains[i].Candidates = append(chains[i].Candidates, localizedmatch.Candidate{Locale: mustLocale(t, fmt.Sprintf("x-node-%d", j))})
		}
	}
	plan, err := localizedmatch.NewPlan(chains, localizedmatch.PlanOptions{MaxDepth: nodes, MaxCandidates: 2 * nodes})
	if err != nil {
		t.Fatal(err)
	}
	allocations := testing.AllocsPerRun(1, func() {
		if result := plan.Resolve(localized.Text{}, chains[0].From); result.Present {
			t.Fatal("missing graph resolved")
		}
	})
	// A linear walk needs at most a small fixed amount per edge. Revisiting
	// the shared Fibonacci-shaped graph costs tens of thousands of allocations.
	if allocations > 1000 {
		t.Fatalf("shared graph allocated %.0f times for %d nodes", allocations, nodes)
	}
}
