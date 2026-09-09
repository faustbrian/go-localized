package config

import (
	localized "github.com/faustbrian/go-localized"
	legacy "github.com/faustbrian/go-localized/localizedconfig"
)

// Error preserves the released localizedconfig error identity.
type Error = legacy.Error

// ErrInvalidValue reports a non-string-map configuration value.
const ErrInvalidValue = legacy.ErrInvalidValue

// Text preserves the released localizedconfig value-hook identity and
// behavior. It snapshots decoded maps through localized.Text construction.
type Text = legacy.Text

// NewText creates a present configuration wrapper.
func NewText(value localized.Text) Text { return legacy.NewText(value) }
