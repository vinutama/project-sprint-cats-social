package cat_repository

import (
	cat_entity "cats-social/entity/cat"
	"context"
)

type CatRepository interface {
	Create(ctx context.Context, req cat_entity.Cat, ownerId string) (cat_entity.Cat, error)
	Edit(ctx context.Context, req cat_entity.Cat, ownerId string, catId string) (cat_entity.Cat, error)
	Search(ctx context.Context, req cat_entity.CatSearch) ([]cat_entity.Cat, error)
	Delete(ctx context.Context, catId string, ownerId string) (cat_entity.Cat, error)
}
