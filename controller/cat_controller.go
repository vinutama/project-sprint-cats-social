package controller

import (
	cat_entity "cats-social/entity/cat"
	exc "cats-social/exceptions"
	auth_service "cats-social/service/auth"
	cat_service "cats-social/service/cat"

	"github.com/gofiber/fiber/v2"
)

type CatController struct {
	CatService  cat_service.CatService
	AuthService auth_service.AuthService
}

func NewCatController(catService cat_service.CatService, authService auth_service.AuthService) CatController {
	return CatController{
		CatService:  catService,
		AuthService: authService,
	}
}

func (controller *CatController) Create(ctx *fiber.Ctx) error {
	catReq := new(cat_entity.CatCreateRequest)
	if err := ctx.BodyParser(catReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}

	resp, err := controller.CatService.Create(ctx, *catReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}

func (controller *CatController) EditCat(ctx *fiber.Ctx) error {
	catReq := new(cat_entity.CatEditRequest)
	if err := ctx.BodyParser(catReq); err != nil {
		return exc.BadRequestException("Failed to parse request body")
	}

	resp, err := controller.CatService.EditCat(ctx, *catReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (controller *CatController) Search(ctx *fiber.Ctx) error {
	catSearchQueries := new(cat_entity.CatSearchQuery)
	if err := ctx.QueryParser(catSearchQueries); err != nil {
		return err
	}

	resp, err := controller.CatService.Search(ctx, *catSearchQueries)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (controller *CatController) Delete(ctx *fiber.Ctx) error {
	resp, err := controller.CatService.Delete(ctx)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}
