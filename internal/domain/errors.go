// Package domain contains core types and sentinel errors shared across layers.
package domain

import "errors"

// ErrNotFound is returned when a requested resource does not exist.
var ErrNotFound = errors.New("not found")

// ErrConflict is returned when a write violates a uniqueness constraint
// or other business invariant (e.g., duplicate slug).
var ErrConflict = errors.New("conflict")

// ErrInvalid is returned when input fails validation at the service layer.
var ErrInvalid = errors.New("invalid")
