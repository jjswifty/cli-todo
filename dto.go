package main

import (
	"errors"
	"fmt"
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

	idsMap := make(map[int]bool, len(data.Todos))

	maxID := 0

	for i, dto := range data.Todos {
		if err := dto.validate(); err != nil {
			return fmt.Errorf("todo at index %d is invalid: %w", i, err)
		}

		id := *dto.ID

		if id > maxID {
			maxID = *dto.ID
		}

		if idsMap[id] {
			return fmt.Errorf("todo at index %d is invalid: duplicate ID", i)
		}

		idsMap[id] = true
	}

	if maxID >= *data.NextID {
		return fmt.Errorf("todo with id %d is invalid: todo id cannot be bigger or equal than nextId(%d)", maxID, *data.NextID)
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

// Конвертирует DTO в доменную модель.
// Предусловие: DTO должен пройти validate() — поля разыменовываются без проверок.
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

// Конвертирует DTO в доменную модель.
// Предусловие: DTO должен пройти validate() — поля разыменовываются без проверок.
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
