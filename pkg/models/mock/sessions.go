package mock

import (
	"time"

	"github.com/erodrigufer/GoTennis/pkg/models"
)

var mockSession = &models.Session{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SessionModel struct{}

// Insert a new session into the db, it returns the id of the newly inserted
// row in the db
func (m *SessionModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}
func (m *SessionModel) Get(id int) (*models.Session, error) {
	switch id {
	case 1:
		return mockSession, nil
	default:
		return nil, models.ErrNoRecord
	}
}
func (m *SessionModel) Latest() ([]*models.Session, error) {
	return []*models.Session{mockSession}, nil
}
