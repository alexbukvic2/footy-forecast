package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus is the lifecycle state of a user account.
type UserStatus string

// User status values. These mirror the Postgres enum user_status.
const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
)

// Valid reports whether s is a defined status value.
func (s UserStatus) Valid() bool {
	switch s {
	case UserStatusActive, UserStatusSuspended:
		return true
	}
	return false
}

// User represents an authenticated player in the system.
type User struct {
	ID          uuid.UUID
	CognitoSub  string
	Email       string
	DisplayName string
	Status      UserStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
