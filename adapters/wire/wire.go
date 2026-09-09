package wire

import (
	localized "github.com/faustbrian/go-localized"
	legacy "github.com/faustbrian/go-localized/localizedwire"
	"github.com/faustbrian/go-wire/jsonwire"
	"github.com/faustbrian/go-wire/msgpackwire"
	"github.com/faustbrian/go-wire/tomlwire"
	"github.com/faustbrian/go-wire/yamlwire"
)

// EncodeJSON encodes canonical keys through wire's bounded JSON encoder.
func EncodeJSON(value localized.Text, options jsonwire.EncodeOptions) ([]byte, error) {
	return legacy.EncodeJSON(value, options)
}

// DecodeJSON validates wire boundaries and localized JSON semantics.
func DecodeJSON(data []byte, options jsonwire.DecodeOptions) (localized.Text, error) {
	return legacy.DecodeJSON(data, options)
}

// EncodeYAML encodes canonical keys through wire's bounded YAML encoder.
func EncodeYAML(value localized.Text, options yamlwire.EncodeOptions) ([]byte, error) {
	return legacy.EncodeYAML(value, options)
}

// DecodeYAML decodes a strict string map and validates canonical locale keys.
func DecodeYAML(data []byte, options yamlwire.DecodeOptions) (localized.Text, error) {
	return legacy.DecodeYAML(data, options)
}

// EncodeTOML encodes canonical keys through wire's bounded TOML encoder.
func EncodeTOML(value localized.Text, options tomlwire.EncodeOptions) ([]byte, error) {
	return legacy.EncodeTOML(value, options)
}

// DecodeTOML decodes a string map and validates canonical locale keys.
func DecodeTOML(data []byte, options tomlwire.DecodeOptions) (localized.Text, error) {
	return legacy.DecodeTOML(data, options)
}

// EncodeMessagePack encodes canonical keys through wire's deterministic
// bounded MessagePack encoder.
func EncodeMessagePack(value localized.Text, options msgpackwire.EncodeOptions) ([]byte, error) {
	return legacy.EncodeMessagePack(value, options)
}

// DecodeMessagePack decodes a strict string map and validates locale keys.
func DecodeMessagePack(data []byte, options msgpackwire.DecodeOptions) (localized.Text, error) {
	return legacy.DecodeMessagePack(data, options)
}
