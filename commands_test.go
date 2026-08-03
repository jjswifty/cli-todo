package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"slices"
	"testing"
	"time"
)

func TestRequireTodos(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    todoList
		wantErr string
	}{
		{
			name: "correctly required todos",
			path: filepath.Join("testdata", "commands", "valid-todos.json"),
			want: todoList{
				NextID: 4,
				Todos: []todo{
					{ID: 0, CreatedAt: mustParseTime(t, "2026-07-08T23:31:15+02:00"), Completed: false, Text: "make dinner"},
					{ID: 1, CreatedAt: mustParseTime(t, "2026-07-09T19:16:36+02:00"), Completed: false, Text: "make dinner"},
					{ID: 2, CreatedAt: mustParseTime(t, "2026-07-09T19:16:46+02:00"), Completed: false, Text: "make dinner asd"},
					{ID: 3, CreatedAt: mustParseTime(t, "2026-07-09T19:16:52+02:00"), Completed: false, Text: "make dinner asd as"},
				},
			},
			wantErr: "",
		},
		{
			name: "returns err if todo list is empty",
			path: filepath.Join("testdata", "commands", "no-todos.json"),
			want: todoList{
				NextID: 0,
				Todos:  nil,
			},
			wantErr: "список задач пуст",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := requireTodos(tt.path)

			fmt.Printf("got: %+v", got)

			checkErr(t, err, tt.wantErr)

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("несовпадение результатов:\nожидали: %+v\nполучили: %+v", tt.want, got)
			}
		})
	}
}

func TestRunAdd(t *testing.T) {
	t.Run("correctly appended todo with add command", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "todos.json")

		baseContent := todoList{
			NextID: 0,
			Todos:  []todo{},
		}

		mustWriteFile(t, path, baseContent)

		err := runAdd(path, "Testing todo")

		if err != nil {
			t.Fatalf("Unexpected error while adding: %v", err)
		}

		todosAfterWrite, err := loadTodos(path)

		if err != nil {
			t.Fatalf("Fatal: error while loading todos: %v", err)
		}

		if len(todosAfterWrite.Todos) != 1 {
			t.Fatalf("Fatal: expected only 1 todo to be appended, got: %d", len(todosAfterWrite.Todos))
		}

		if todosAfterWrite.NextID != 1 {
			t.Errorf("NextID should increment after todo append")
		}

		writtenTodo := todosAfterWrite.Todos[0]

		if writtenTodo.Text != "Testing todo" {
			t.Errorf("added todo got wrong text")
		}

		if writtenTodo.Completed != false {
			t.Errorf("added todo got wrong completion status, expected false")
		}

		if writtenTodo.ID != 0 {
			t.Errorf("added todo should have id 0, got: %d", writtenTodo.ID)
		}
	})

	t.Run("returns error for invalid path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing", "todos.json")

		err := runAdd(path, "Testing todo")

		checkErr(t, err, "error while file creating")
	})

	t.Run("returns error when todos file is not writable", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("проверка Unix-прав файла неприменима на Windows")
		}

		path := filepath.Join(t.TempDir(), "todos.json")
		mustWriteFile(t, path, todoList{NextID: 0, Todos: []todo{}})

		if err := os.Chmod(path, 0o400); err != nil {
			t.Fatalf("не удалось убрать право на запись у тестового файла: %v", err)
		}

		t.Cleanup(func() {
			if err := os.Chmod(path, 0o600); err != nil {
				t.Errorf("не удалось восстановить права тестового файла: %v", err)
			}
		})

		err := runAdd(path, "Testing todo")

		checkErr(t, err, "error while file creating")
	})
}

func TestRunDone(t *testing.T) {
	baseContent := todoList{
		NextID: 1,
		Todos: []todo{
			{ID: 0, CreatedAt: mustParseTime(t, "2026-07-08T23:31:15+02:00"), Completed: false, Text: "make dinner"},
		},
	}

	t.Run("correctly marks todo as completed", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "todos.json")
		mustWriteFile(t, path, baseContent)

		err := runDone(path, 0)
		if err != nil {
			t.Fatalf("unexpected error while marking todo as completed: %v", err)
		}

		loadedTodos, err := loadTodos(path)
		if err != nil {
			t.Fatalf("unexpected error while loading testing data: %v", err)
		}

		if len(loadedTodos.Todos) != 1 {
			t.Fatalf("expected one todo after runDone, got %d", len(loadedTodos.Todos))
		}

		if !loadedTodos.Todos[0].Completed {
			t.Error("todo was not marked as completed")
		}
	})

	t.Run("returns error if todo with id does not exist", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "todos.json")
		mustWriteFile(t, path, baseContent)

		err := runDone(path, 5)

		checkErr(t, err, "todo с id 5 не найдено")
	})
}

func TestRunUndone(t *testing.T) {
	baseContent := todoList{
		NextID: 1,
		Todos: []todo{
			{ID: 0, CreatedAt: mustParseTime(t, "2026-07-08T23:31:15+02:00"), Completed: true, Text: "make dinner"},
		},
	}

	t.Run("correctly marks todo as not completed", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "todos.json")
		mustWriteFile(t, path, baseContent)

		err := runUndone(path, 0)
		if err != nil {
			t.Fatalf("unexpected error while marking todo as not completed: %v", err)
		}

		loadedTodos, err := loadTodos(path)
		if err != nil {
			t.Fatalf("unexpected error while loading testing data: %v", err)
		}

		if len(loadedTodos.Todos) != 1 {
			t.Fatalf("expected one todo after runUndone, got %d", len(loadedTodos.Todos))
		}

		if loadedTodos.Todos[0].Completed {
			t.Error("todo was not marked as not completed")
		}
	})

	t.Run("returns error if todo with id does not exist", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "todos.json")
		mustWriteFile(t, path, baseContent)

		err := runUndone(path, 5)

		checkErr(t, err, "todo с id 5 не найдено")
	})
}

