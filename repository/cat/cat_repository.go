package cat_repository

import (
	cat_entity "cats-social/entity/cat"
	"context"

	"github.com/jackc/pgx/v5"
)

type CatRepository interface {
	Create(ctx context.Context, tx pgx.Tx, req cat_entity.Cat, ownerId string) (cat_entity.Cat, error)
}