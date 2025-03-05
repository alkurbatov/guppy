package database

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/alkurbatov/guppy/internal/compute/parser"
)

//go:generate go tool mockgen -source database.go -destination engine_mock.go -package database Engine
type Engine interface {
	// Set добавляет ключ в хранилище или модифицирует значение существующего.
	Set(key, value string) error

	// Get возвращает значение ключа из хранилища.
	Get(key string) (string, error)

	// Del удаляет ключ из хранилища.
	Del(key string) error
}

// Database key-value база данных.
type Database struct {
	logger *zap.Logger
	engine Engine
}

// New создает новый объект [Database].
func New(logger *zap.Logger, engine Engine) Database {
	return Database{
		logger: logger,
		engine: engine,
	}
}

// ProcessRequest обрабатывает запрос.
func (db Database) ProcessRequest(input string) (string, error) {
	query, err := parser.ParseText(input)
	if err != nil {
		db.logger.Error("cannot parse request",
			zap.Error(err), zap.String("input", input))

		return "", fmt.Errorf("parse input: %w", err)
	}

	switch query.Cmd {
	case parser.SET:
		key := query.Args[0]
		value := query.Args[1]

		if err = db.engine.Set(key, value); err != nil {
			db.logger.Error("command SET has failed",
				zap.Error(err), zap.String("key", key), zap.String("value", value))

			return "", fmt.Errorf("set key value: %w", err)
		}

		db.logger.Info("set value of a key", zap.String("key", key), zap.String("value", value))

		return "OK", nil

	case parser.GET:
		key := query.Args[0]
		var value string

		value, err = db.engine.Get(key)
		if err != nil {
			db.logger.Error("command GET has failed",
				zap.Error(err), zap.String("key", key))

			return "", fmt.Errorf("get key: %w", err)
		}

		db.logger.Info("got value of a key", zap.String("key", key), zap.String("value", value))

		return value, nil

	case parser.DEL:
		key := query.Args[0]

		err = db.engine.Del(key)
		if err != nil {
			db.logger.Error("command DEL has failed",
				zap.Error(err), zap.String("key", key))

			return "", fmt.Errorf("delete key: %w", err)
		}

		db.logger.Info("removed a key", zap.String("key", key))

		return "OK", nil

	default:
		return "", parser.ErrUnknownCommand
	}
}
