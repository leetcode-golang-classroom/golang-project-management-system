package storage

import "time"

type Task struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	ProjectID     int64     `json:"projectID"`
	AssgignedToID int64     `json:"assignedTo"`
	CreatedAt     time.Time `json:"createdAt"`
}

type ErrrorResponse struct {
	Error string `json:"error"`
}

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}
