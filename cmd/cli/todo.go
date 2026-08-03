package main

import (
	"time"
)

type todoList struct {
	NextID int
	Todos  []todo
}

type todo struct {
	ID        int
	CreatedAt time.Time
	Completed bool
	Text      string
}

func newTodoFromText(text string, id int) todo {
	return todo{
		Text:      text,
		Completed: false,
		CreatedAt: time.Now(),
		ID:        id,
	}
}
