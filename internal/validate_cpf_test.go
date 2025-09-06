package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCpf(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected bool
	}{
		"valid cpf 1": {
			input:    "97456321558",
			expected: true,
		},
		"valid cpf 2": {
			input:    "974.563.215-58",
			expected: true,
		},
		"valid cpf 3": {
			input:    "71428793860",
			expected: true,
		},
		"valid cpf 4": {
			input:    "87748248800",
			expected: true,
		},
		"invalid cpf 1": {
			input:    "",
			expected: false,
		},
		"invalid cpf 2": {
			input:    "11111111111",
			expected: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, validateCpf(tt.input))
		})
	}
}
