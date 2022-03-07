package mysql

import (
	"database/sql"

	"github.com/erodrigufer/goTennis/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

// Insert a new record to the users table
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate login attempt with email and password
// If correct, return the user ID
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Fetch details for a specific user based on its user ID
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
