// Package to explicitly handle the sessions mysql db
package mysql

import (
	"database/sql"

	"github.com/erodrigufer/GoTennis/pkg/models"
)

// Define a SessionModel type which wraps a sql.DB connection pool
type SessionModel struct {
	DB *sql.DB
}

// Insert new session into the db
func (m *SessionModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

// Get Session from db, using its id
func (m *SessionModel) Get(id int) (*models.Session, error) {
	return nil, nil
}

// Return the 10 most recently created sessions
func (m *Sessionmodel) Latest() ([]*models.Session, error) {
	return nil, nil
}
