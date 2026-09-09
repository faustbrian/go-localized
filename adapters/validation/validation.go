package validation

import (
	localized "github.com/faustbrian/go-localized"
	legacy "github.com/faustbrian/go-localized/localizedvalidation"
	validationcore "github.com/faustbrian/go-validation"
)

// Error preserves the released localizedvalidation error identity.
type Error = legacy.Error

const (
	// ErrBytesExceeded reports text exceeding a byte limit.
	ErrBytesExceeded = legacy.ErrBytesExceeded
	// ErrControlCharacter reports a prohibited Unicode control character.
	ErrControlCharacter = legacy.ErrControlCharacter
	// ErrEmpty reports a present-empty localized string.
	ErrEmpty = legacy.ErrEmpty
	// ErrInvalidForm reports an unsupported normalization form.
	ErrInvalidForm = legacy.ErrInvalidForm
	// ErrInvalidRule reports a nil or invalid rule.
	ErrInvalidRule = legacy.ErrInvalidRule
	// ErrLinesExceeded reports text exceeding a line limit.
	ErrLinesExceeded = legacy.ErrLinesExceeded
	// ErrRunesExceeded reports text exceeding a rune limit.
	ErrRunesExceeded = legacy.ErrRunesExceeded
	// ErrWhitespace reports text containing only Unicode whitespace.
	ErrWhitespace = legacy.ErrWhitespace
)

// Rule preserves the released localizedvalidation rule contract.
type Rule = legacy.Rule

// Form preserves the released localizedvalidation normalization identity.
type Form = legacy.Form

const (
	// NFC applies canonical composition.
	NFC = legacy.NFC
	// NFD applies canonical decomposition.
	NFD = legacy.NFD
	// NFKC applies compatibility composition.
	NFKC = legacy.NFKC
	// NFKD applies compatibility decomposition.
	NFKD = legacy.NFKD
)

// Validate applies every rule to every present value in deterministic order.
func Validate(value localized.Text, rules ...Rule) error {
	return legacy.Validate(value, rules...)
}

// Validator adapts immutable localized text rules to go-validation.
func Validator(rules ...Rule) validationcore.Validator[localized.Text] {
	return legacy.Validator(rules...)
}

// RequireNonEmpty rejects present-empty strings.
func RequireNonEmpty() Rule { return legacy.RequireNonEmpty() }

// RequireNonWhitespace rejects empty and Unicode whitespace-only strings.
func RequireNonWhitespace() Rule { return legacy.RequireNonWhitespace() }

// MaxBytes bounds UTF-8 encoded bytes.
func MaxBytes(maximum int) Rule { return legacy.MaxBytes(maximum) }

// MaxRunes bounds Unicode code points without claiming grapheme semantics.
func MaxRunes(maximum int) Rule { return legacy.MaxRunes(maximum) }

// MaxLines bounds newline-delimited logical lines.
func MaxLines(maximum int) Rule { return legacy.MaxLines(maximum) }

// NoControlCharacters rejects Unicode controls except tab and line endings.
func NoControlCharacters() Rule { return legacy.NoControlCharacters() }

// Normalize returns a new value with the selected Unicode transform.
func Normalize(value localized.Text, form Form) (localized.Text, error) {
	return legacy.Normalize(value, form)
}
