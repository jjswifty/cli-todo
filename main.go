// Package for TODOs
package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/fatih/color"
)

func main() {
	allowedCommands := [1]string{"add"}

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("No command was specified; terminating program.")
		os.Exit(1)
	}

	command := args[0]

	if !slices.Contains(allowedCommands[:], command) {
		fmt.Println("Unknown command. Allowed commands:", allowedCommands)
		os.Exit(1)
	}

	if len(args) == 1 {
		fmt.Printf("No argument was provided for command: \"%s\", terminating... \n", args[0])
		os.Exit(1)
	}

	argument := strings.Join(args[1:], " ")

	switch command {
	case "add":
		loadedTodos, err := loadTodos()

		if err != nil {
			fmt.Println(err)

			os.Exit(1)
		}

		createdTodo := newTodoFromText(argument, loadedTodos.NextID)

		loadedTodos.Todos = append(loadedTodos.Todos, createdTodo)
		loadedTodos.NextID++

		err = saveTodos(loadedTodos)

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		color.Green("Задача \"%s\" успешно добавлена.", createdTodo.Text)

		os.Exit(0)
	default:
		panic("NO SUCH COMMAND WE GONNA DIIIE, lol")
	}
}
