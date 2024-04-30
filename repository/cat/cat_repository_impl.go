package cat_repository

import (
	cat_entity "cats-social/entity/cat"
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type CatRepositoryImpl struct {
}

func NewCatRepository() CatRepository {
	return &CatRepositoryImpl{}
}

func (repository *CatRepositoryImpl) Create(ctx context.Context, tx pgx.Tx, cat cat_entity.Cat, ownerId string) (cat_entity.Cat, error) {
	var catId string
	var catCreatedAt time.Time

	query := `INSERT INTO cats (id, name, race, sex, age_in_month, description, image_urls, user_id) 
	SELECT 
		gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7
	WHERE EXISTS (
		SELECT 1 FROM users WHERE id = CAST($7 as VARCHAR)
	)
	RETURNING id, created_at;`
	if err := tx.QueryRow(ctx, query, cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, cat.Description, strings.Join(cat.ImageURLs, "||"), ownerId).Scan(&catId, &catCreatedAt); err != nil {
		tx.Rollback(ctx)
		return cat_entity.Cat{}, err
	}

	cat.Id = catId
	cat.IsDeleted = false
	cat.HasMatched = false
	cat.CreatedAt = catCreatedAt.Format(time.RFC3339)

	if err := tx.Commit(ctx); err != nil {
		return cat_entity.Cat{}, err
	}

	return cat, nil
}
