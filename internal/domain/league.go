package domain

import (
	"time"

	"github.com/google/uuid"
)

// LeagueMemberRole distinguishes league owners from regular members.
type LeagueMemberRole string

// LeagueMemberRoleOwner and LeagueMemberRoleMember are the valid roles.
const (
	LeagueMemberRoleOwner  LeagueMemberRole = "owner"
	LeagueMemberRoleMember LeagueMemberRole = "member"
)

// League is a private competition group scoped to a tournament.
type League struct {
	ID           uuid.UUID
	TournamentID uuid.UUID
	OwnerID      uuid.UUID
	Name         string
	Code         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// LeagueMember represents a user's membership in a league.
type LeagueMember struct {
	LeagueID uuid.UUID
	UserID   uuid.UUID
	Role     LeagueMemberRole
	JoinedAt time.Time
}

// CreateLeagueInput carries the caller-supplied fields for league creation.
type CreateLeagueInput struct {
	TournamentID uuid.UUID
	Name         string
}
