package match

// Error is a stable privacy-safe matching error identity.
type Error string

// Error implements error.
func (e Error) Error() string { return string(e) }

const (
	// ErrCandidateLimit reports a preference/chain/candidate admission count failure.
	ErrCandidateLimit Error = "localized match: candidate limit exceeded"
	// ErrDepthLimit reports invalid or excessive configured-chain depth.
	ErrDepthLimit Error = "localized match: fallback depth exceeded"
	// ErrDuplicateCandidate reports repeated canonical candidates or source chains.
	ErrDuplicateCandidate Error = "localized match: duplicate candidate"
	// ErrFallbackCycle reports a cycle in configured fallback chains.
	ErrFallbackCycle Error = "localized match: fallback cycle"
	// ErrInvalidCandidate reports an invalid locale or fallback candidate kind.
	ErrInvalidCandidate Error = "localized match: invalid fallback candidate"
	// ErrInvalidWeight reports a nonfinite or out-of-range quality weight.
	ErrInvalidWeight Error = "localized match: invalid weight"
)
