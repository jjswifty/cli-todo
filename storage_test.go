package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadTodos(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    todoList
		wantErr string
	}{
		{
			name: "correct json",
			path: filepath.Join("testdata", "storage", "valid-todos.json"),
			want: todoList{
				NextID: 4,
				Todos: []todo{
					{ID: 0, CreatedAt: "2026-07-08T23:31:15+02:00", Completed: false, Text: "make dinner"},
					{ID: 1, CreatedAt: "2026-07-09T19:16:36+02:00", Completed: false, Text: "make dinner"},
					{ID: 2, CreatedAt: "2026-07-09T19:16:46+02:00", Completed: false, Text: "make dinner asd"},
					{ID: 3, CreatedAt: "2026-07-09T19:16:52+02:00", Completed: false, Text: "make dinner asd as"},
				},
			},
			wantErr: "",
		},
		{
			name:    "handles wrong filepath",
			path:    filepath.Join("testdata", "storage", "non-exist.json"),
			wantErr: "",
			want:    todoList{},
		},
		{
			name:    "read error",
			path:    t.TempDir(),
			wantErr: "failed to read file",
			want:    todoList{},
		},
		{
			name:    "empty json",
			path:    filepath.Join("testdata", "storage", "empty-file.json"),
			wantErr: "файл по пути testdata/storage/empty-file.json существует, но пустой — удалите его или восстановите содержимое",
			want:    todoList{},
		},
		{
			name:    "broken json",
			path:    filepath.Join("testdata", "storage", "broken-file"),
			wantErr: "не удалось конвертировать JSON в структуру. Возможно, JSON поврежден.",
			want:    todoList{},
		},
		{
			name:    "read error",
			path:    filepath.Join("testdata", "storage", "invalid-todos-next-id-bigger.json"),
			wantErr: "invalid todos file",
			want:    todoList{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadTodos(tt.path)

			checkErr(t, err, tt.wantErr)

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("несовпадение структур:\nожидали: %+v\nполучили: %+v", tt.want, got)
			}
		})
	}
}

func TestSaveTodos_WriteError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nope", "todos.json") // папки nope не существует

	err := saveTodos(todoList{}, path)

	checkErr(t, err, "error while file creating")
}

func TestSaveTodos_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "todos.json")

	original := todoList{
		NextID: 1,
		Todos: []todo{
			{ID: 0, CreatedAt: "2026-07-08T23:31:15+02:00", Completed: false, Text: "make dinner"},
		},
	}

	if err := saveTodos(original, path); err != nil {
		t.Fatalf("unexpected save error: %q", err)
	}

	got, err := loadTodos(path)

	if err != nil {
		t.Fatalf("unexpected load error: %q", err)
	}

	if !reflect.DeepEqual(original, got) {
		t.Errorf("round trip изменил данные:\nбыло:  %+v\nстало: %+v", original, got)
	}
}
