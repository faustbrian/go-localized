package match

// Operation identifies a bounded resolution operation.
type Operation uint8

const (
	// OperationMatch identifies standards-aligned preference matching.
	OperationMatch Operation = iota
	// OperationFallback identifies configured fallback resolution.
	OperationFallback
)

// Event contains bounded low-cardinality resolution telemetry. It intentionally
// excludes requested, resolved, and content values.
type Event struct {
	// Operation identifies preference matching or configured fallback.
	Operation Operation
	// Kind is the completed resolution classification.
	Kind Kind
	// CandidateCount is accepted preference count or configured edge count, not actual
	// probes or locale identities.
	CandidateCount int
}

// Observer receives completed resolution events. notify isolates panics but
// calls remain synchronous; blocking observers delay return and concurrent
// readers require a concurrency-safe implementation. No lock is held or worker
// detached, and external callback effects remain caller-owned.
type Observer interface {
	// Observe receives a completed content-free event synchronously. Implementations
	// own bounded execution and concurrent-call safety.
	Observe(Event)
}

// ObserverFunc adapts a function to Observer.
type ObserverFunc func(Event)

// Observe implements Observer.
func (function ObserverFunc) Observe(event Event) { function(event) }

func notify(observer Observer, event Event) {
	if observer == nil {
		return
	}
	func() {
		defer func() { _ = recover() }()
		observer.Observe(event)
	}()
}
