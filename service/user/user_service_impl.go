package user_service

import (
	user_entity "cats-social/entity/user"
	helpers "cats-social/helpers"
	userRep "cats-social/repository/user"
	authService "cats-social/service/auth"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserServiceImpl struct {
	UserRepository userRep.UserRepository
	DBPool         *pgxpool.Pool
	AuthService    authService.AuthService
}

func NewUserService(userRepository userRep.UserRepository, dbPool *pgxpool.Pool, authService authService.AuthService) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DBPool:         dbPool,
		AuthService:    authService,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, req user_entity.UserRegisterRequest) (user_entity.UserRegisterResponse, error) {
	tx, err := service.DBPool.Begin(ctx)
	if err != nil {
		return user_entity.UserRegisterResponse{}, err
	}
	defer tx.Rollback(ctx)

	hashPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return user_entity.UserRegisterResponse{}, err
	}

	user := user_entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashPassword,
	}
	userRegistered, err := userRep.NewUserRepository().Register(ctx, tx, user)
	if err != nil {
		return user_entity.UserRegisterResponse{}, err
	}

	token, err := authService.NewAuthService().GenerateToken(ctx, userRegistered.Id)
	if err != nil {
		return user_entity.UserRegisterResponse{}, err
	}

	return user_entity.UserRegisterResponse{
		Message: "User registered",
		Data: &user_entity.UserData{
			Email:       userRegistered.Email,
			Name:        userRegistered.Name,
			AccessToken: token,
		},
	}, nil
}
