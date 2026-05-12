package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

// uuidFromString parses a UUID string into pgtype.UUID.
func uuidFromString(s string) (pgtype.UUID, error) {
	var u pgtype.UUID
	if err := u.Scan(s); err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid uuid %q: %w", s, err)
	}
	return u, nil
}

// numericFromFloat64 converts a float64 to pgtype.Numeric.
func numericFromFloat64(f float64) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	if err := n.Scan(f); err != nil {
		return pgtype.Numeric{}, fmt.Errorf("convert float64 to numeric: %w", err)
	}
	return n, nil
}
