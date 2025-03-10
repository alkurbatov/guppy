package parser

// Command команда базы данных.
type Command string

const (
	SET Command = "SET"
	GET Command = "GET"
	DEL Command = "DEL"
)

var commandArgsCount = map[Command]int{
	SET: 2,
	GET: 1,
	DEL: 1,
}

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
