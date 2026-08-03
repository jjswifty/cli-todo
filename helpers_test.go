package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

func checkErr(t *testing.T, err error, wantErr string) {
	t.Helper()

	// в кейсе не ожидается ошибка, значит мы не можем сюда попасть
	if wantErr == "" {
		if err != nil {
			t.Fatalf("unexpected error: %q", err.Error())
		}

		return
	}

	// если ожидали, но не получили
	if err == nil {
		t.Fatalf("ожидали ошибку %q, но не получили ее.", wantErr)
	}

	if !strings.Contains(err.Error(), wantErr) {
		t.Errorf("ожидали ошибку %q, но получили %q", wantErr, err.Error())
	}
}

func mustParseTime(t *testing.T, s string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("некорректная дата в тесте: %v", err)
	}

	return parsed
}

func mustWriteFile(t *testing.T, path string, content todoList) {
	t.Helper()

	data, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("не удалось сериализовать тестовые задачи: %v", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("не удалось подготовить тестовый файл: %v", err)
	}
}
