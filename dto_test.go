package main

import (
	"reflect"
	"testing"
)

// Функция фабрика для создания тестовых данных
// mods - вариадический параметр, принимает функции модификаторы, которые получают на вход
// ссылку на созданный DTO объект. Таким образом мы можем модифицировать
// тестовые данные прямо в тесте
func makeTodoDTO(mods ...func(dto *todoDTO)) todoDTO {
	defaultDTO := todoDTO{
		ID:        new(0),
		Text:      new("Hi"),
		Completed: new(false),
		CreatedAt: new("2026-07-09T19:16:36+02:00"),
	}

	for _, v := range mods {
		v(&defaultDTO)
	}

	return defaultDTO
}

func makeTodosFileDTO(mods ...func(dto *todosFileDTO)) todosFileDTO {
	defaultTodosFileDTO := todosFileDTO{
		Todos:  []todoDTO{},
		NextID: new(1),
	}

	for _, v := range mods {
		v(&defaultTodosFileDTO)
	}

	return defaultTodosFileDTO
}

func TestTodoDTO_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   todoDTO
		wantErr string
	}{
		{
			name:    "correct structure",
			input:   makeTodoDTO(),
			wantErr: "",
		},
		{
			name:    "missing id",
			input:   makeTodoDTO(func(dto *todoDTO) { dto.ID = nil }),
			wantErr: "missing field `id`",
		},
		{
			name:    "negative id",
			input:   makeTodoDTO(func(dto *todoDTO) { dto.ID = new(-1) }),
			wantErr: "field `id` must be positive or 0",
		},
		{
			name:    "missing text",
			input:   makeTodoDTO(func(dto *todoDTO) { dto.Text = nil }),
			wantErr: "missing field `text`",
		},
		{
			name:    "missing createdAt",
			input:   makeTodoDTO(func(dto *todoDTO) { dto.CreatedAt = nil }),
			wantErr: "missing field `createdAt`",
		},
		{
			name:    "empty createdAt",
			input:   makeTodoDTO(func(dto *todoDTO) { dto.CreatedAt = new(string) }),
			wantErr: "field `createdAt` cannot be empty",
		},
		{
			name:    "missing completed",
			input:   makeTodoDTO(func(dto *todoDTO) { dto.Completed = nil }),
			wantErr: "missing field `completed`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.validate()
			// если ожидаем ошибку
			checkErr(t, err, tt.wantErr)
		})
	}
}

func TestTodosFileDTO_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   todosFileDTO
		wantErr string
	}{
		{
			name:    "correct structure with empty todos",
			input:   makeTodosFileDTO(),
			wantErr: "",
		},
		{
			name: "correct structure with todos",
			input: makeTodosFileDTO(func(dto *todosFileDTO) {
				for i := 0; i < 3; i++ {
					dto.Todos = append(dto.Todos, makeTodoDTO(func(dto *todoDTO) { dto.ID = new(i) }))
					*dto.NextID++
				}
			}),
			wantErr: "",
		},
		{
			name:    "missing nextId",
			input:   makeTodosFileDTO(func(dto *todosFileDTO) { dto.NextID = nil }),
			wantErr: "missing field `nextId`",
		},
		{
			name:    "negative nextId",
			input:   makeTodosFileDTO(func(dto *todosFileDTO) { dto.NextID = new(-1) }),
			wantErr: "field nextId must be positive or 0",
		},
		{
			name:    "missing todos",
			input:   makeTodosFileDTO(func(dto *todosFileDTO) { dto.Todos = nil }),
			wantErr: "missing field `todos`",
		},
		{
			name: "invalid todo",
			input: makeTodosFileDTO(func(dto *todosFileDTO) {
				dto.Todos = []todoDTO{
					makeTodoDTO(func(dto *todoDTO) { dto.ID = nil }),
				}
			}),
			wantErr: "todo at index 0 is invalid",
		},
		{
			name: "duplicate todo id",
			input: makeTodosFileDTO(func(dto *todosFileDTO) {
				dto.Todos = []todoDTO{
					makeTodoDTO(func(dto *todoDTO) { dto.ID = new(1) }),
					makeTodoDTO(func(dto *todoDTO) { dto.ID = new(1) }),
					makeTodoDTO(func(dto *todoDTO) { dto.ID = new(1) }),
				}
			}),
			wantErr: "todo at index 1 is invalid",
		},
		{
			name: "todo id equals nextId",
			input: makeTodosFileDTO(func(dto *todosFileDTO) {
				dto.NextID = new(5)
				dto.Todos = []todoDTO{
					makeTodoDTO(func(dto *todoDTO) { dto.ID = new(5) }),
				}
			}),
			wantErr: "todo with id 5 is invalid: todo id cannot be bigger or equal than nextId(5)",
		},
		{
			name: "todo id bigger than nextId",
			input: makeTodosFileDTO(func(dto *todosFileDTO) {
				dto.NextID = new(5)
				dto.Todos = []todoDTO{
					makeTodoDTO(func(dto *todoDTO) { dto.ID = new(6) }),
				}
			}),
			wantErr: "todo with id 6 is invalid: todo id cannot be bigger or equal than nextId(5)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.validate()
			// если ожидаем ошибку
			checkErr(t, err, tt.wantErr)
		})
	}
}

func TestTodoList_RoundTrip(t *testing.T) {
	original := todoList{
		NextID: 1,
		Todos: []todo{
			{
				ID:        0,
				CreatedAt: "2026-07-09T19:16:36+02:00",
				Completed: false,
				Text:      "testing",
			},
		},
	}

	// Нет смысла проверять, что структура после toDTO невалидна, это гарантирует
	// validate на этапе после анмаршала
	got := original.toDTO().toDomain()

	if !reflect.DeepEqual(original, got) {
		t.Errorf("round trip изменил данные:\nбыло: %+v\nстало: %+v", original, got)
	}
}
