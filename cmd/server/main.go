// Package for server
package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"
)

type Todo struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	Completed bool      `json:"completed"`
}

type CreateTodoRequest struct {
	Text string `json:"text"`
}

type PatchTodoRequest struct {
	Completed *bool `json:"completed"`
}

var todoCache = make(map[int]Todo)

var cacheMutex sync.RWMutex

var nextTodoID = 1

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("GET /todos", handleGetTodos)
	mux.HandleFunc("GET /todos/{id}", handleGetTodoByID)
	mux.HandleFunc("POST /todos", handleCreateTodo)
	mux.HandleFunc("DELETE /todos/{id}", handleDeleteTodoByID)
	mux.HandleFunc("PATCH /todos/{id}", handlePatchTodoByID)

	fmt.Println("Server started on http://localhost:8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("fatal: unable to start server: %v", err)
		return
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Body)

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello, World!"))
	if err != nil {
		return
	}
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if req.Text == "" {
		http.Error(w, "no text for todo was provided", http.StatusUnprocessableEntity)

		return
	}

	cacheMutex.Lock()

	id := nextTodoID
	nextTodoID++

	createdTodo := Todo{
		ID:        id,
		CreatedAt: time.Now(),
		Completed: false,
		Text:      req.Text,
	}

	todoCache[id] = createdTodo

	cacheMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")

	j, err := json.Marshal(createdTodo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(j)

	if err != nil {
		return
	}
}

func handleGetTodos(w http.ResponseWriter, _ *http.Request) {
	cacheMutex.RLock()
	todos := slices.Collect(maps.Values(todoCache))
	cacheMutex.RUnlock()

	slices.SortFunc(todos, func(a, b Todo) int {
		return cmp.Compare(a.ID, b.ID)
	})

	j, err := json.Marshal(todos)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(j)

	if err != nil {
		return
	}
}

func handleGetTodoByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	cacheMutex.RLock()
	todo, ok := todoCache[id]
	cacheMutex.RUnlock()

	if !ok {
		http.Error(w, "todo not found", http.StatusNotFound)

		return
	}

	j, err := json.Marshal(todo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(j)

	if err != nil {
		return
	}
}

func handleDeleteTodoByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	cacheMutex.Lock()

	_, ok := todoCache[id]
	if ok {
		delete(todoCache, id)
	}

	cacheMutex.Unlock()

	if !ok {
		http.Error(w, "todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handlePatchTodoByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	var req PatchTodoRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if req.Completed == nil {
		http.Error(w, "no completion status was provided", http.StatusUnprocessableEntity)

		return
	}

	cacheMutex.Lock()
	todo, ok := todoCache[id]

	if ok {
		todo.Completed = *req.Completed
		todoCache[id] = todo
	}

	cacheMutex.Unlock()

	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)

		return
	}

	j, err := json.Marshal(todo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(j)

	if err != nil {
		return
	}
}
