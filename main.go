package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type todoDTO struct {
	CreatedAt *string `json:"createdAt"`
	Completed *bool   `json:"completed"`
	Text      *string `json:"text"`
}

type todosFileDTO struct {
	Todos []todoDTO `json:"todos"`
}

type todosFile struct {
	Todos []todo `json:"todos"`
}

type todo struct {
	CreatedAt string `json:"createdAt"`
	Completed bool   `json:"completed"`
	Text      string `json:"text"`
}

func (data todoDTO) validate() error {
	if data.Text == nil {
		return errors.New("missing field `text`")
	}

	if data.CreatedAt == nil {
		return errors.New("missing field `createdAt`")
	}

	if data.Completed == nil {
		return errors.New("missing field `completed`")
	}

	return nil
}

func (data todoDTO) toDomain() todo {
	return todo{
		// это хоть и указатели, но они указаывают на исходную строку. копировать строку нет смысла, она immutable
		CreatedAt: *data.CreatedAt,
		Text:      *data.Text,
		Completed: *data.Completed,
	}
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

		fmt.Println(string(fileContent), "filecontent")

		var rawData todosFileDTO

		// создаем go struct из json:

		err = json.Unmarshal(fileContent, &rawData)

		if err != nil {
			fmt.Printf("Не удалось конвертировать JSON в структуру. Возможно, JSON поврежден. Err: %v", err)

			os.Exit(1)
		}

		fmt.Println("After unmarashal", rawData)

		var currentTodos []todo

		for i, dto := range rawData.Todos {
			if err := dto.validate(); err != nil {

				fmt.Printf("Error in field todos[%d]: %v. Raw value: %+v", i, err, dto)
				os.Exit(1)
			}

			currentTodos = append(currentTodos, dto.toDomain())
		}

		currentTodos = append(currentTodos, createdTodo)

		fmt.Printf("%+v", currentTodos)

		constructedJson, err := json.Marshal(todosFile{Todos: currentTodos})

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
		CreatedAt: time.Now().String(),
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
