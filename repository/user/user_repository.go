package user_repository

import (
	user_entity "cats-social/entity/user"
	"context"
)

type UserRepository interface {
	Register(ctx context.Context, req user_entity.User) (user_entity.User, error)
	Login(ctx context.Context, req user_entity.User) (user_entity.User, error)
}
