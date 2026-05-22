// Package ctxutil provides typed context helpers for request-scoped values.
package ctxutil

import (
	"context"

	"github.com/alexbukvic2/footy-forecast/internal/domain"
)

type ctxKey string

const ctxKeyUser ctxKey = "user"

// WithUser returns a new context carrying u.
func WithUser(ctx context.Context, u domain.User) context.Context {
	return context.WithValue(ctx, ctxKeyUser, u)
}

// UserFromCtx extracts the User stored by WithUser.
// Returns (User, false) if no user is present — a programming error indicating
// the handler was called without the auth middleware.
func UserFromCtx(ctx context.Context) (domain.User, bool) {
	u, ok := ctx.Value(ctxKeyUser).(domain.User)
	return u, ok
}
