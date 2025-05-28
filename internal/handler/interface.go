package handler

import (
	"context"
	"workout-tracker/internal/model/user"
	"workout-tracker/internal/service/auth"
)

type AuthService interface {
	GetUserByUserID(ctx context.Context, id int) (*user.User, error)
}

var _ AuthService = (*auth.AuthService)(nil)
