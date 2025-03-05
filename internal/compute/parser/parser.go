package parser

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrBadInput       = errors.New("not a DB query")
	ErrUnknownCommand = errors.New("unknown command")
	ErrBadSymbol      = errors.New("input contains invalid symbols")
	ErrNotEnoughArgs  = errors.New("not enough arguments")
	ErrTooManyArgs    = errors.New("too many arguments")
)

var isValid = regexp.MustCompile(`^[a-zA-Z0-9*/_]+$`).MatchString

func ParseText(input string) (Query, error) {
	tokens := strings.Split(input, " ")

	if len(tokens) < 2 {
		return Query{}, ErrBadInput
	}

	cmd := Command(tokens[0])
	if err := validateCommand(cmd); err != nil {
		return Query{}, err
	}

	args := stripSpaces(tokens[1:])
	if err := validateArgs(cmd, args); err != nil {
		return Query{}, err
	}

	return Query{
		Cmd:  cmd,
		Args: args,
	}, nil
}

func stripSpaces(src []string) []string {
	rv := make([]string, 0, len(src))

	for i := range src {
		s := strings.TrimSpace(src[i])
		if s == "" {
			continue
		}

		rv = append(rv, s)
	}

	return rv
}

func validateCommand(cmd Command) error {
	switch cmd {
	case SET, GET, DEL:
		return nil
	default:
		return ErrUnknownCommand
	}
}

func validateArgs(cmd Command, args []string) error {
	count := 1
	if cmd == SET {
		count = 2
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
