// Package encoding provides alternate canonical localized representations.
package encoding

import (
	"bytes"
	"encoding/json"
	"io"
	"unicode/utf8"

	localized "github.com/faustbrian/go-localized/v3"
)

const defaultMaxInputBytes = 9 << 20

// Entry is the stable entry-array wire representation.
type Entry struct {
	Locale string `json:"locale"`
	Text   string `json:"text"`
}

// DecodeOptions bounds entry-array input and constructed values.
type DecodeOptions struct {
	MaxInputBytes int
	Limits        localized.Limits
}

// MarshalEntries encodes deterministic canonical entry order.
func MarshalEntries(value localized.Text) ([]byte, error) {
	entries := value.Entries()
	wire := make([]Entry, len(entries))
	for i, entry := range entries {
		wire[i] = Entry{Locale: entry.Locale.String(), Text: entry.Text}
	}
	return json.Marshal(wire)
}

// UnmarshalEntries strictly decodes one bounded entry array.
func UnmarshalEntries(data []byte, options DecodeOptions) (localized.Text, error) {
	maxInput := options.MaxInputBytes
	if maxInput < 0 {
		return localized.Text{}, localized.ErrLimitExceeded
	}
	if maxInput == 0 {
		maxInput = defaultMaxInputBytes
	}
	if len(data) > maxInput {
		return localized.Text{}, localized.ErrLimitExceeded
	}
	if !utf8.Valid(data) {
		return localized.Text{}, localized.ErrInvalidUTF8
	}
	limits := options.Limits
	if limits == (localized.Limits{}) {
		limits = localized.DefaultLimits()
	}
	if limits.MaxLocales < 0 || limits.MaxTagBytes < 0 || limits.MaxTextBytes < 0 || limits.MaxTotalBytes < 0 {
		return localized.Text{}, localized.ErrLimitExceeded
	}
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		return localized.Text{}, localized.ErrNullValue
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	first, err := decoder.Token()
	if err != nil || first != json.Delim('[') {
		return localized.Text{}, localized.ErrInvalidEncoding
	}
	entries := make([]localized.Pair, 0)
	total := 0
	for decoder.More() {
		if len(entries) >= limits.MaxLocales {
			return localized.Text{}, localized.ErrLimitExceeded
		}
		var entry Entry
		if err := decoder.Decode(&entry); err != nil {
			return localized.Text{}, localized.ErrInvalidEncoding
		}
		if len(entry.Locale) > limits.MaxTagBytes || len(entry.Text) > limits.MaxTextBytes || len(entry.Text) > limits.MaxTotalBytes-total {
			return localized.Text{}, localized.ErrLimitExceeded
		}
		total += len(entry.Text)
		entries = append(entries, localized.Pair{Locale: entry.Locale, Text: entry.Text})
	}
	if _, err := decoder.Token(); err != nil {
		return localized.Text{}, localized.ErrInvalidEncoding
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return localized.Text{}, localized.ErrTrailingInput
		}
		return localized.Text{}, localized.ErrInvalidEncoding
	}
	return localized.TextFromPairsWithOptions(localized.ConstructionOptions{Limits: limits}, entries...)
}
