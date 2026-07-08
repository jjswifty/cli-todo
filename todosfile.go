package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"
)

type todosFileDTO struct {
	NextID *int      `json:"nextId"`
	Todos  []todoDTO `json:"todos"`
}
type todoDTO struct {
	ID        *int    `json:"id"`
	CreatedAt *string `json:"createdAt"`
	Completed *bool   `json:"completed"`
	Text      *string `json:"text"`
}

type todoList struct {
	NextID int
	Todos  []todo
}

type todo struct {
	ID        int
	CreatedAt string
	Completed bool
	Text      string
}

func (data todosFileDTO) validate() error {
	if data.NextID == nil {
		return errors.New("missing field `nextId`")
	}

	if *data.NextID < 0 {
		return errors.New("field nextId must be positive or 0")
	}

	// todos может и не быть, нельзя сразу распаковывать
	if data.Todos == nil {
		return errors.New("missing field `todos`")
	}

	if len(data.Todos) == 0 {
		return nil
	}

	// todo: переписать на мапу, так как здесь сейчас o(n^2)
	var ids []int

	maxID := 0

	for i, dto := range data.Todos {
		if err := dto.validate(); err != nil {
			return fmt.Errorf("todo at index %d is invalid: %w", i, err)
		}

		id := *dto.ID

		if id > maxID {
			maxID = *dto.ID
		}

		if slices.Contains(ids, id) {
			return fmt.Errorf("todo at index %d is invalid: duplicate ID", i)
		}

		ids = append(ids, id)
	}

	if maxID >= *data.NextID {
		return fmt.Errorf("todo with id %d is invalid: todo id cannot be bigger than nextId(%d)", maxID, *data.NextID)
	}

	return nil
}

func (data todoDTO) validate() error {
	if data.ID == nil {
		return errors.New("missing field `id`")
	}

	if *data.ID < 0 {
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

func (data todosFileDTO) toDomain() todoList {
	var domainedTodos []todo

	for _, v := range data.Todos {
		domainedTodos = append(domainedTodos, v.toDomain())
	}

	return todoList{
		NextID: *data.NextID,
		Todos:  domainedTodos,
	}
}

func (data todoDTO) toDomain() todo {
	return todo{
		// это хоть и указатели, но они указывают на исходную строку. Копировать строку нет смысла, она immutable
		ID:        *data.ID,
		CreatedAt: *data.CreatedAt,
		Text:      *data.Text,
		Completed: *data.Completed,
	}
}

func (t todo) toDTO() todoDTO {
	return todoDTO{
		ID:        &t.ID,
		CreatedAt: &t.CreatedAt,
		Completed: &t.Completed,
		Text:      &t.Text,
	}
}

func (list todoList) toDTO() todosFileDTO {
	dtos := make([]todoDTO, 0, len(list.Todos))

	for _, t := range list.Todos {
		dtos = append(dtos, t.toDTO())
	}

	return todosFileDTO{
		NextID: &list.NextID,
		Todos:  dtos,
	}
}

func loadTodos() (todoList, error) {
	fileContent, err := os.ReadFile("todos.json")

	if errors.Is(err, os.ErrNotExist) {
		// файла нет — начинаем с чистого листа, это не ошибка
		return todoList{}, nil
	}

	if err != nil {
		return todoList{}, fmt.Errorf("failed to read todos.json: %w", err)
	}

	if len(fileContent) == 0 {
		return todoList{}, errors.New("todos.json существует, но пустой — удалите его или восстановите содержимое")
	}

	var rawData todosFileDTO

	// создаем go struct из json:
	err = json.Unmarshal(fileContent, &rawData)

	if err != nil {
		return todoList{}, fmt.Errorf("не удалось конвертировать JSON в структуру. Возможно, JSON поврежден. Err: %w", err)
	}

	err = rawData.validate()

	if err != nil {
		return todoList{}, fmt.Errorf("invalid todos file: %w", err)
	}

	return rawData.toDomain(), nil
}

func saveTodos(list todoList) error {
	constructedJSON, err := json.Marshal(list.toDTO())

	if err != nil {
		return fmt.Errorf("ошибка создания JSON: %w", err)
	}

	if err := os.WriteFile("todos.json", constructedJSON, 0777); err != nil {
		return fmt.Errorf("error while file creating, probably security: %w", err)
	}

	return nil
}

func newTodoFromText(text string, id int) todo {
	return todo{
		Text:      text,
		Completed: false,
		CreatedAt: time.Now().Format(time.RFC3339),
		ID:        id,
	}
}
