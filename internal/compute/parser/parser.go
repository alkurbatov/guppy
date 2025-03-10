package parser

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmptyInput     = errors.New("empty input")
	ErrUnknownCommand = errors.New("unknown command")
	ErrBadSymbol      = errors.New("input contains invalid symbols")
	ErrNotEnoughArgs  = errors.New("not enough arguments")
	ErrTooManyArgs    = errors.New("too many arguments")
)

var isValid = regexp.MustCompile(`^[a-zA-Z0-9*/_]+$`).MatchString

type Parser struct{}

func (p Parser) ParseText(input string) (Query, error) {
	tokens := strings.Fields(input)

	if len(tokens) == 0 {
		return Query{}, ErrEmptyInput
	}

	cmd := Command(tokens[0])
	args := tokens[1:]
	if err := validate(cmd, args); err != nil {
		return Query{}, err
	}

	return Query{
		Cmd:  cmd,
		Args: args,
	}, nil
}

func validate(cmd Command, args []string) error {
	count, ok := commandArgsCount[cmd]
	if !ok {
		return ErrUnknownCommand
	}

	if len(args) < count {
		return ErrNotEnoughArgs
	}

	if len(args) > count {
		return ErrTooManyArgs
	}

	for _, arg := range args {
		if !isValid(arg) {
			return ErrBadSymbol
		}
	}

	return nil
}
