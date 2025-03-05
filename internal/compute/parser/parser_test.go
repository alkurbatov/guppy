package parser_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alkurbatov/guppy/internal/compute/parser"
)

func TestParseText(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected parser.Query
	}{
		{
			name:     "Set key",
			input:    "SET weather_2_pm cold_moscow_weather",
			expected: parser.NewQuery(parser.SET, "weather_2_pm", "cold_moscow_weather"),
		},
		{
			name:     "Get key",
			input:    "GET /etc/nginx/config",
			expected: parser.NewQuery(parser.GET, "/etc/nginx/config"),
		},
		{
			name:     "Capital letters key",
			input:    "GET MOSCOW",
			expected: parser.NewQuery(parser.GET, "MOSCOW"),
		},
		{
			name:     "Too many spaces",
			input:    "SET   abc    2",
			expected: parser.NewQuery(parser.SET, "abc", "2"),
		},
		{
			name:     "Del key",
			input:    "DEL user_****",
			expected: parser.NewQuery(parser.DEL, "user_****"),
		},
		{
			name:     "Short args",
			input:    "SET a b",
			expected: parser.NewQuery(parser.SET, "a", "b"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parser.ParseText(tc.input)

			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestParseText_BadInput(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "Empty input",
			input:    "",
			expected: parser.ErrBadInput,
		},
		{
			name:     "Only spaces",
			input:    "     ",
			expected: parser.ErrUnknownCommand,
		},
		{
			name:     "Malformed input",
			input:    "ksdsjkdjjdjs",
			expected: parser.ErrBadInput,
		},
		{
			name:     "Unknown command",
			input:    "SELECT 123",
			expected: parser.ErrUnknownCommand,
		},
		{
			name:     "Unicode symbols",
			input:    "GET ммм",
			expected: parser.ErrBadSymbol,
		},
		{
			name:     "Unexpected punctuation",
			input:    "GET a,",
			expected: parser.ErrBadSymbol,
		},
		{
			name:     "Lower case command",
			input:    "del user_****",
			expected: parser.ErrUnknownCommand,
		},
		{
			name:     "Mixed case command",
			input:    "DeL user_****",
			expected: parser.ErrUnknownCommand,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parser.ParseText(tc.input)

			require.ErrorIs(t, err, tc.expected)
		})
	}
}

func TestParseText_ArgsError(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "SET: No args",
			input:    "SET ",
			expected: parser.ErrNotEnoughArgs,
		},
		{
			name:     "SET: Not enough args",
			input:    "SET key",
			expected: parser.ErrNotEnoughArgs,
		},
		{
			name:     "SET: Too many args",
			input:    "SET key value1 value2",
			expected: parser.ErrTooManyArgs,
		},
		{
			name:     "GET: No args",
			input:    "GET ",
			expected: parser.ErrNotEnoughArgs,
		},
		{
			name:     "GET: Too many args",
			input:    "GET key value",
			expected: parser.ErrTooManyArgs,
		},
		{
			name:     "DEL: No args",
			input:    "DEL ",
			expected: parser.ErrNotEnoughArgs,
		},
		{
			name:     "DEL: Too many args",
			input:    "DEL key value",
			expected: parser.ErrTooManyArgs,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parser.ParseText(tc.input)

			require.ErrorIs(t, err, tc.expected)
		})
	}
}
