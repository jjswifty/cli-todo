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

type todosFileDTO struct {
	NextId *int      `json:"nextId"`
	Todos  []todoDTO `json:"todos"`
}
type todoDTO struct {
	Id        *int    `json:"id"`
	CreatedAt *string `json:"createdAt"`
	Completed *bool   `json:"completed"`
	Text      *string `json:"text"`
}

type todosFile struct {
	NextId int    `json:"nextId"`
	Todos  []todo `json:"todos"`
}

type todo struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"createdAt"`
	Completed bool   `json:"completed"`
	Text      string `json:"text"`
}

func (data todosFileDTO) validate() error {
	if data.NextId == nil {
		return errors.New("missing field `nextId`")
	}

	if *data.NextId < 0 {
		return errors.New("field nextId must be positive or 0")
	}

	// todos может и не быть, нельзя сразу распаковывать
	if data.Todos == nil {
		return errors.New("missing field `todos`")
	}

	var ids []int

	maxId := 0

	for i, dto := range data.Todos {
		if err := dto.validate(); err != nil {
			return fmt.Errorf("todo at index %d is invalid: %w", i, err)
		}

		id := *dto.Id

		if id > maxId {
			maxId = *dto.Id
		}

		if slices.Contains(ids, id) {
			return fmt.Errorf("todo at index %d is invalid: duplicate ID", i)
		}

		ids = append(ids, id)
	}

	if maxId >= *data.NextId {
		return fmt.Errorf("todo with id %d is invalid: todo id cannot be bigger than nextId(%d)", maxId, *data.NextId)
	}

	return nil
}

func (data todoDTO) validate() error {
	if data.Id == nil {
		return errors.New("missing field `id`")
	}

	if *data.Id < 0 {
		return errors.New("field `id` must be positive or 0")
	}

	if data.Text == nil {
		return errors.New("missing field `text`")
	}

	if data.CreatedAt == nil {
		return errors.New("missing field `createdAt`")
	}

	if *data.CreatedAt == "" {
		return errors.New("field `createdAt` cannot be empty")
	}

	if data.Completed == nil {
		return errors.New("missing field `completed`")
	}

	return nil
}

func (data todosFileDTO) toDomain() todosFile {
	var domainedTodos []todo

	for _, v := range data.Todos {
		domainedTodos = append(domainedTodos, v.toDomain())
	}

	return todosFile{
		NextId: *data.NextId,
		Todos:  domainedTodos,
	}
}

func (data todoDTO) toDomain() todo {
	return todo{
		// это хоть и указатели, но они указывают на исходную строку. Копировать строку нет смысла, она immutable
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

		/**
		Сначала мы должны проверить существует ли мас
		получаем джсон из файла -> превр. в обычный объект ->
		1 (?). Файл оказался пустым, создаем объект сами
		2. добавляем в этот объект нашу туду -> превращаем в джсон -> записываем файл заново
		*/

		_, err := os.Stat("todos.json")

		if err != nil {
			// файла не существует, сразу создадим жсон и закинем внутрь

			createdTodo := createTodoFromText(argument, 0)

			createFileAndTodos(createdTodo)

			os.Exit(0)
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

			createdTodo := createTodoFromText(argument, 0)

			createFileAndTodos(createdTodo)

			os.Exit(0)
		}

		var rawData todosFileDTO

		// создаем go struct из json:

		err = json.Unmarshal(fileContent, &rawData)

		if err != nil {
			fmt.Printf("Не удалось конвертировать JSON в структуру. Возможно, JSON поврежден. Err: %v", err)

			os.Exit(1)
		}

		err = rawData.validate()

		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		todosJsonFile := rawData.toDomain()

		createdTodo := createTodoFromText(argument, todosJsonFile.NextId)

		todosJsonFile.Todos = append(todosJsonFile.Todos, createdTodo)
		todosJsonFile.NextId++

		constructedJson, err := json.Marshal(todosJsonFile)

		if err != nil {
			fmt.Printf("Ошибка создания JSON: %v", err)
			os.Exit(1)
		}

		if err := os.WriteFile("todos.json", constructedJson, 0777); err != nil {
			fmt.Print("Error while file creating, probably security:", err)
			os.Exit(1)
		}

		fmt.Printf("Задача \"%s\" успешно добавлена.", createdTodo.Text)
	default:
		panic("NO SUCH COMMAND WE GONNA DIIIE, lol")
	}
}

func createTodoFromText(text string, id int) todo {
	return todo{
		Text:      text,
		Completed: false,
		CreatedAt: time.Now().String(),
		ID:        id,
	}
}

func createFileAndTodos(todoObj todo) {
	var todos []todo

	todosObject := todosFile{
		NextId: todoObj.ID + 1,
		Todos:  append(todos, todoObj),
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
