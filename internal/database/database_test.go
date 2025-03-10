package database_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/alkurbatov/guppy/internal/compute/parser"
	"github.com/alkurbatov/guppy/internal/database"
	"github.com/alkurbatov/guppy/internal/storage/engine"
)

func TestHandleRequest(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		setup    func(mockParser *database.MockParser, mockEngine *database.MockEngine)
		expected string
	}{
		{
			name:  "Set key to value",
			input: "SET weather_2_pm cold_moscow_weather",
			setup: func(mockParser *database.MockParser, mockEngine *database.MockEngine) {
				q := parser.NewQuery(parser.SET, "weather_2_pm", "cold_moscow_weather")
				mockParser.EXPECT().ParseText("SET weather_2_pm cold_moscow_weather").Return(q, nil)

				mockEngine.EXPECT().Set("weather_2_pm", "cold_moscow_weather").Return(nil)
			},
			expected: "OK",
		},
		{
			name:  "Get a key",
			input: "GET /etc/nginx/config",
			setup: func(mockParser *database.MockParser, mockEngine *database.MockEngine) {
				q := parser.NewQuery(parser.GET, "/etc/nginx/config")
				mockParser.EXPECT().ParseText("GET /etc/nginx/config").Return(q, nil)

				mockEngine.EXPECT().Get("/etc/nginx/config").Return("some-value", nil)
			},
			expected: "some-value",
		},
		{
			name:  "Delete a key",
			input: "DEL user_****",
			setup: func(mockParser *database.MockParser, mockEngine *database.MockEngine) {
				q := parser.NewQuery(parser.DEL, "user_****")
				mockParser.EXPECT().ParseText("DEL user_****").Return(q, nil)

				mockEngine.EXPECT().Del("user_****")
			},
			expected: "OK",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockParser := database.NewMockParser(ctrl)
			mockEngine := database.NewMockEngine(ctrl)
			tc.setup(mockParser, mockEngine)
			sut := database.New(zap.NewNop(), mockParser, mockEngine)

			result, err := sut.ProcessRequest(tc.input)

			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestHandleRequest_DependenciesFailure(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		setup    func(mockParser *database.MockParser, mockEngine *database.MockEngine)
		expected error
	}{
		{
			name:  "Bad query",
			input: "SELECT FROM *",
			setup: func(mockParser *database.MockParser, _ *database.MockEngine) {
				var q parser.Query
				mockParser.EXPECT().ParseText("SELECT FROM *").Return(q, parser.ErrUnknownCommand)
			},
			expected: parser.ErrUnknownCommand,
		},
		{
			name:  "Unexpected command from parser",
			input: "GET /etc/nginx/config",
			setup: func(mockParser *database.MockParser, _ *database.MockEngine) {
				q := parser.NewQuery("XXX")
				mockParser.EXPECT().ParseText("GET /etc/nginx/config").Return(q, nil)
			},
			expected: parser.ErrUnknownCommand,
		},
		{
			name:  "Set key to empty value",
			input: "SET weather_2_pm cold_moscow_weather",
			setup: func(mockParser *database.MockParser, mockEngine *database.MockEngine) {
				q := parser.NewQuery(parser.SET, "weather_2_pm", "cold_moscow_weather")
				mockParser.EXPECT().ParseText("SET weather_2_pm cold_moscow_weather").Return(q, nil)

				mockEngine.EXPECT().Set("weather_2_pm", "cold_moscow_weather").Return(engine.ErrEmptyKey)
			},
			expected: engine.ErrEmptyKey,
		},
		{
			name:  "Get unexisting key",
			input: "GET /etc/nginx/config",
			setup: func(mockParser *database.MockParser, mockEngine *database.MockEngine) {
				q := parser.NewQuery(parser.GET, "/etc/nginx/config")
				mockParser.EXPECT().ParseText("GET /etc/nginx/config").Return(q, nil)

				mockEngine.EXPECT().Get("/etc/nginx/config").Return("", engine.ErrKeyNotFound)
			},
			expected: engine.ErrKeyNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockParser := database.NewMockParser(ctrl)
			mockEngine := database.NewMockEngine(ctrl)
			tc.setup(mockParser, mockEngine)
			sut := database.New(zap.NewNop(), mockParser, mockEngine)

			_, err := sut.ProcessRequest(tc.input)

			require.ErrorIs(t, err, tc.expected)
		})
	}
}
