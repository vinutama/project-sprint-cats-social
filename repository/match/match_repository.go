package match_repository

import (
	match_entity "cats-social/entity/match"
	"context"

	"github.com/jackc/pgx/v5"
)

type MatchRepository interface {
	Create(ctx context.Context, tx pgx.Tx, req match_entity.Match, userId string) error
}
