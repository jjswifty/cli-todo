// Package for TODOs
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
)

const pathToTodos = "todos.json"

func main() {

	rawArgs := os.Args[1:]

	var mainError error

	cmd, err := parseArgs(rawArgs)

	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	switch c := cmd.(type) {
	case addCommand:
		mainError = runAdd(pathToTodos, c.text)

	case listAllCommand:
		mainError = runList(pathToTodos)
	case listOneCommand:
		mainError = runListByID(pathToTodos, c.id)
	case doneCommand:
		mainError = runDone(pathToTodos, c.id)
	case undoneCommand:
		mainError = runUndone(pathToTodos, c.id)
	case deleteCommand:
		mainError = runDelete(pathToTodos, c.id)
	default:
		mainError = fmt.Errorf("unknown command")
	}

	if mainError != nil {
		if errors.Is(mainError, errNoTodos) {
			color.Yellow("У вас еще нет задач. Добавьте с помощью команды `add`.")
		} else {
			color.Red(mainError.Error())
		}

		os.Exit(1)
	}

	os.Exit(0)
}
