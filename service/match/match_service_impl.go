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

	matchRegistered, err := matchRep.NewMatchRepository().Create(userCtx, tx, match, userId)
	if err != nil {
		return match_entity.MatchCreateResponse{}, err
	}

	return match_entity.MatchCreateResponse{
		Message: "Match request success, waiting for response from receiver",
		Data: &match_entity.MatchCreateDataResponse{
			Id:        matchRegistered.Id,
			CreatedAt: matchRegistered.CreatedAt,
		},
	}, nil
}

func (service *matchServiceImpl) Approve(ctx *fiber.Ctx, req match_entity.MatchApproveRequest) (match_entity.MatchApproveResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return match_entity.MatchApproveResponse{}, exc.BadRequestException(fmt.Sprintf("%s", err))
	}

	userCtx := ctx.UserContext()
	tx, err := service.DBPool.Begin(userCtx)
	if err != nil {
		return match_entity.MatchApproveResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}
	defer tx.Rollback(userCtx)

	userId, err := authService.NewAuthService().GetValidUser(ctx)
	if err != nil {
		return match_entity.MatchApproveResponse{}, exc.UnauthorizedException("Unauthorized")
	}
	match := match_entity.Match{
		Id: req.MatchId,
	}

	if err := matchRep.NewMatchRepository().Approve(userCtx, tx, match, userId); err != nil {
		return match_entity.MatchApproveResponse{}, err
	}

	return match_entity.MatchApproveResponse{
		Message: "Congratulations your cat is matched!",
	}, nil
}

func (service *matchServiceImpl) Get(ctx *fiber.Ctx) (match_entity.MatchGetResponse, error) {
	userCtx := ctx.UserContext()
	tx, err := service.DBPool.Begin(userCtx)
	if err != nil {
		return match_entity.MatchGetResponse{}, exc.InternalServerException(fmt.Sprintf("Internal Server Error: %s", err))
	}
	defer tx.Rollback(userCtx)

	userId, err := authService.NewAuthService().GetValidUser(ctx)
	if err != nil {
		return match_entity.MatchGetResponse{}, exc.UnauthorizedException("Unauthorized")
	}

	data, err := matchRep.NewMatchRepository().Get(userCtx, tx, userId)
	if err != nil {
		return match_entity.MatchGetResponse{}, err
	}

	return match_entity.MatchGetResponse{
		Message: "Successfully get matches",
		Data:    &data,
	}, nil

}

func (service *matchServiceImpl) Delete(ctx *fiber.Ctx, params match_entity.MatchDeleteParams) (match_entity.MatchDeleteResponse, error) {
	if err := service.Validator.Struct(params); err != nil {
		return match_entity.MatchDeleteResponse{}, exc.BadRequestException(fmt.Sprintf("%s", err))
	}

	userCtx := ctx.Context()
	tx, err := service.DBPool.Begin(userCtx)
	if err != nil {
		return match_entity.MatchDeleteResponse{}, exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}
	defer tx.Rollback(userCtx)

	userId, err := authService.NewAuthService().GetValidUser(ctx)
	if err != nil {
		return match_entity.MatchDeleteResponse{}, exc.UnauthorizedException("Unauthorized")
	}
	match := match_entity.Match{
		Id: params.Id,
	}
	if err := matchRep.NewMatchRepository().Delete(userCtx, tx, match, userId); err != nil {
		return match_entity.MatchDeleteResponse{}, err
	}

	return match_entity.MatchDeleteResponse{
		Message: "Match request successfully deleted",
	}, nil
}
