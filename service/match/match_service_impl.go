package match_service

import (
	match_entity "cats-social/entity/match"
	exc "cats-social/exceptions"
	matchRep "cats-social/repository/match"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

type matchServiceImpl struct {
	MatchRepository matchRep.MatchRepository
	Validator       *validator.Validate
}

func NewMatchService(matchRepository matchRep.MatchRepository, validator *validator.Validate) MatchService {
	return &matchServiceImpl{
		MatchRepository: matchRepository,
		Validator:       validator,
	}
}

func (service *matchServiceImpl) Create(ctx *fiber.Ctx, req match_entity.MatchCreateRequest) (match_entity.MatchCreateResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return match_entity.MatchCreateResponse{}, exc.BadRequestException(fmt.Sprintf("%s", err))
	}

	userCtx := ctx.UserContext()
	userId := ctx.Locals("userId").(string)
	match := match_entity.Match{
		Message:       req.Message,
		CatIssuerId:   req.UserCatId,
		CatReceiverId: req.MatchCatId,
	}

	matchRegistered, err := service.MatchRepository.Create(userCtx, match, userId)
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

func (service *matchServiceImpl) Approve(ctx *fiber.Ctx, req match_entity.MatchActionRequest) (match_entity.MatchActionResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return match_entity.MatchActionResponse{}, exc.BadRequestException(fmt.Sprintf("%s", err))
	}

	userCtx := ctx.UserContext()
	userId := ctx.Locals("userId").(string)
	match := match_entity.Match{
		Id: req.MatchId,
	}

	if err := service.MatchRepository.Approve(userCtx, match, userId); err != nil {
		return match_entity.MatchActionResponse{}, err
	}

	return match_entity.MatchActionResponse{
		Message: "Congratulations your cat is matched!",
	}, nil
}

func (service *matchServiceImpl) Reject(ctx *fiber.Ctx, req match_entity.MatchActionRequest) (match_entity.MatchActionResponse, error) {
	if err := service.Validator.Struct(req); err != nil {
		return match_entity.MatchActionResponse{}, exc.BadRequestException(fmt.Sprintf("%s", err))
	}

	userCtx := ctx.UserContext()
	userId := ctx.Locals("userId").(string)
	match := match_entity.Match{
		Id: req.MatchId,
	}

	if err := service.MatchRepository.Reject(userCtx, match, userId); err != nil {
		return match_entity.MatchActionResponse{}, err
	}

	return match_entity.MatchActionResponse{
		Message: "Successfully rejected the match request",
	}, nil
}

func (service *matchServiceImpl) Get(ctx *fiber.Ctx) (match_entity.MatchGetResponse, error) {
	userCtx := ctx.UserContext()
	userId := ctx.Locals("userId").(string)

	data, err := service.MatchRepository.Get(userCtx, userId)
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
	userId := ctx.Locals("userId").(string)
	match := match_entity.Match{
		Id: params.Id,
	}
	if err := service.MatchRepository.Delete(userCtx, match, userId); err != nil {
		return match_entity.MatchDeleteResponse{}, err
	}

	return match_entity.MatchDeleteResponse{
		Message: "Match request successfully deleted",
	}, nil
}
