package cat_service

import (
	cat_entity "cats-social/entity/cat"
	exc "cats-social/exceptions"
	catRep "cats-social/repository/cat"
	authService "cats-social/service/auth"
	"fmt"
	"strconv"
	"strings"
	"time"

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

func NewCatService(catRepository catRep.CatRepository, dbPool *pgxpool.Pool, authService authService.AuthService, validator *validator.Validate) CatService {
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

	userCtx := ctx.UserContext()
	tx, err := service.DBPool.Begin(userCtx)
	if err != nil {
		return cat_entity.CatCreateResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}
	defer tx.Rollback(userCtx)

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

	catRegistered, err := catRep.NewCatRepository().Create(userCtx, tx, cat, userId)
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
			CreatedAt: catRegistered.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (service *CatServiceImpl) EditCat(ctx *fiber.Ctx, req cat_entity.CatEditRequest) (cat_entity.CatEditResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return cat_entity.CatEditResponse{}, exc.BadRequestException(fmt.Sprintf("%s", err))
	}

	userCtx := ctx.UserContext()
	catId := ctx.Params("id")
	tx, err := service.DBPool.Begin(userCtx)
	if err != nil {
		return cat_entity.CatEditResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}
	defer tx.Rollback(userCtx)

	userId, err := authService.NewAuthService().GetValidUser(ctx)
	if err != nil {
		return cat_entity.CatEditResponse{}, exc.UnauthorizedException("Unauthorized")
	}
	cat := cat_entity.Cat{
		Name:        req.Name,
		Race:        req.Race,
		Sex:         req.Sex,
		AgeInMonth:  req.AgeInMonth,
		Description: req.Description,
		ImageURLs:   req.ImageURLs,
	}

	editedcat, err := catRep.NewCatRepository().Edit(userCtx, tx, cat, catId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return cat_entity.CatEditResponse{}, exc.NotFoundException("User/Cat id is not found/match")
		}
		return cat_entity.CatEditResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}

	return cat_entity.CatEditResponse{
		Message: "success",
		Data: &cat_entity.CatEditDataResponse{
			Id: editedcat.Id,
		},
	}, nil
}
func (service *CatServiceImpl) Search(ctx *fiber.Ctx, searchQueries cat_entity.CatSearchQuery) (cat_entity.CatSearchResponse, error) {
	if err := service.Validator.Struct(searchQueries); err != nil {
		return cat_entity.CatSearchResponse{}, exc.BadRequestException(fmt.Sprintf("%s", err))
	}

	userCtx := ctx.UserContext()
	tx, err := service.DBPool.Begin(ctx.UserContext())
	if err != nil {
		return cat_entity.CatSearchResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}
	defer tx.Rollback(ctx.UserContext())

	userId, err := authService.NewAuthService().GetValidUser(ctx)
	if err != nil {
		return cat_entity.CatSearchResponse{}, exc.UnauthorizedException("Unauthorized")
	}

	if strings.ToLower(searchQueries.HasMatched) != "true" || strings.ToLower(searchQueries.HasMatched) != "false" {
		searchQueries.HasMatched = ""
	}
	if strings.ToLower(searchQueries.Owned) != "true" || strings.ToLower(searchQueries.Owned) != "false" {
		searchQueries.Owned = ""
	}

	cat := cat_entity.CatSearch{
		Id:           searchQueries.Id,
		Race:         searchQueries.Race,
		Sex:          searchQueries.Sex,
		HasMatched:   searchQueries.HasMatched,
		Owned:        searchQueries.Owned,
		UserId:       userId,
		AgeCondition: "!=",
		Name:         searchQueries.Search,
		Limit:        5,
		Offset:       0,
	}

	if searchQueries.AgeInMonth != "" {
		if strings.Contains(searchQueries.AgeInMonth, ">") || strings.Contains(searchQueries.AgeInMonth, "<") || strings.Contains(searchQueries.AgeInMonth, "=") {
			age, _ := strconv.Atoi(searchQueries.AgeInMonth[1:len(searchQueries.AgeInMonth)])

			cat.AgeCondition = fmt.Sprintf("%c", searchQueries.AgeInMonth[0])
			cat.AgeInMonth = age
		} else {
			cat.AgeCondition = "="
			age, _ := strconv.Atoi(searchQueries.AgeInMonth)
			cat.AgeInMonth = age
		}
	}
	if searchQueries.Limit != "" {
		cat.Limit, _ = strconv.Atoi(searchQueries.Limit)
	}
	if searchQueries.Offset != "" {
		cat.Offset, _ = strconv.Atoi(searchQueries.Offset)
	}

	catSearched, err := catRep.NewCatRepository().Search(userCtx, tx, cat)
	if err != nil {
		return cat_entity.CatSearchResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}

	data := []cat_entity.CatSearchDataResponse{}
	for _, cat := range catSearched {
		catData := cat_entity.CatSearchDataResponse{
			Id:          cat.Id,
			Name:        cat.Name,
			Race:        cat.Race,
			Sex:         cat.Sex,
			AgeInMonth:  cat.AgeInMonth,
			ImageURLs:   cat.ImageURLs,
			Description: cat.Description,
			HasMatched:  cat.HasMatched,
			CreatedAt:   cat.CreatedAt.Format(time.RFC3339),
		}

		data = append(data, catData)
	}

	return cat_entity.CatSearchResponse{
		Messagge: "success",
		Data:     &data,
	}, nil
}

func (service *CatServiceImpl) Delete(ctx *fiber.Ctx) (cat_entity.CatDeleteResponse, error) {
	userCtx := ctx.UserContext()
	tx, err := service.DBPool.Begin(userCtx)
	if err != nil {
		return cat_entity.CatDeleteResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}

	defer tx.Rollback(ctx.UserContext())
	userId, err := authService.NewAuthService().GetValidUser(ctx)
	if err != nil {
		return cat_entity.CatDeleteResponse{}, exc.UnauthorizedException("Unauthorized")
	}

	catId := ctx.Params("id")
	deletedCat, err := catRep.NewCatRepository().Delete(userCtx, tx, catId, userId)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return cat_entity.CatDeleteResponse{}, exc.NotFoundException("Invalid user id")
		}
		return cat_entity.CatDeleteResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}

	return cat_entity.CatDeleteResponse{
		Message: "success",
		Data: &cat_entity.CatDeleteDataResponse{
			Id: deletedCat.Id,
		},
	}, nil

}
