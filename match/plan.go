package match

import (
	"fmt"

	"github.com/faustbrian/go-international/v4/locale"
	localized "github.com/faustbrian/go-localized/v5"
)

// CandidateKind selects exact lookup or locale-layer parent traversal.
type CandidateKind uint8

const (
	// ExactLocale checks only the candidate's canonical tag.
	ExactLocale CandidateKind = iota
	// ParentRange walks Tag.Parent according to the locale dependency.
	ParentRange
)

// Candidate is one ordered configured fallback operation.
type Candidate struct {
	// Kind selects exact identity or strict-parent traversal before following a configured
	// chain.
	Kind CandidateKind
	// Locale is canonicalized and copied during plan construction.
	Locale locale.Tag
}

// Chain configures candidates to try after From is absent.
type Chain struct {
	// From is the unique canonical source identity for this fallback chain.
	From locale.Tag
	// Candidates is an ordered copied list; duplicate Kind/Locale pairs are rejected.
	Candidates []Candidate
}

// PlanOptions bounds a graph and provides an optional final default.
// MaxCandidates separately bounds source-chain count and candidate-edge count.
type PlanOptions struct {
	// Default is a copied optional final exact locale, not charged against candidate-edge
	// count.
	Default *locale.Tag
	// MaxDepth is a required positive configured-chain path ceiling, not a locale-parent
	// or callback time budget.
	MaxDepth int
	// MaxCandidates separately caps source chains and total candidate edges; zero allows
	// neither.
	MaxCandidates int
	// Observer is borrowed for the plan lifetime and can be invoked concurrently by
	// readers.
	Observer Observer
}

// Plan is an immutable validated fallback graph.
type Plan struct {
	chains     map[string]Chain
	defaultTag *locale.Tag
	observer   Observer
	candidates int
}

// NewPlan admits source/edge counts before copying, canonicalizes candidates and
// validates cycles and longest configured-chain depth. The optional default is
// copied but not charged as a candidate edge. Observer is borrowed for the plan
// lifetime; immutable readers may invoke it concurrently. No locale-parent depth
// or callback-time budget is added by these graph limits.
func NewPlan(chains []Chain, options PlanOptions) (Plan, error) {
	if options.MaxDepth <= 0 {
		return Plan{}, ErrDepthLimit
	}
	if options.MaxCandidates < 0 {
		return Plan{}, ErrCandidateLimit
	}
	if len(chains) > options.MaxCandidates {
		return Plan{}, ErrCandidateLimit
	}
	// Reject the complete edge budget before copying any caller container.
	remaining := options.MaxCandidates
	for _, chain := range chains {
		if len(chain.Candidates) > remaining {
			return Plan{}, ErrCandidateLimit
		}
		remaining -= len(chain.Candidates)
	}
	owned := make(map[string]Chain, len(chains))
	total := 0
	for _, chain := range chains {
		canonicalFrom, err := chain.From.Canonical()
		if err != nil {
			return Plan{}, ErrInvalidCandidate
		}
		chain.From = canonicalFrom
		key := canonicalFrom.String()
		if _, exists := owned[key]; exists {
			return Plan{}, ErrDuplicateCandidate
		}
		copy := Chain{From: chain.From, Candidates: append([]Candidate(nil), chain.Candidates...)}
		seen := make(map[string]struct{}, len(copy.Candidates))
		for i, candidate := range copy.Candidates {
			if candidate.Kind > ParentRange {
				return Plan{}, ErrInvalidCandidate
			}
			canonical, err := candidate.Locale.Canonical()
			if err != nil {
				return Plan{}, ErrInvalidCandidate
			}
			candidate.Locale = canonical
			copy.Candidates[i] = candidate
			candidateKey := fmt.Sprintf("%d:%s", candidate.Kind, candidate.Locale.String())
			if _, exists := seen[candidateKey]; exists {
				return Plan{}, ErrDuplicateCandidate
			}
			seen[candidateKey] = struct{}{}
			total++
		}
		owned[key] = copy
	}
	if err := validateGraph(owned, options.MaxDepth); err != nil {
		return Plan{}, err
	}
	var defaultTag *locale.Tag
	if options.Default != nil {
		copy, err := options.Default.Canonical()
		if err != nil {
			return Plan{}, ErrInvalidCandidate
		}
		defaultTag = &copy
	}
	return Plan{
		chains: owned, defaultTag: defaultTag,
		observer: options.Observer, candidates: total,
	}, nil
}

