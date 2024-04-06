package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/leetcode-golang-classroom/golang-project-management-system/config"
	"github.com/leetcode-golang-classroom/golang-project-management-system/internal/storage"
)

type UserService struct {
	store storage.Store
}

var errEmailRequired = errors.New("email is required")
var errFirstNameRequired = errors.New("firstName is required")
var errLastNameRequired = errors.New("lastName is required")
var errPasswordRequired = errors.New("password is required")

func NewUserService(s storage.Store) *UserService {
	return &UserService{store: s}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", s.HandleUserRegister).Methods("POST")
}

func (s *UserService) HandleUserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var payload *storage.User
	err = json.Unmarshal(body, &payload)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, storage.ErrrorResponse{Error: "Invalid request payload"})
		return
	}

	if err := validateUserPayload(payload); err != nil {
		WriteJSON(w, http.StatusBadRequest, storage.ErrrorResponse{Error: err.Error()})
		return
	}
	hashedPW, err := HashPassword(payload.Password)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, storage.ErrrorResponse{Error: err.Error()})
		return
	}
	payload.Password = hashedPW
	u, err := s.store.CreateUser(payload)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, storage.ErrrorResponse{Error: "Error creating user"})
		return
	}
	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, storage.ErrrorResponse{Error: "Error creating session"})
		return
	}
	WriteJSON(w, http.StatusCreated, token)
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	secret := []byte(config.Envs.JWTSecret)
	token, err := CreateJWT(secret, id)
	if err != nil {
		return "", err
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}

func CreateJWT(secret []byte, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateUserPayload(user *storage.User) error {
	if user.Email == "" {
		return errEmailRequired
	}
	if user.FirstName == "" {
		return errFirstNameRequired
	}
	if user.LastName == "" {
		return errLastNameRequired
	}
	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}
