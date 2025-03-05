package engine

import "errors"

var (
	ErrEmptyKey    = errors.New("key is empty")
	ErrEmptyValue  = errors.New("value is empty")
	ErrKeyNotFound = errors.New("key not found")
)

// InMemory реализация key-value хранилища в оперативной памяти.
type InMemory struct {
	data map[string]string
}

// NewInMemory создает новый объект [InMemory].
func NewInMemory() *InMemory {
	return &InMemory{
		data: make(map[string]string, 0),
	}
}

// Set добавляет ключ в хранилище или модифицирует значение существующего.
func (db *InMemory) Set(key, value string) error {
	if key == "" {
		return ErrEmptyKey
	}
	if value == "" {
		return ErrEmptyValue
	}

	db.data[key] = value

	return nil
}

// Get возвращает значение ключа из хранилища.
func (db *InMemory) Get(key string) (string, error) {
	value, ok := db.data[key]
	if !ok {
		return "", ErrKeyNotFound
	}

	return value, nil
}

// Del удаляет ключ из хранилища.
func (db *InMemory) Del(key string) error {
	if _, err := db.Get(key); err != nil {
		return err
	}

	delete(db.data, key)

	return nil
}
