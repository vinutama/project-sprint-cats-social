package match_service

import (
	match_entity "cats-social/entity/match"
	exc "cats-social/exceptions"
	matchRep "cats-social/repository/match"
	authService "cats-social/service/auth"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type matchServiceImpl struct {
	MatchRepository matchRep.MatchRepository
	DBPool          *pgxpool.Pool
	AuthService     authService.AuthService
	Validator       *validator.Validate
}

func NewMatchService(matchRepository matchRep.MatchRepository, dbPool *pgxpool.Pool, authService authService.AuthService, validator *validator.Validate) MatchService {
	return &matchServiceImpl{
		MatchRepository: matchRepository,
		DBPool:          dbPool,
		AuthService:     authService,
		Validator:       validator,
	}
}

func (service *matchServiceImpl) Create(ctx *fiber.Ctx, req match_entity.MatchCreateRequest) (match_entity.MatchCreateResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return match_entity.MatchCreateResponse{}, exc.BadRequestException(fmt.Sprintf("%s", err))
	}

	userCtx := ctx.UserContext()
	tx, err := service.DBPool.Begin(userCtx)
	if err != nil {
		return match_entity.MatchCreateResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}
	defer tx.Rollback(userCtx)

	userId, err := authService.NewAuthService().GetValidUser(ctx)
	if err != nil {
		return match_entity.MatchCreateResponse{}, exc.UnauthorizedException("Unauthorized")
	}
	match := match_entity.Match{
		Message:       req.Message,
		CatIssuerId:   req.UserCatId,
		CatReceiverId: req.MatchCatId,
	}

	if err := matchRep.NewMatchRepository().Create(userCtx, tx, match, userId); err != nil {
		return match_entity.MatchCreateResponse{}, err
	}

	return match_entity.MatchCreateResponse{
		Message: "Cat Successfully matched",
	}, nil

}
