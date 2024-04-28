package user_service

import (
	user_entity "cats-social/entity/user"
	"context"
)

type UserService interface {
	Register(ctx context.Context, req user_entity.UserRegisterRequest) (user_entity.UserRegisterResponse, error)
}
