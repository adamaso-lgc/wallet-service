package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUUIDFromString_Valid(t *testing.T) {
	const valid = "550e8400-e29b-41d4-a716-446655440000"

	u, err := UUIDFromString(valid)

	require.NoError(t, err)
	assert.True(t, u.Valid)
}

func TestUUIDFromString_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"not a uuid", "not-a-uuid"},
		{"too short", "550e8400-e29b-41d4"},
		{"wrong chars", "gggggggg-gggg-gggg-gggg-gggggggggggg"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := UUIDFromString(tc.input)
			require.Error(t, err)
		})
	}
}

func TestNumericFromFloat64(t *testing.T) {
	tests := []struct {
		name  string
		input float64
	}{
		{"zero", 0},
		{"positive", 100.50},
		{"negative", -25.75},
		{"large", 9_999_999.99999999},
		{"small fraction", 0.00000001},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n, err := NumericFromFloat64(tc.input)

			require.NoError(t, err)
			assert.True(t, n.Valid)
		})
	}
}
