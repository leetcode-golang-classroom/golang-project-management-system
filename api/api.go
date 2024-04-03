package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leetcode-golang-classroom/golang-project-management-system/internal/storage"
)

type APIServer struct {
	addr  string
	store storage.Store
}

func NewAPIServer(addr string, store storage.Store) *APIServer {
	return &APIServer{addr: addr, store: store}
}
func (s *APIServer) Serve() {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	// registering our services
	taskService := NewTaskService(s.store)
	taskService.RegisterRoutes(subRouter)

	log.Println("Starting the API server at", s.addr)

	log.Fatal(http.ListenAndServe(s.addr, subRouter))

}
