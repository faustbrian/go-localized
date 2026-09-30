package postgres_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	localized "github.com/faustbrian/go-localized"
	"github.com/faustbrian/go-localized/postgres"
)

func TestSecurityPersistenceRejectsOversizedRowsBeforeParsing(t *testing.T) {
	rows := make([]postgres.Row, 129)
	for i := range rows {
		rows[i].Locale = "invalid_locale"
	}
	t.Run("count", func(t *testing.T) {
		if _, err := postgres.FromRows(rows); !errors.Is(err, localized.ErrLimitExceeded) {
			t.Fatalf("rows error = %v", err)
		}
	})
	t.Run("tag", func(t *testing.T) {
		if _, err := postgres.FromRows([]postgres.Row{{Locale: strings.Repeat("x", 256)}}); !errors.Is(err, localized.ErrLimitExceeded) {
			t.Fatalf("tag bytes error = %v", err)
		}
	})
}

func TestSecurityPersistenceAcceptsExactRowCount(t *testing.T) {
	rows := make([]postgres.Row, localized.DefaultLimits().MaxLocales)
	for i := range rows {
		rows[i] = postgres.Row{Locale: fmt.Sprintf("x-entry-%d", i), Text: "text"}
	}
	value, err := postgres.FromRows(rows)
	if err != nil || value.Len() != len(rows) {
		t.Fatalf("exact row count = %d, %v", value.Len(), err)
	}
}

func TestSecurityPersistenceAcceptsMaximumTagBytes(t *testing.T) {
	// A valid BCP 47 private-use tag at the documented 255-byte ceiling.
	tag := "en-x-" + strings.Repeat("abcdefgh-", 27) + "abcdefg"
	value, err := postgres.FromRows([]postgres.Row{{Locale: tag, Text: "text"}})
	if err != nil || value.Len() != 1 || value.Entries()[0].Locale.String() != tag {
		t.Fatalf("maximum row tag = %v, %v", value.Entries(), err)
	}
}

func TestSecurityPGXNilDestinationReturnsError(t *testing.T) {
	var target *localized.Text
	if err := postgres.JSONBCodec().Unmarshal([]byte(`{}`), target); !errors.Is(err, postgres.ErrUnsupportedDatabaseType) {
		t.Fatalf("nil target error = %v", err)
	}
}
