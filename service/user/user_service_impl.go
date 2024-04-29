package user_service

import (
	user_entity "cats-social/entity/user"
	exc "cats-social/exceptions"
	helpers "cats-social/helpers"
	userRep "cats-social/repository/user"
	authService "cats-social/service/auth"
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserRepository userRep.UserRepository
	DBPool         *pgxpool.Pool
	AuthService    authService.AuthService
	Validator      *validator.Validate
}

func NewUserService(userRepository userRep.UserRepository, dbPool *pgxpool.Pool, authService authService.AuthService, validator *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DBPool:         dbPool,
		AuthService:    authService,
		Validator:      validator,
	}
}

func (service *UserServiceImpl) Register(ctx context.Context, req user_entity.UserRegisterRequest) (user_entity.UserRegisterResponse, error) {
	// validate by rule we defined in _request_entity.go
	if err := service.Validator.Struct(req); err != nil {
		return user_entity.UserRegisterResponse{}, exc.BadRequestException(fmt.Sprintf("Bad request: %s", err))
	}

	tx, err := service.DBPool.Begin(ctx)
	if err != nil {
		return user_entity.UserRegisterResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
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
		if strings.Contains(err.Error(), "duplicate key value") {
			return user_entity.UserRegisterResponse{}, exc.ConflictException("User with this email already registered")
		}
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
			return user_entity.UserLoginResponse{}, exc.NotFoundException("User is not found")
		}

		return user_entity.UserLoginResponse{}, err
	}

	if _, err = helpers.ComparePassword(userLogin.Password, req.Password); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return user_entity.UserLoginResponse{}, exc.BadRequestException("Invalid password")
		}

		return user_entity.UserLoginResponse{}, err
	}

	token, err := authService.NewAuthService().GenerateToken(ctx, userLogin.Id)
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
