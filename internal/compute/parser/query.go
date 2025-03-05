package parser

// Command команда базы данных.
type Command string

const (
	SET Command = "SET"
	GET Command = "GET"
	DEL Command = "DEL"
)

// Query запрос к базе данных.
type Query struct {
	// Command идентификатор выполняемой команды.
	Cmd Command

	// Args аргументы запроса.
	Args []string
}

// NewQuery создает новый объект [Query].
func NewQuery(cmd Command, args ...string) Query {
	return Query{
		Cmd:  cmd,
		Args: args,
	}
}
