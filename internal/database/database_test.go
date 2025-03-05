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
		setup    func(m *database.MockEngine)
		expected string
	}{
		{
			name:  "Set key to value",
			input: "SET weather_2_pm cold_moscow_weather",
			setup: func(m *database.MockEngine) {
				m.EXPECT().Set("weather_2_pm", "cold_moscow_weather").Return(nil)
			},
			expected: "OK",
		},
		{
			name:  "Get a key",
			input: "GET /etc/nginx/config",
			setup: func(m *database.MockEngine) {
				m.EXPECT().Get("/etc/nginx/config").Return("some-value", nil)
			},
			expected: "some-value",
		},
		{
			name:  "Delete a key",
			input: "DEL user_****",
			setup: func(m *database.MockEngine) {
				m.EXPECT().Del("user_****").Return(nil)
			},
			expected: "OK",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := database.NewMockEngine(ctrl)
			tc.setup(m)
			sut := database.New(zap.NewNop(), m)

			result, err := sut.ProcessRequest(tc.input)

			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestHandleRequest_EngineFailure(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		setup    func(m *database.MockEngine)
		expected error
	}{
		{
			name:  "Set key to value",
			input: "SET weather_2_pm cold_moscow_weather",
			setup: func(m *database.MockEngine) {
				m.EXPECT().Set("weather_2_pm", "cold_moscow_weather").Return(engine.ErrEmptyKey)
			},
			expected: engine.ErrEmptyKey,
		},
		{
			name:  "Get a key",
			input: "GET /etc/nginx/config",
			setup: func(m *database.MockEngine) {
				m.EXPECT().Get("/etc/nginx/config").Return("", engine.ErrKeyNotFound)
			},
			expected: engine.ErrKeyNotFound,
		},
		{
			name:  "Delete a key",
			input: "DEL user_****",
			setup: func(m *database.MockEngine) {
				m.EXPECT().Del("user_****").Return(engine.ErrKeyNotFound)
			},
			expected: engine.ErrKeyNotFound,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			m := database.NewMockEngine(ctrl)
			tc.setup(m)
			sut := database.New(zap.NewNop(), m)

			_, err := sut.ProcessRequest(tc.input)

			require.ErrorIs(t, err, tc.expected)
		})
	}
}

func TestHandleRequest_BadQuery(t *testing.T) {
	sut := database.New(zap.NewNop(), nil)

	_, err := sut.ProcessRequest("SELECT FROM *")

	require.ErrorIs(t, err, parser.ErrUnknownCommand)
}