func TestRunDelete(t *testing.T) {
	baseContent := todoList{
		NextID: 1,
		Todos: []todo{
			{ID: 0, CreatedAt: mustParseTime(t, "2026-07-08T23:31:15+02:00"), Completed: false, Text: "make dinner"},
		},
	}

	t.Run("correctly deletes todo", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "todos.json")
		mustWriteFile(t, path, baseContent)

		err := runDelete(path, 0)
		if err != nil {
			t.Fatalf("unexpected error while deleting todo: %v", err)
		}

		loadedTodos, err := loadTodos(path)
		if err != nil {
			t.Fatalf("unexpected error while loading testing data: %v", err)
		}

		if len(loadedTodos.Todos) != 0 {
			t.Errorf("expected todo to be deleted, got %d todos", len(loadedTodos.Todos))
		}
	})

	t.Run("returns error if todo with id does not exist", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "todos.json")
		mustWriteFile(t, path, baseContent)

		err := runDelete(path, 5)

		checkErr(t, err, "todo with id 5 not found")
	})
}

func TestDeleteTodoByID(t *testing.T) {
	baseContent := todoList{
		NextID: 4,
		Todos: []todo{
			{ID: 0, CreatedAt: mustParseTime(t, "2026-07-08T23:31:15+02:00"), Completed: false, Text: "make dinner"},
			{ID: 1, CreatedAt: mustParseTime(t, "2026-07-09T19:16:36+02:00"), Completed: false, Text: "make dinner 2"},
			{ID: 2, CreatedAt: mustParseTime(t, "2026-07-09T19:16:46+02:00"), Completed: false, Text: "make dinner asd"},
			{ID: 3, CreatedAt: mustParseTime(t, "2026-07-09T19:16:52+02:00"), Completed: false, Text: "make dinner asd as"},
		},
	}

	t.Run("returns error for invalid path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "missing", "todos.json")

		err := deleteTodoByID(path, 1)

		checkErr(t, err, "список задач пуст")
	})

	t.Run("correctly delete todo from list", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "todos.json")

		mustWriteFile(t, path, baseContent)

		err := deleteTodoByID(path, 1)

		if err != nil {
			t.Fatalf("unexpected error occurred while todo deletion: %v", err)
		}

		todosAfterDeletion, err := loadTodos(path)

		if err != nil {
			t.Fatalf("unexpected error occurred while loading testing data: %v", err)
		}

		stillExists := slices.ContainsFunc(todosAfterDeletion.Todos, func(t todo) bool {
			return t.ID == 1
		})

		if stillExists {
			t.Errorf("element was not successfully deleted")
		}
	})

	t.Run("returns error if todo with id does not exist", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "todos.json")

		mustWriteFile(t, path, baseContent)

		err := deleteTodoByID(path, 5)

		checkErr(t, err, "todo with id 5 not found")
	})
}

func TestToggleTodoCompletionStatus(t *testing.T) {
	baseContent := todoList{
		NextID: 1,
		Todos: []todo{
			{ID: 0, CreatedAt: mustParseTime(t, "2026-07-08T23:31:15+02:00"), Completed: false, Text: "make dinner"},
		},
	}

	t.Run("correctly changes todo completion status", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "todos.json")
		mustWriteFile(t, path, baseContent)

		err := toggleTodoCompletionStatus(path, 0, true)
		if err != nil {
			t.Fatalf("unexpected error while changing todo completion status: %v", err)
		}

		loadedTodos, err := loadTodos(path)
		if err != nil {
			t.Fatalf("unexpected error while loading testing data: %v", err)
		}

		if !loadedTodos.Todos[0].Completed {
			t.Error("todo completion status was not changed")
		}
	})

	t.Run("returns error if todo with id does not exist", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "todos.json")
		mustWriteFile(t, path, baseContent)

		err := toggleTodoCompletionStatus(path, 5, true)

		checkErr(t, err, "todo с id 5 не найдено")
	})
}

type testInputStruct struct {
	list []todo
	id   int
}

func TestGetTodoByIDFromList(t *testing.T) {
	t.Run("incorrect id", func(t *testing.T) {
		input := testInputStruct{
			id: -2323,
			list: []todo{{
				ID:        5,
				CreatedAt: time.Time{},
				Completed: false,
				Text:      "mock",
			}, {
				ID:        2,
				CreatedAt: time.Time{},
				Completed: false,
				Text:      "mock",
			}},
		}

		got, err := getTodoByIDFromList(input.list, input.id)

		checkErr(t, err, "todo с id -2323 не найдено")

		if got != nil {
			t.Fatalf("функция не должна была ничего вернуть, однако мы получили %+v", got)
		}
	})

	t.Run("returns correct ptr to original backing array", func(t *testing.T) {
		input := testInputStruct{
			id: 2,
			list: []todo{{
				ID:        5,
				CreatedAt: time.Time{},
				Completed: false,
				Text:      "mock",
			}, {
				ID:        2,
				CreatedAt: time.Time{},
				Completed: false,
				Text:      "mock",
			}},
		}

		originalPtr := &input.list[1]

		got, err := getTodoByIDFromList(input.list, input.id)

		checkErr(t, err, "")

		if originalPtr != got {
			t.Errorf("несовпадение адресов:\nожидали: %p\nполучили: %p", originalPtr, got)
		}
	})
}
