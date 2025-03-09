package main

import (
	"bufio"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/alkurbatov/guppy/internal/compute/parser"
	"github.com/alkurbatov/guppy/internal/database"
	"github.com/alkurbatov/guppy/internal/storage/engine"
)

func run() int {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintln(os.Stderr, "create logger", err)
		return 1
	}
	defer logger.Sync() //nolint: errcheck //нет смысла обрабатывать ошибку, приложение завершается

	var queryParser parser.Parser
	inMemoryEngine := engine.NewInMemory()
	db := database.New(logger, queryParser, inMemoryEngine)

	fmt.Println("Enter query:")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Print(">>> ")

		var result string

		result, err = db.ProcessRequest(scanner.Text())
		if err != nil {
			fmt.Println("error: ", err)
			continue
		}

		fmt.Println(result)
	}

	if err = scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
