package match_repository

import (
	match_entity "cats-social/entity/match"
	"context"
)

type MatchRepository interface {
	Create(ctx context.Context, req match_entity.Match, userId string) (match_entity.Match, error)
	Approve(ctx context.Context, req match_entity.Match, userId string) error
	Reject(ctx context.Context, req match_entity.Match, userId string) error
	Get(ctx context.Context, userId string) ([]match_entity.MatchGetDataResponse, error)
	Delete(ctx context.Context, req match_entity.Match, userId string) error
}
