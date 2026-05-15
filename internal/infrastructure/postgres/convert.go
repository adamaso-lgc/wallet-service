package postgres

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UUIDToString formats a pgtype.UUID back to the standard hyphenated string.
func UUIDToString(u pgtype.UUID) string {
	return uuid.UUID(u.Bytes).String()
}

// NumericToFloat64 converts a pgtype.Numeric to float64 via its string
func NumericToFloat64(n pgtype.Numeric) (float64, error) {
	v, err := n.Value()
	if err != nil {
		return 0, fmt.Errorf("numeric to float64: %w", err)
	}
	s, ok := v.(string)
	if !ok {
		return 0, fmt.Errorf("unexpected numeric value type: %T", v)
	}
	return strconv.ParseFloat(s, 64)
}

func UUIDFromString(s string) (pgtype.UUID, error) {
	var u pgtype.UUID
	if err := u.Scan(s); err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid uuid %q: %w", s, err)
	}
	return u, nil
}

// NumericFromFloat64 converts a float64 to pgtype.Numeric.
func NumericFromFloat64(f float64) (pgtype.Numeric, error) {
	var n pgtype.Numeric
	s := strconv.FormatFloat(f, 'f', -1, 64)
	if err := n.Scan(s); err != nil {
		return pgtype.Numeric{}, fmt.Errorf("convert float64 to numeric: %w", err)
	}
	return n, nil
}
