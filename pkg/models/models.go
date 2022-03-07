package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// Error for when a user tries to login in with an incorrect email address or
	// password
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// Error for when a user tries to sign up with an email adress that is already
	// found in the database (not unique)
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

type Session struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}
