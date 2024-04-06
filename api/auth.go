package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/leetcode-golang-classroom/golang-project-management-system/config"
	"github.com/leetcode-golang-classroom/golang-project-management-system/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from request (auth header)
		tokenString := GetTokenFromRequest(r)
		// vaildate the token
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Println("failed to authenticate token", err)
			permissionDenied(w)
			return
		}
		if !token.Valid {
			log.Println("failed to authenticate token")
			permissionDenied(w)
			return
		}
		// get the userId from the token
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userID"].(string)
		_, err = store.GetUserByID(userID)
		if err != nil {
			log.Println("failed to get user")
			permissionDenied(w)
			return
		}
		// call the handler func and continue to the endpoint
		handlerFunc(w, r)
	}
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, storage.ErrrorResponse{
		Error: fmt.Errorf("permission denied").Error(),
	})
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}
	return ""
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := config.Envs.JWTSecret
	return jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func HashPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
