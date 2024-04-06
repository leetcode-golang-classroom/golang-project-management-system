package storage

import (
	"database/sql"
)

type Store interface {
	// Users
	CreateUser(u *User) (*User, error)
	GetUserByID(id string) (*User, error)
	// Tasks
	CreateTask(t *Task) (*Task, error)
	GetTask(id string) (*Task, error)
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateTask(t *Task) (*Task, error) {
	rows, err := s.db.Exec(`INSERT INTO tasks (name, status, project_id, assigned_to) VALUES (?, 
		, ?, ?, ?)`, t.Name, t.Status, t.ProjectID, t.AssgignedToID)
	if err != nil {
		return nil, err
	}
	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}
	t.ID = id
	return t, nil
}

func (s *Storage) CreateUser(u *User) (*User, error) {
	rows, err := s.db.Exec(`INSERT INTO users (email, firstName, lastName, password) VALUES
	(?,?,?,?)`, u.Email, u.FirstName, u.LastName, u.Password)
	if err != nil {
		return nil, err
	}
	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}
	u.ID = id
	return u, nil
}

func (s *Storage) GetTask(id string) (*Task, error) {
	var t Task

	err := s.db.QueryRow("SELECT id, name, status, projectId, assignedToID, createdAt FROM tasks WHERE id = ?", id).
		Scan(&t.ID, &t.Name, &t.ProjectID, &t.AssgignedToID, &t.CreatedAt)

	return &t, err

}

func (s *Storage) GetUserByID(id string) (*User, error) {
	var u User
	err := s.db.QueryRow("SELECT id, email, firstName, lastName, createdAt FROM users WHERE id = ?", id).
		Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt)

	return &u, err
}
