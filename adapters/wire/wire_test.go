package wire_test

import (
	"testing"

	localized "github.com/faustbrian/go-localized"
	wire "github.com/faustbrian/go-localized/adapters/wire"
	"github.com/faustbrian/go-wire/jsonwire"
	"github.com/faustbrian/go-wire/msgpackwire"
	"github.com/faustbrian/go-wire/tomlwire"
	"github.com/faustbrian/go-wire/yamlwire"
)

func TestCanonicalWireAdapterPreservesAllFormats(t *testing.T) {
	t.Parallel()

	value, _ := localized.TextFromMap(map[string]string{"en": "Hello"})
	tests := []struct {
		name   string
		encode func(localized.Text) ([]byte, error)
		decode func([]byte) (localized.Text, error)
	}{
		{"json", func(value localized.Text) ([]byte, error) { return wire.EncodeJSON(value, jsonwire.EncodeOptions{}) }, func(data []byte) (localized.Text, error) { return wire.DecodeJSON(data, jsonwire.DecodeOptions{}) }},
		{"yaml", func(value localized.Text) ([]byte, error) { return wire.EncodeYAML(value, yamlwire.EncodeOptions{}) }, func(data []byte) (localized.Text, error) { return wire.DecodeYAML(data, yamlwire.DecodeOptions{}) }},
		{"toml", func(value localized.Text) ([]byte, error) { return wire.EncodeTOML(value, tomlwire.EncodeOptions{}) }, func(data []byte) (localized.Text, error) { return wire.DecodeTOML(data, tomlwire.DecodeOptions{}) }},
		{"messagepack", func(value localized.Text) ([]byte, error) {
			return wire.EncodeMessagePack(value, msgpackwire.EncodeOptions{})
		}, func(data []byte) (localized.Text, error) {
			return wire.DecodeMessagePack(data, msgpackwire.DecodeOptions{})
		}},
	}
	for _, test := range tests {
		encoded, err := test.encode(value)
		if err != nil {
			t.Fatalf("%s encode error = %v", test.name, err)
		}
		decoded, err := test.decode(encoded)
		if err != nil || !decoded.Equal(value) {
			t.Fatalf("%s decode = %v, %v", test.name, decoded.Entries(), err)
		}
	}
}
