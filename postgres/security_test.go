package postgres_test

import (
	"errors"
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

func TestSecurityPGXNilDestinationReturnsError(t *testing.T) {
	var target *localized.Text
	if err := postgres.JSONBCodec().Unmarshal([]byte(`{}`), target); !errors.Is(err, postgres.ErrUnsupportedDatabaseType) {
		t.Fatalf("nil target error = %v", err)
	}
}