func validateGraph(chains map[string]Chain, maxDepth int) error {
	const (
		unseen uint8 = iota
		active
		done
	)
	states := make(map[string]uint8, len(chains))
	heights := make(map[string]int, len(chains))
	var visit func(string, int) (int, error)
	visit = func(key string, depth int) (int, error) {
		if depth > maxDepth {
			return 0, ErrDepthLimit
		}
		if states[key] == active {
			return 0, ErrFallbackCycle
		}
		if states[key] == done {
			return heights[key], nil
		}
		states[key] = active
		height := 1
		for _, candidate := range chains[key].Candidates {
			next := candidate.Locale.String()
			if candidate.Kind == ParentRange && next == key {
				continue
			}
			if _, exists := chains[next]; exists {
				childHeight, err := visit(next, depth+1)
				if err != nil {
					return 0, err
				}
				height = max(height, 1+childHeight)
				if height > maxDepth {
					return 0, ErrDepthLimit
				}
			}
		}
		states[key] = done
		heights[key] = height
		return height, nil
	}
	for key := range chains {
		if _, err := visit(key, 1); err != nil {
			return err
		}
	}
	return nil
}

// Resolve checks the requested locale, traverses its configured graph, and
// finally checks the optional default without changing the source value. Empty
// strings are successful presence. ParentRange starts with the strict parent of
// its candidate, excluding the candidate itself; parent identities come from the
// pinned dependency. Invalid requests return Missing without observation. The
// zero Plan still performs exact lookup, with no configured fallback/default.
func (p Plan) Resolve(value localized.Text, requested locale.Tag) Result {
	canonicalRequested, err := requested.Canonical()
	if err != nil {
		return Result{Kind: Missing}
	}
	requested = canonicalRequested
	var result Result
	if text, ok := value.Get(requested); ok {
		result = present(Exact, requested, requested, text)
		notify(p.observer, Event{Operation: OperationFallback, Kind: result.Kind, CandidateCount: p.candidates})
		return result
	}
	if resolved, ok := p.resolveChain(value, requested, requested, make(map[string]bool, len(p.chains))); ok {
		result = resolved
		notify(p.observer, Event{Operation: OperationFallback, Kind: result.Kind, CandidateCount: p.candidates})
		return result
	}
	if p.defaultTag != nil {
		if text, ok := value.Get(*p.defaultTag); ok {
			result = present(Default, requested, *p.defaultTag, text)
			notify(p.observer, Event{Operation: OperationFallback, Kind: result.Kind, CandidateCount: p.candidates})
			return result
		}
	}
	result = Result{Kind: Missing, Requested: requested}
	notify(p.observer, Event{Operation: OperationFallback, Kind: result.Kind, CandidateCount: p.candidates})
	return result
}

func (p Plan) resolveChain(value localized.Text, requested, current locale.Tag, visited map[string]bool) (Result, bool) {
	key := current.String()
	if visited[key] {
		return Result{}, false
	}
	chain, exists := p.chains[key]
	if !exists {
		return Result{}, false
	}
	// A shared subgraph has the same outcome on every path for this immutable
	// value. Visit it once to prevent exponential work on missing values.
	visited[key] = true
	for _, candidate := range chain.Candidates {
		switch candidate.Kind {
		case ExactLocale:
			if text, ok := value.Get(candidate.Locale); ok {
				return present(Fallback, requested, candidate.Locale, text), true
			}
		case ParentRange:
			for parent, ok := candidate.Locale.Fallback(locale.FallbackParent); ok; parent, ok = parent.Fallback(locale.FallbackParent) {
				if text, ok := value.Get(parent); ok {
					return present(Fallback, requested, parent, text), true
				}
			}
		}
		if candidate.Locale.String() == current.String() {
			continue
		}
		if result, ok := p.resolveChain(value, requested, candidate.Locale, visited); ok {
			return result, true
		}
	}
	return Result{}, false
}
