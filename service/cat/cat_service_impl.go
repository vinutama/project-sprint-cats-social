package cat_service

import (
	cat_entity "cats-social/entity/cat"
	exc "cats-social/exceptions"
	catRep "cats-social/repository/cat"
	authService "cats-social/service/auth"
	"fmt"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatServiceImpl struct {
	CatRepository catRep.CatRepository
	DBPool        *pgxpool.Pool
	AuthService   authService.AuthService
	Validator     *validator.Validate
}

func NewUserService(catRepository catRep.CatRepository, dbPool *pgxpool.Pool, authService authService.AuthService, validator *validator.Validate) CatService {
	return &CatServiceImpl{
		CatRepository: catRepository,
		DBPool:        dbPool,
		AuthService:   authService,
		Validator:     validator,
	}
}

func (service *CatServiceImpl) Create(ctx *fiber.Ctx, req cat_entity.CatCreateRequest) (cat_entity.CatCreateResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return cat_entity.CatCreateResponse{}, exc.BadRequestException(fmt.Sprintf("%s", err))
	}

	tx, err := service.DBPool.Begin(ctx.UserContext())
	if err != nil {
		return cat_entity.CatCreateResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}
	defer tx.Rollback(ctx.UserContext())

	userId, err := authService.NewAuthService().GetValidUser(ctx)
	if err != nil {
		return cat_entity.CatCreateResponse{}, exc.UnauthorizedException("Unauthorized")
	}
	cat := cat_entity.Cat{
		Name:        req.Name,
		Race:        req.Race,
		Sex:         req.Sex,
		AgeInMonth:  req.AgeInMonth,
		Description: req.Description,
		ImageURLs:   req.ImageURLs,
	}

	catRegistered, err := catRep.NewCatRepository().Create(ctx.UserContext(), tx, cat, userId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return cat_entity.CatCreateResponse{}, exc.BadRequestException("Invalid user id")
		}
		return cat_entity.CatCreateResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}

	return cat_entity.CatCreateResponse{
		Message: "success",
		Data: &cat_entity.CatCreateDataResponse{
			Id:        catRegistered.Id,
			CreatedAt: catRegistered.CreatedAt,
		},
	}, nil
}
