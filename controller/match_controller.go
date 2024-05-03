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
		return exc.BadRequestException("Failed to parse request body")
	}

	resp, err := controller.MatchService.Create(ctx, *matchReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (controller *MatchController) Approve(ctx *fiber.Ctx) error {
	matchApproveReq := new(match_entity.MatchActionRequest)
	if err := ctx.BodyParser(matchApproveReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}

	resp, err := controller.MatchService.Approve(ctx, *matchApproveReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (controller *MatchController) Reject(ctx *fiber.Ctx) error {
	matchRejectReq := new(match_entity.MatchActionRequest)
	if err := ctx.BodyParser(matchRejectReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}

	resp, err := controller.MatchService.Reject(ctx, *matchRejectReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (controller *MatchController) Get(ctx *fiber.Ctx) error {
	resp, err := controller.MatchService.Get(ctx)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
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
