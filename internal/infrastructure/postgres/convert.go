package postgres

import (
	"fmt"
	"strconv"

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
// pgtype.Numeric.Scan does not accept float64 directly; we marshal to the
// string representation first (e.g. "100.5") so pgtype can parse it exactly.
func numericFromFloat64(f float64) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	s := strconv.FormatFloat(f, 'f', -1, 64)
	if err := n.Scan(s); err != nil {
		return pgtype.Numeric{}, fmt.Errorf("convert float64 to numeric: %w", err)
	}
	return n, nil
}
