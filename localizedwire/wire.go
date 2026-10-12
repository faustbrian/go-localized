// Package localizedwire adapts Text to the pinned wire format packages.
// Inputs are borrowed during synchronous decoding; returned Text owns its
// containers. Upstream options/errors retain upstream behavior and may expose
// input details. Localized construction uses default value limits, not a new
// total heap/time budget; these helpers own no transport or I/O lifecycle.
package localizedwire

import (
	localized "github.com/faustbrian/go-localized/v5"
	"github.com/faustbrian/go-wire/v3/jsonwire"
	"github.com/faustbrian/go-wire/v3/msgpackwire"
	"github.com/faustbrian/go-wire/v3/tomlwire"
	"github.com/faustbrian/go-wire/v3/yamlwire"
)

func stringMap(value localized.Text) map[string]string {
	result := make(map[string]string, value.Len())
	for _, entry := range value.Entries() {
		result[entry.Locale.String()] = entry.Text
	}
	return result
}

// EncodeJSON encodes canonical keys through wire's bounded JSON encoder.
func EncodeJSON(value localized.Text, options jsonwire.EncodeOptions) ([]byte, error) {
	return jsonwire.Encode(stringMap(value), options)
}

// DecodeJSON validates wire boundaries, then applies canonical duplicate and
// locale semantics from the localized JSON decoder. The raw bytes are decoded
// twice: upstream shape admission, then strict core decoding under default value
// limits and the converted MaxBytes input bound. Neither pass retains input.
func DecodeJSON(data []byte, options jsonwire.DecodeOptions) (localized.Text, error) {
	var shape map[string]string
	if err := jsonwire.Decode(data, &shape, options); err != nil {
		return localized.Text{}, err
	}
	maxBytes := int(options.MaxBytes)
	return localized.DecodeJSON(data, localized.DecodeOptions{MaxInputBytes: maxBytes})
}

// EncodeYAML encodes canonical keys through wire's bounded YAML encoder.
func EncodeYAML(value localized.Text, options yamlwire.EncodeOptions) ([]byte, error) {
	return yamlwire.Encode(stringMap(value), options)
}

// DecodeYAML delegates string-map parsing to wire then applies default localized
// construction. Duplicate raw-key handling is upstream; map conversion cannot
// recover occurrences already collapsed by the upstream decoder.
func DecodeYAML(data []byte, options yamlwire.DecodeOptions) (localized.Text, error) {
	var values map[string]string
	if err := yamlwire.Decode(data, &values, options); err != nil {
		return localized.Text{}, err
	}
	return localized.TextFromMap(values)
}

// EncodeTOML encodes canonical keys through wire's bounded TOML encoder.
func EncodeTOML(value localized.Text, options tomlwire.EncodeOptions) ([]byte, error) {
	return tomlwire.Encode(stringMap(value), options)
}

// DecodeTOML delegates string-map parsing to wire then applies default localized
// construction. Raw duplicate-key visibility and parser policy belong to wire.
func DecodeTOML(data []byte, options tomlwire.DecodeOptions) (localized.Text, error) {
	var values map[string]string
	if err := tomlwire.Decode(data, &values, options); err != nil {
		return localized.Text{}, err
	}
	return localized.TextFromMap(values)
}

// EncodeMessagePack encodes canonical keys through wire's deterministic
// bounded MessagePack encoder.
func EncodeMessagePack(value localized.Text, options msgpackwire.EncodeOptions) ([]byte, error) {
	return msgpackwire.Encode(stringMap(value), options)
}

// DecodeMessagePack delegates string-map parsing to wire then applies default
// localized construction. Raw duplicate-key visibility belongs to wire; canonical
// duplicates still visible in distinct map keys are rejected by core.
func DecodeMessagePack(data []byte, options msgpackwire.DecodeOptions) (localized.Text, error) {
	var values map[string]string
	if err := msgpackwire.Decode(data, &values, options); err != nil {
		return localized.Text{}, err
	}
	return localized.TextFromMap(values)
}
