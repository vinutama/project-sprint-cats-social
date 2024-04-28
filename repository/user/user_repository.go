package user_repository

import (
	user_entity "cats-social/entity/user"
	"context"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	Register(ctx context.Context, tx pgx.Tx, req user_entity.User) (user_entity.User, error)
	// Login(ctx context.Context, tx pgx.Tx, req user_entity.User) (user_entity.User, error)
}
