package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leetcode-golang-classroom/golang-project-management-system/internal/storage"
)

type TaskService struct {
	store storage.Store
}

var errNameRequired = errors.New("name is required")
var errProjectIDRequired = errors.New("project id is required")
var errUserIDRequired = errors.New("user id is required")

func NewTaskService(s storage.Store) *TaskService {
	return &TaskService{store: s}
}
func (s *TaskService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", WithJWTAuth(s.HandleCreateTask, s.store)).Methods("POST")
	r.HandleFunc("/tasks/{id}", WithJWTAuth(s.HandleGetTask, s.store)).Methods("GET")
}

func (s *TaskService) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, storage.ErrrorResponse{Error: "Invalid request payload"})
		return
	}
	defer r.Body.Close()
	var task *storage.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, storage.ErrrorResponse{Error: "Invalid request payload"})
		return
	}
	if err := validateTaskPayload(task); err != nil {
		WriteJSON(w, http.StatusBadRequest, storage.ErrrorResponse{Error: err.Error()})
		return
	}

	t, err := s.store.CreateTask(task)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, storage.ErrrorResponse{Error: "Error creating task"})
		return
	}
	WriteJSON(w, http.StatusCreated, t)
}

func (s *TaskService) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		WriteJSON(w, http.StatusBadRequest, storage.ErrrorResponse{Error: "id is required"})
		return
	}

	t, err := s.store.GetTask(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, storage.ErrrorResponse{Error: "task not found"})
		return
	}
	WriteJSON(w, http.StatusOK, t)
}

func validateTaskPayload(t *storage.Task) error {
	if t.Name == "" {
		return errNameRequired
	}
	if t.ProjectID == 0 {
		return errProjectIDRequired
	}
	if t.AssgignedToID == 0 {
		return errUserIDRequired
	}
	return nil
}
