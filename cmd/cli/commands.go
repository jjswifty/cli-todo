package main

import (
	"errors"
	"fmt"
	"slices"

	"github.com/fatih/color"
)

var errNoTodos = errors.New("список задач пуст")

func requireTodos(path string) (todoList, error) {
	list, err := loadTodos(path)

	if err != nil {
		return todoList{}, err
	}

	if len(list.Todos) == 0 {
		return todoList{}, errNoTodos // гвард: провал — это значение
	}

	return list, nil
}

func runAdd(path string, text string) error {
	loadedTodos, err := loadTodos(path)

	if err != nil {
		return err
	}

	createdTodo := newTodoFromText(text, loadedTodos.NextID)

	loadedTodos.Todos = append(loadedTodos.Todos, createdTodo)
	loadedTodos.NextID++

	err = saveTodos(loadedTodos, path)

	if err != nil {
		return err
	}

	color.Green("Задача \"%s\" успешно добавлена.", createdTodo.Text)

	return nil
}

func runList(path string) error {
	loadedTodos, err := requireTodos(path)

	if err != nil {
		return err
	}

	fmt.Println("Список задач:")

	for _, t := range loadedTodos.Todos {
		printTodo(t)
	}

	return nil
}

func runListByID(path string, id int) error {
	loadedTodos, err := requireTodos(path)

	if err != nil {
		return err
	}

	foundTodo, err := getTodoByIDFromList(loadedTodos.Todos, id)

	if err != nil {
		return err
	}

	printTodo(*foundTodo)

	return nil
}

func runDone(path string, id int) error {
	err := toggleTodoCompletionStatus(path, id, true)

	if err != nil {
		return err
	}

	return nil
}

func runUndone(path string, id int) error {
	err := toggleTodoCompletionStatus(path, id, false)

	if err != nil {
		return err
	}

	return nil
}

func runDelete(path string, id int) error {
	err := deleteTodoByID(path, id)

	if err != nil {
		return err
	}

	return nil
}

func toggleTodoCompletionStatus(path string, id int, status bool) error {
	loadedTodos, err := requireTodos(path)

	if err != nil {
		return err
	}

	foundTodo, err := getTodoByIDFromList(loadedTodos.Todos, id)

	if err != nil {
		return err
	}

	foundTodo.Completed = status

	err = saveTodos(loadedTodos, path)

	if err != nil {
		return err
	}

	printTodo(*foundTodo)

	return nil
}

func deleteTodoByID(path string, id int) error {
	loadedTodos, err := requireTodos(path)

	if err != nil {
		return err
	}

	before := len(loadedTodos.Todos)

	var found todo

	loadedTodos.Todos = slices.DeleteFunc(loadedTodos.Todos, func(t todo) bool {
		if t.ID == id {
			found = t
			return true
		}
		return false
	})

	if len(loadedTodos.Todos) == before {
		return fmt.Errorf("todo with id %d not found", id)
	}

	err = saveTodos(loadedTodos, path)

	if err != nil {
		return err
	}

	fmt.Println("Успешно удалили следующую задачу:")

	printTodo(found)

	return nil
}

// приходит указатель на лист
// мы возвращаем тудушку из оригинального листа, и при ее изменении будет меняться и
// оригинальный лист
// list [1, 2, 3]
// getTodoByIdFromList(2)
// меняем его параметры
// list [1, 2 [CHANGED], 3]
func getTodoByIDFromList(list []todo, id int) (*todo, error) {
	for i := range list {
		if list[i].ID == id {
			return &list[i], nil
		}
	}

	return nil, fmt.Errorf("todo с id %d не найдено", id)
}

func printTodo(t todo) {
	faint := color.New(color.Faint).SprintFunc()
	done := color.New(color.FgGreen).SprintFunc()

	created := t.CreatedAt.Format("02.01.2006 15:04")

	mark := " "
	if t.Completed {
		mark = done("x")
	}

	fmt.Printf("[%s] %s %s %s\n", mark, faint(fmt.Sprintf("#%d", t.ID)), faint(t.Text), faint("("+created+")"))
}
