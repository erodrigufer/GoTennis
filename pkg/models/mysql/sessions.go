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
	// SQL-command to execute, `` to write command over 2 lines for readability
	// ? is a placeholder parameter, since we would otherwise be using untrusted
	// unsanitized user input data
	stmt := `INSERT INTO sessions (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	// Use the Exec() method on the embedded connection pool to execute the
	// statement. The first parameter is the SQL statement, followed by the
	// title, content and expiry values for the placeholder parameters. This
	// method returns a sql.Result object, which contains some basic
	// information about what happened when the statement was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	// Use the LastInsertId() method on the result object to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// Get Session from db, using its id
func (m *SessionModel) Get(id int) (*models.Session, error) {
	return nil, nil
}

// Return the 10 most recently created sessions
func (m *Sessionmodel) Latest() ([]*models.Session, error) {
	return nil, nil
}
