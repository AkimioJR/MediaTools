package wordmatch_test

import (
	"MediaTools/internal/pkg/wordmatch"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseOffsetExpr(t *testing.T) {
	tests := []struct {
		expr     string
		episode  int
		expected int
	}{
		{"EP+1", 5, 6},
		{"EP-2", 10, 8},
		{"2*EP", 3, 6},
		{"2*EP-4", 5, 6},
		{"2*EP+3", 4, 11},
		{"2*EP", 5, 10},
		{"EP", 7, 7},
		{"-1*EP", 5, -5},
		{"2*EP+3", 5, 13},
		{"EP+6", 5, 11},
		{"2*EP-2", 5, 8},
	}

	for _, test := range tests {
		result, err := wordmatch.ParseOffsetExpr(test.expr, test.episode)
		require.NoError(t, err, fmt.Sprintf("Failed to parse expression: %s with episode: %d", test.expr, test.episode))
		require.Equal(t, test.expected, result, fmt.Sprintf("Expected %d but got %d for expression: %s with episode: %d", test.expected, result, test.expr, test.episode))
	}
}
