// Package localizedconfig adapts Text to config's value-hook contracts.
// Hooks mutate caller-owned nullable wrappers, snapshot input maps on success and
// retain immutable payloads. No configuration source, reload or watcher is owned.
package localizedconfig

import (
	"reflect"

	localized "github.com/faustbrian/go-localized/v5"
)

// Error is a stable privacy-safe configuration error identity.
type Error string

// Error implements error.
func (e Error) Error() string { return string(e) }

// ErrInvalidValue reports a non-string-map configuration value.
const ErrInvalidValue Error = "localized config: invalid value"

// Text is a presence-aware config value hook.
type Text struct {
	// Localized is the immutable payload; meaningful only when Valid is true.
	Localized localized.Text
	// Valid distinguishes present empty content from null/absent wrapper state; it is
	// caller-owned mutable state.
	Valid bool
}

// NewText creates a present configuration wrapper.
func NewText(value localized.Text) Text {
	return Text{Localized: value, Valid: true}
}

// ConfigTextTarget asks textual sources to decode a string map before the hook.
func (Text) ConfigTextTarget() reflect.Type {
	return reflect.TypeFor[map[string]string]()
}

// UnmarshalConfigValue snapshots string-valued maps under default construction
// limits and assigns only on success; prior state survives failure. Untyped nil
// clears to invalid, while a typed nil map constructs present-empty Text. A nil
// receiver returns ErrInvalidValue; input maps must not change during the call.
func (t *Text) UnmarshalConfigValue(input any) error {
	if t == nil {
		return ErrInvalidValue
	}
	if input == nil {
		*t = Text{}
		return nil
	}
	values := make(map[string]string)
	switch typed := input.(type) {
	case map[string]string:
		if len(typed) > localized.DefaultLimits().MaxLocales {
			return localized.ErrLimitExceeded
		}
		for key, value := range typed {
			values[key] = value
		}
	case map[string]any:
		if len(typed) > localized.DefaultLimits().MaxLocales {
			return localized.ErrLimitExceeded
		}
		for key, value := range typed {
			text, ok := value.(string)
			if !ok {
				return ErrInvalidValue
			}
			values[key] = text
		}
	default:
		return ErrInvalidValue
	}
	decoded, err := localized.TextFromMap(values)
	if err != nil {
		return err
	}
	*t = NewText(decoded)
	return nil
}
