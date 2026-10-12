package localized

// Error is a stable privacy-safe error identity.
type Error string

// Error implements error.
func (e Error) Error() string { return string(e) }

const (
	// ErrConflict reports an overlapping locale under RejectConflict.
	ErrConflict Error = "localized: merge conflict"
	// ErrDuplicateLocale reports repeated canonical locale identities.
	ErrDuplicateLocale Error = "localized: duplicate locale"
	// ErrInvalidEncoding reports malformed or incompatible boundary input without its
	// content.
	ErrInvalidEncoding Error = "localized: invalid encoding"
	// ErrInvalidLocale reports a locale rejected by the pinned parser/canonicalizer.
	ErrInvalidLocale Error = "localized: invalid locale"
	// ErrInvalidPolicy reports unsupported option values or missing transform policy.
	ErrInvalidPolicy Error = "localized: invalid policy"
	// ErrInvalidUTF8 reports malformed raw UTF-8 bytes, not every escaped JSON surrogate.
	ErrInvalidUTF8 Error = "localized: invalid UTF-8"
	// ErrLimitExceeded reports construction or parser byte/count refusal.
	ErrLimitExceeded Error = "localized: limit exceeded"
	// ErrLocaleRejected reports an explicit LocalePolicy denial.
	ErrLocaleRejected Error = "localized: locale rejected"
	// ErrMissingLocale reports missing exact presence without matching or fallback.
	ErrMissingLocale Error = "localized: missing locale"
	// ErrNullValue reports null where strict localized representation is required.
	ErrNullValue Error = "localized: null value"
	// ErrResolverRequired reports ResolveConflict without a resolver.
	ErrResolverRequired Error = "localized: resolver required"
	// ErrTrailingInput reports an additional value after the decoded boundary.
	ErrTrailingInput Error = "localized: trailing input"
)
