package localized

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/faustbrian/go-international/v3/locale"
)

// JSONMode controls legacy compatibility during decoding.
type JSONMode uint8

const (
	// StrictJSON accepts only a JSON object with strict BCP 47 keys.
	StrictJSON JSONMode = iota
	// PermissiveJSON additionally accepts null and underscore locale separators.
	PermissiveJSON
)

const defaultMaxJSONBytes = 9437184

// DecodeOptions bounds JSON parsing and selects strictness.
type DecodeOptions struct {
	Mode          JSONMode
	MaxInputBytes int
	Limits        Limits
}

// EncodeJSON returns the deterministic canonical JSON object representation.
func EncodeJSON(value Text) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteByte('{')
	for i, entry := range value.entries {
		if i > 0 {
			buffer.WriteByte(',')
		}
		key, _ := json.Marshal(entry.Locale.String())
		text, _ := json.Marshal(entry.Text)
		buffer.Write(key)
		buffer.WriteByte(':')
		buffer.Write(text)
	}
	buffer.WriteByte('}')
	return buffer.Bytes(), nil
}

// DecodeJSON parses one bounded localized JSON object without partial results.
func DecodeJSON(data []byte, options DecodeOptions) (Text, error) {
	if options.Mode > PermissiveJSON {
		return Text{}, ErrInvalidPolicy
	}
	maxInput := options.MaxInputBytes
	if maxInput < 0 {
		return Text{}, fmt.Errorf("%w: parser input", ErrLimitExceeded)
	}
	if maxInput == 0 {
		maxInput = defaultMaxJSONBytes
	}
	if len(data) > maxInput {
		return Text{}, fmt.Errorf("%w: parser input", ErrLimitExceeded)
	}
	if !utf8.Valid(data) {
		return Text{}, ErrInvalidUTF8
	}
	limits, err := constructionLimits(options.Limits)
	if err != nil {
		return Text{}, err
	}
	trimmed := bytes.TrimSpace(data)
	if bytes.Equal(trimmed, []byte("null")) {
		if options.Mode == PermissiveJSON {
			return Text{}, nil
		}
		return Text{}, ErrNullValue
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	first, err := decoder.Token()
	if err != nil {
		return Text{}, ErrInvalidEncoding
	}
	delimiter, ok := first.(json.Delim)
	if !ok || delimiter != '{' {
		return Text{}, ErrInvalidEncoding
	}

	entries := make([]Entry, 0)
	total := 0
	for decoder.More() {
		if len(entries) >= limits.MaxLocales {
			return Text{}, fmt.Errorf("%w: locale count", ErrLimitExceeded)
		}
		keyToken, err := decoder.Token()
		if err != nil {
			return Text{}, ErrInvalidEncoding
		}
		raw := keyToken.(string)
		if len(raw) > limits.MaxTagBytes {
			return Text{}, fmt.Errorf("%w: tag bytes", ErrLimitExceeded)
		}
		tag, err := parseJSONLocale(raw, options.Mode)
		if err != nil {
			return Text{}, err
		}
		valueToken, err := decoder.Token()
		if err != nil {
			return Text{}, ErrInvalidEncoding
		}
		value, ok := valueToken.(string)
		if !ok {
			return Text{}, ErrInvalidEncoding
		}
		if len(value) > limits.MaxTextBytes || len(value) > limits.MaxTotalBytes-total {
			return Text{}, fmt.Errorf("%w: text bytes", ErrLimitExceeded)
		}
		total += len(value)
		entries = append(entries, Entry{Locale: tag, Text: value})
	}
	if _, err := decoder.Token(); err != nil {
		return Text{}, ErrInvalidEncoding
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return Text{}, ErrTrailingInput
		}
		return Text{}, ErrInvalidEncoding
	}
	return NewTextWithLimits(limits, entries...)
}

func parseJSONLocale(raw string, mode JSONMode) (locale.Tag, error) {
	if mode == StrictJSON && strings.ContainsAny(raw, "_ \t\r\n") {
		return locale.Tag{}, ErrInvalidLocale
	}
	if mode == PermissiveJSON {
		raw = strings.ReplaceAll(raw, "_", "-")
	}
	if len(raw) > defaultMaxTagBytes {
		return locale.Tag{}, fmt.Errorf("%w: tag bytes", ErrLimitExceeded)
	}
	tag, err := locale.Parse(raw)
	if err != nil {
		return locale.Tag{}, ErrInvalidLocale
	}
	return tag, nil
}

// MarshalJSON implements json.Marshaler with canonical object encoding.
func (t Text) MarshalJSON() ([]byte, error) { return EncodeJSON(t) }

// UnmarshalJSON implements json.Unmarshaler using strict bounded decoding.
func (t *Text) UnmarshalJSON(data []byte) error {
	if t == nil {
		return ErrInvalidEncoding
	}
	decoded, err := DecodeJSON(data, DecodeOptions{})
	if err != nil {
		return err
	}
	*t = decoded
	return nil
}
