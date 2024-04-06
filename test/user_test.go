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

func TestCreateUser(t *testing.T) {
	ms := &storage.MockStore{}
	service := api.NewUserService(ms)
	t.Run("should return error if email is empty", func(t *testing.T) {
		payload := &storage.User{
			Email: "",
		}
		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/users/register", service.HandleUserRegister)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Error("invalid status code, it should failed")
		}
	})
	t.Run("should create a user", func(t *testing.T) {
		payload := &storage.User{
			Email:     "tom@gmail.com",
			FirstName: "Tom",
			LastName:  "Andison",
			Password:  "asd",
		}
		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/users/register", service.HandleUserRegister)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}
