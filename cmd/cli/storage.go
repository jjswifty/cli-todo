package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func loadTodos(path string) (todoList, error) {
	fileContent, err := os.ReadFile(path)

	if errors.Is(err, os.ErrNotExist) {
		// файла нет — начинаем с чистого листа, это не ошибка
		return todoList{}, nil
	}

	if err != nil {
		return todoList{}, fmt.Errorf("%s failed to read file: %w", path, err)
	}

	if len(fileContent) == 0 {
		return todoList{}, fmt.Errorf("файл по пути %s существует, но пустой — удалите его или восстановите содержимое", path)
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

func saveTodos(list todoList, path string) error {
	constructedJSON, err := json.Marshal(list.toDTO())

	if err != nil {
		return fmt.Errorf("ошибка создания JSON: %w", err)
	}

	if err := os.WriteFile(path, constructedJSON, 0644); err != nil {
		return fmt.Errorf("error while file creating, probably security: %w", err)
	}

	return nil
}
