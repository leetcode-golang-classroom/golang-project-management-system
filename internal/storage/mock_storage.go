package storage

// Mocks

type MockStore struct{}

func (m *MockStore) CreateUser() error {
	return nil
}

func (m *MockStore) CreateTask(t *Task) (*Task, error) {
	return &Task{}, nil
}

func (m *MockStore) GetTask(id string) (*Task, error) {
	return &Task{}, nil
}

func (m *MockStore) GetUserByID(id string) (*User, error) {
	return &User{}, nil
}
