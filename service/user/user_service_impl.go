package user_service

import (
	user_entity "cats-social/entity/user"
	helpers "cats-social/helpers"
	userRep "cats-social/repository/user"
	authService "cats-social/service/auth"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
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

func (service *UserServiceImpl) Login(ctx context.Context, req user_entity.UserLoginRequest) (user_entity.UserLoginResponse, error) {
	tx, err := service.DBPool.Begin(ctx)
	if err != nil {
		return user_entity.UserLoginResponse{}, err
	}
	defer tx.Rollback(ctx)

	user := user_entity.User{
		Email: req.Email,
	}

	userLogin, err := userRep.NewUserRepository().Login(ctx, tx, user)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return user_entity.UserLoginResponse{
				Message: "User not found",
				Status:  404,
			}, nil
		}

		return user_entity.UserLoginResponse{}, err
	}

	if _, err = helpers.ComparePassword(userLogin.Password, req.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return user_entity.UserLoginResponse{
				Message: "Invalid password",
				Status:  400,
			}, nil
		}

		return user_entity.UserLoginResponse{}, err
	}

	token, err := authService.NewAuthService().GenerateToken(ctx, user.Id)
	if err != nil {
		return user_entity.UserLoginResponse{}, err
	}

	return user_entity.UserLoginResponse{
		Message: "User logged successfully",
		Data: &user_entity.UserData{
			Email:       userLogin.Email,
			Name:        userLogin.Name,
			AccessToken: token,
		},
	}, nil

}
