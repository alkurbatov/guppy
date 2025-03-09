package database

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/alkurbatov/guppy/internal/compute/parser"
)

//go:generate go tool mockgen -source database.go -destination database_mock.go -package database

type Parser interface {
	ParseText(input string) (parser.Query, error)
}

type Engine interface {
	// Set добавляет ключ в хранилище или модифицирует значение существующего.
	Set(key, value string) error

	// Get возвращает значение ключа из хранилища.
	Get(key string) (string, error)

	// Del удаляет ключ из хранилища.
	Del(key string)
}

// Database key-value база данных.
type Database struct {
	logger *zap.Logger
	parser Parser
	engine Engine
}

// New создает новый объект [Database].
func New(logger *zap.Logger, parser Parser, engine Engine) Database {
	return Database{
		logger: logger,
		parser: parser,
		engine: engine,
	}
}

// ProcessRequest обрабатывает запрос.
func (db Database) ProcessRequest(input string) (string, error) {
	query, err := db.parser.ParseText(input)
	if err != nil {
		db.logger.Error("cannot parse request",
			zap.Error(err), zap.String("input", input))

		return "", fmt.Errorf("parse input: %w", err)
	}

	switch query.Cmd {
	case parser.SET:
		return db.set(query.Args[0], query.Args[1])

	case parser.GET:
		return db.get(query.Args[0])

	case parser.DEL:
		return db.del(query.Args[0])

	default:
		return "", parser.ErrUnknownCommand
	}
}

func (db Database) set(key, value string) (string, error) {
	if err := db.engine.Set(key, value); err != nil {
		db.logger.Error("command SET has failed",
			zap.Error(err), zap.String("key", key), zap.String("value", value))

		return "", fmt.Errorf("set key value: %w", err)
	}

	db.logger.Info("set value of a key", zap.String("key", key), zap.String("value", value))

	return "OK", nil
}

func (db Database) get(key string) (string, error) {
	value, err := db.engine.Get(key)
	if err != nil {
		db.logger.Error("command GET has failed",
			zap.Error(err), zap.String("key", key))

		return "", fmt.Errorf("get key: %w", err)
	}

	db.logger.Info("got value of a key", zap.String("key", key), zap.String("value", value))

	return value, nil
}

func (db Database) del(key string) (string, error) {
	db.engine.Del(key)
	db.logger.Info("removed a key", zap.String("key", key))

	return "OK", nil
}
