package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set up test environment
	dbFile = "test_tasks.json"
	loadTasks()
	code := m.Run()
	// Clean up
	os.Remove(dbFile)
	os.Exit(code)
}

func TestGetTasks(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tasksHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var tasks []Task
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
}

func TestCreateTask(t *testing.T) {
	var jsonStr = []byte(`{"title":"Test Task","done":false}`)
	req, err := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tasksHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var task Task
	if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if task.Title != "Test Task" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			task.Title, "Test Task")
	}
}
