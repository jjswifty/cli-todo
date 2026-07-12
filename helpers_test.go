package main

import (
	"strings"
	"testing"
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
