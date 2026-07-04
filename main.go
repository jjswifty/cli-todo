package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type todo struct {
	CreatedAt time.Time `json:"created_at"`
	Completed bool      `json:"completed"`
	Text      string    `json:"text"`
}

type todosFile struct {
	Todos []todo `json:"todos"`
}

func main() {
	allowedCommands := [1]string{"add"}

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("No command was specified; terminating program.")
		os.Exit(0)
	}

	command := args[0]

	if slices.Contains(allowedCommands[:], command) == false {
		fmt.Println("Unknown command. Allowed commands:", allowedCommands)
		os.Exit(0)
	}

	if len(args) == 1 {
		fmt.Printf("No argument was provided for command: \"%s\", terminating...", args[0])
		os.Exit(0)
	}

	argument := strings.Join(args[1:], " ")

	switch command {
	case "add":

		createdTodo := createTodoFromText(argument)

		/**
		Сначала мы должны проверить существует ли мас
		получаем джсон из файла -> превр. в обычный объект ->
		1 (?). файл оказался пустым, создаем объект сами
		2. добавляем в этот объект нашу туду -> превращаем в джсон -> записываем файл заново
		*/

		_, err := os.Stat("todos.json")

		if err != nil {
			// файла не существует, сразу создадим жсон и закинем внутрь
			createFileAndTodos(createdTodo)

			os.Exit(1)
		}

		/**
		Сценарий 2; файл существует, но теперь нужно проверить его валидность
		Если он валиден, то просто получим жсон, проверим наличие поля todos, потом если он массив
		то заапендим
		*/

		fileContent, err := os.ReadFile("todos.json")

		if err != nil {
			fmt.Printf("Ошибка чтения файла. Создадим свой. Err: %s", err)

			err = os.Remove("todos.json")

			if err != nil {
				fmt.Printf("Не удалось пересоздать существующий файл, terminating... Err: %s", err)
				os.Exit(1)
			}

			createFileAndTodos(createdTodo)

			os.Exit(1)
		}

		var data todosFile
		// создаем go struct из json:

		err = json.Unmarshal(fileContent, &data)

		if err != nil {
			fmt.Printf("Не удалось конвертировать JSON в структуру. Возможно, JSON поврежден. Err: %v", err)

			os.Exit(1)
		}

		data.Todos = append(data.Todos, createdTodo)

		fmt.Printf("%+v", data)

		constructedJson, err := json.Marshal(data)

		if err != nil {
			fmt.Printf("Ошибка создания JSON: %v", err)
			os.Exit(0)
		}

		fmt.Print(string(constructedJson))

		if err := os.WriteFile("todos.json", constructedJson, 0777); err != nil {
			fmt.Print("Error while file creating, probably security:", err)
			os.Exit(0)
		}

	default:
		panic("NO SUCH COMMAND WE GONNA DIIIE, lol")
	}
}

func createTodoFromText(text string) todo {
	return todo{
		Text:      text,
		Completed: false,
		CreatedAt: time.Now(),
	}
}

func createFileAndTodos(todoObj todo) {
	var todos []todo

	todosObject := todosFile{
		Todos: append(todos, todoObj),
	}

	constructedJson, err := json.Marshal(todosObject)

	if err != nil {
		fmt.Print("Ошибка создания JSON:", err)
		os.Exit(0)
	}

	if err := os.WriteFile("todos.json", constructedJson, 0777); err != nil {
		fmt.Print("Error while file creating, probably security:", err)
		os.Exit(0)
	}

	fmt.Printf("Задача \"%s\" успешно добавлена.", todoObj.Text)
}
