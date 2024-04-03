package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/leetcode-golang-classroom/golang-project-management-system/api"
	"github.com/leetcode-golang-classroom/golang-project-management-system/internal/storage"
)

func TestCreateTask(t *testing.T) {
	ms := &storage.MockStore{}
	service := api.NewTaskService(ms)
	t.Run("should return an error if name is empty", func(t *testing.T) {
		payload := &storage.Task{
			Name: "",
		}
		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/tasks", service.HandleCreateTask)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Error("invalid status code, it should failed")
		}
	})
	t.Run("should create a task", func(t *testing.T) {
		payload := &storage.Task{
			Name:          "Creating a REST API in go",
			ProjectID:     1,
			AssgignedToID: 42,
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/tasks", service.HandleCreateTask)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

func TestGetTask(t *testing.T) {
	ms := &storage.MockStore{}
	service := api.NewTaskService(ms)

	t.Run("should return the task", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/tasks/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/tasks/{id}", service.HandleGetTask)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("invalid status code, it should failed")
		}
	})
}
