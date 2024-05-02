package match_service

import (
	match_entity "cats-social/entity/match"

	"github.com/gofiber/fiber/v2"
)

type MatchService interface {
	Create(ctx *fiber.Ctx, req match_entity.MatchCreateRequest) (match_entity.MatchCreateResponse, error)
	Get(ctx *fiber.Ctx) (match_entity.MatchGetResponse, error)
}
