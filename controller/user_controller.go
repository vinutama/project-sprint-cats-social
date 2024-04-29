package controller

import (
	user_entity "cats-social/entity/user"
	auth_service "cats-social/service/auth"
	user_service "cats-social/service/user"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService user_service.UserService
	AuthService auth_service.AuthService
}

func NewUserController(userService user_service.UserService, authService auth_service.AuthService) UserController {
	return UserController{
		UserService: userService,
		AuthService: authService,
	}
}

func (controller *UserController) Register(ctx *fiber.Ctx) error {
	userReq := new(user_entity.UserRegisterRequest)
	if err := ctx.BodyParser(userReq); err != nil {
		return err
	}
	resp, err := controller.UserService.Register(ctx.UserContext(), *userReq)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(resp)

}

func (controller *UserController) Login(ctx *fiber.Ctx) error {
	userReq := new(user_entity.UserLoginRequest)
	if err := ctx.BodyParser(userReq); err != nil {
		return err
	}

	resp, err := controller.UserService.Login(ctx.UserContext(), *userReq)
	if err != nil {
		return err
	}

	return ctx.Status(resp.Status).JSON(resp)
}
