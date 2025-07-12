package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	tasks  []Task
	nextID int
	mutex  = &sync.Mutex{}
)

const dbFile = "tasks.json"

func loadTasks() {
	mutex.Lock()
	defer mutex.Unlock()
	file, err := os.ReadFile(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			tasks = []Task{}
			nextID = 1
			return
		}
		log.Fatalf("Error reading tasks file: %v", err)
	}
	if err := json.Unmarshal(file, &tasks); err != nil {
		log.Fatalf("Error unmarshalling tasks: %v", err)
	}
	if len(tasks) > 0 {
		nextID = tasks[len(tasks)-1].ID + 1
	} else {
		nextID = 1
	}
}

func saveTasks() {
	mutex.Lock()
	defer mutex.Unlock()
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling tasks: %v", err)
	}
	if err := os.WriteFile(dbFile, data, 0644); err != nil {
		log.Fatalf("Error writing tasks file: %v", err)
	}
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		getTasks(w, r)
	case http.MethodPost:
		createTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodPut:
		updateTask(w, r)
	case http.MethodDelete:
		deleteTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mutex.Lock()
	task.ID = nextID
	nextID++
	tasks = append(tasks, task)
	mutex.Unlock()
	saveTasks()
	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, task := range tasks {
		if task.ID == id {
			tasks[i] = updatedTask
			saveTasks()
			json.NewEncoder(w).Encode(updatedTask)
			return
		}
	}

	http.NotFound(w, r)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			saveTasks()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.NotFound(w, r)
}

func main() {
	loadTasks()

	http.Handle("/", http.FileServer(http.Dir("./frontend/")))

	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/tasks/", taskHandler)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
