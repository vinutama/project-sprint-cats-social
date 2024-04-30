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
		return err
	}

	resp, err := controller.CatService.Create(ctx, *catReq)
	if err != nil {
		return exc.Exception(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(resp)
}