package engine_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alkurbatov/guppy/internal/storage/engine"
)

func fillInMemoryEngine(t *testing.T) *engine.InMemory {
	t.Helper()

	db := engine.NewInMemory()
	err := db.Set("key", "value")
	require.NoError(t, err)

	return db
}

func TestInMemory_Set(t *testing.T) {
	sut := fillInMemoryEngine(t)

	result, err := sut.Get("key")

	require.NoError(t, err)
	require.Equal(t, "value", result)
}

func TestInMemory_Update(t *testing.T) {
	sut := fillInMemoryEngine(t)

	err := sut.Set("key", "new-value")
	require.NoError(t, err)

	result, err := sut.Get("key")
	require.NoError(t, err)
	require.Equal(t, "new-value", result)
}

func TestInMemory_SetBadData(t *testing.T) {
	tt := []struct {
		name     string
		key      string
		value    string
		expected error
	}{
		{
			name:     "Set empty key",
			key:      "",
			value:    "value",
			expected: engine.ErrEmptyKey,
		},
		{
			name:     "Set empty value",
			key:      "key",
			value:    "",
			expected: engine.ErrEmptyValue,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sut := engine.NewInMemory()

			err := sut.Set(tc.key, tc.value)

			require.ErrorIs(t, err, tc.expected)
		})
	}
}

func TestInMemory_GetUnknownKey(t *testing.T) {
	sut := engine.NewInMemory()

	_, err := sut.Get("xxx")

	require.ErrorIs(t, err, engine.ErrKeyNotFound)
}

func TestInMemory_Del(t *testing.T) {
	sut := fillInMemoryEngine(t)

	err := sut.Del("key")
	require.NoError(t, err)

	_, err = sut.Get("key")
	require.ErrorIs(t, err, engine.ErrKeyNotFound)
}

func TestInMemory_DelUnknownKey(t *testing.T) {
	sut := engine.NewInMemory()

	err := sut.Del("xxx")

	require.ErrorIs(t, err, engine.ErrKeyNotFound)
}
