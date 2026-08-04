package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// sealed interface with marker method (ts analogue of discriminated union)
type command interface {
	isCommand()
}

type addCommand struct{ text string }
type listAllCommand struct{}
type listOneCommand struct{ id int }
type doneCommand struct{ id int }
type undoneCommand struct{ id int }
type deleteCommand struct{ id int }

func (addCommand) isCommand()     {}
func (listAllCommand) isCommand() {}
func (listOneCommand) isCommand() {}
func (doneCommand) isCommand()    {}
func (undoneCommand) isCommand()  {}
func (deleteCommand) isCommand()  {}

// Получает на вход rawArgs без 1 системного аргумента
func parseArgs(args []string) (command, error) {
	allowedCommands := [5]string{"add", "list", "done", "undone", "delete"}

	if len(args) == 0 {
		return nil, fmt.Errorf("no command was specified")
	}

	// allowedCommand
	switch args[0] {
	case "add":
		todoText := strings.TrimSpace(strings.Join(args[1:], " "))

		if todoText == "" {
			return nil, errors.New("todo text cannot be empty")
		}

		return addCommand{
			text: todoText,
		}, nil
	case "list":
		if len(args) == 1 {
			return listAllCommand{}, nil
		}

		id, err := parseIDArg(args)

		if err != nil {
			return nil, fmt.Errorf("list: %w", err)
		}

		return listOneCommand{
			id,
		}, nil
	case "done":
		id, err := parseIDArg(args)
		if err != nil {
			return nil, fmt.Errorf("done: %w", err)
		}

		return doneCommand{
			id,
		}, nil
	case "undone":
		id, err := parseIDArg(args)
		if err != nil {
			return nil, fmt.Errorf("undone: %w", err)
		}

		return undoneCommand{
			id,
		}, nil
	case "delete":
		id, err := parseIDArg(args)
		if err != nil {
			return nil, fmt.Errorf("delete: %w", err)
		}

		return deleteCommand{
			id,
		}, nil
	default:
		return nil, fmt.Errorf("unknown command, allowed commands: %s", allowedCommands)
	}
}

func parseIDArg(args []string) (int, error) {
	if len(args) == 0 {
		return 0, fmt.Errorf("no command was provided, args: %+v", args)
	}

	if len(args) != 2 {
		return 0, fmt.Errorf("expected exactly 1 integer argument for id, but got %d", len(args)-1)
	}

	if strings.HasPrefix(args[1], "+") || strings.HasPrefix(args[1], "-") {
		return 0, fmt.Errorf("id must not contain a sign")
	}

	id, err := strconv.Atoi(args[1])

	if err != nil {
		return 0, fmt.Errorf("id is not correct %q: \n%w", args[1], err)
	}

	return id, nil
}
