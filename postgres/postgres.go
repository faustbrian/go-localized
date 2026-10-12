// Package postgres provides explicit SQL values and pgx JSONB codecs.
// It owns no database connection, pool, transaction, query or migration lifecycle.
// Already-fetched data is decoded synchronously; database cancellation and ingress
// byte admission remain caller-owned. Nullable wrappers and codecs are mutable
// caller-owned objects, unlike their immutable localized payloads.
package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"strings"

	"github.com/faustbrian/go-international/v4/locale"
	localized "github.com/faustbrian/go-localized/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// Error is a stable privacy-safe persistence error identity.
type Error string

// Error implements error.
func (e Error) Error() string { return string(e) }

// ErrUnsupportedDatabaseType reports a value outside JSON text forms.
const ErrUnsupportedDatabaseType Error = "localized postgres: unsupported database type"

// Text is a nullable localized value for database/sql boundaries.
type Text struct {
	// Localized is the immutable payload; meaningful only when Valid is true.
	Localized localized.Text
	// Valid distinguishes present empty content from null/absent wrapper state; it is
	// caller-owned mutable state.
	Valid bool
}

// NewText creates a non-null SQL value.
func NewText(value localized.Text) Text {
	return Text{Localized: value, Valid: true}
}

// Value implements driver.Valuer. Invalid values become SQL NULL.
func (t Text) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return localized.EncodeJSON(t.Localized)
}

// Scan assigns only after successful strict default-limit JSON decoding; SQL NULL
// resets the wrapper. Byte input is borrowed during decoding, not retained. String
// input is copied before byte admission, so cap driver row sizes before this call.
// JSON null is not SQL NULL and is rejected. A nil receiver returns a fixed error.
func (t *Text) Scan(source any) error {
	if t == nil {
		return ErrUnsupportedDatabaseType
	}
	if source == nil {
		*t = Text{}
		return nil
	}
	var data []byte
	switch value := source.(type) {
	case []byte:
		data = value
	case string:
		data = []byte(value)
	default:
		return ErrUnsupportedDatabaseType
	}
	decoded, err := localized.DecodeJSON(data, localized.DecodeOptions{})
	if err != nil {
		return err
	}
	*t = NewText(decoded)
	return nil
}

// JSONBCodec returns a pgx codec that preserves generic JSONB behavior while
// using localized's canonical encoding and strict decoding for Text values. Each
// call returns a new mutable codec; callers own registration and synchronization.
// A nil *localized.Text encodes JSON null, which strict localized decoding rejects.
// Generic branches use encoding/json and do not inherit core content-free errors.
func JSONBCodec() *pgtype.JSONBCodec {
	return &pgtype.JSONBCodec{
		Marshal: func(value any) ([]byte, error) {
			switch typed := value.(type) {
			case localized.Text:
				return localized.EncodeJSON(typed)
			case *localized.Text:
				if typed == nil {
					return []byte("null"), nil
				}
				return localized.EncodeJSON(*typed)
			default:
				return json.Marshal(value)
			}
		},
		Unmarshal: func(data []byte, target any) error {
			if typed, ok := target.(*localized.Text); ok {
				if typed == nil {
					return ErrUnsupportedDatabaseType
				}
				decoded, err := localized.DecodeJSON(data, localized.DecodeOptions{})
				if err != nil {
					return err
				}
				*typed = decoded
				return nil
			}
			return json.Unmarshal(data, target)
		},
	}
}

// Row is a normalized persistence representation with a canonical locale key.
type Row struct {
	// Locale is the canonical row key on output and a strict raw key on input.
	Locale string
	// Text is the stored UTF-8 content, including empty text.
	Text string
}

// Rows returns deterministic caller-owned normalized rows.
func Rows(value localized.Text) []Row {
	entries := value.Entries()
	rows := make([]Row, len(entries))
	for i, entry := range entries {
		rows[i] = Row{Locale: entry.Locale.String(), Text: entry.Text}
	}
	return rows
}

// FromRows admits row count/raw tag bytes under defaults before constructing an
// immutable value. Input is borrowed during the call; canonical duplicates reject
// and failure returns zero Text rather than a partial row set.
func FromRows(rows []Row) (localized.Text, error) {
	limits := localized.DefaultLimits()
	if len(rows) > limits.MaxLocales {
		return localized.Text{}, localized.ErrLimitExceeded
	}
	entries := make([]localized.Entry, 0, len(rows))
	for _, row := range rows {
		if len(row.Locale) > limits.MaxTagBytes {
			return localized.Text{}, localized.ErrLimitExceeded
		}
		if strings.ContainsAny(row.Locale, "_ \t\r\n") {
			return localized.Text{}, localized.ErrInvalidLocale
		}
		tag, err := locale.Parse(row.Locale)
		if err != nil {
			return localized.Text{}, localized.ErrInvalidLocale
		}
		entries = append(entries, localized.Entry{Locale: tag, Text: row.Text})
	}
	return localized.NewText(entries...)
}
