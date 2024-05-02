package controller

import (
	match_entity "cats-social/entity/match"
	exc "cats-social/exceptions"
	auth_service "cats-social/service/auth"
	match_service "cats-social/service/match"

	"github.com/gofiber/fiber/v2"
)

type MatchController struct {
	MatchService match_service.MatchService
	AuthService  auth_service.AuthService
}

func NewMatchController(matchService match_service.MatchService, authService auth_service.AuthService) MatchController {
	return MatchController{
		MatchService: matchService,
		AuthService:  authService,
	}
}

func (controller *MatchController) Create(ctx *fiber.Ctx) error {
	matchReq := new(match_entity.MatchCreateRequest)
	if err := ctx.BodyParser(matchReq); err != nil {
		return err
	}

	resp, err := controller.MatchService.Create(ctx, *matchReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (controller *MatchController) Delete(ctx *fiber.Ctx) error {
	matchParams := new(match_entity.MatchDeleteParams)
	if err := ctx.ParamsParser(matchParams); err != nil {
		return err
	}

	resp, err := controller.MatchService.Delete(ctx, *matchParams)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
