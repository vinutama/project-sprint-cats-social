package cat_service

import (
	cat_entity "cats-social/entity/cat"

	"github.com/gofiber/fiber/v2"
)

type CatService interface {
	Create(ctx *fiber.Ctx, req cat_entity.CatCreateRequest) (cat_entity.CatCreateResponse, error)
	Search(ctx *fiber.Ctx, queries cat_entity.CatSearchQueries) (cat_entity.CatSearchResponse, error)
}
