package cat_repository

import (
	cat_entity "cats-social/entity/cat"
	"context"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
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
	if err := tx.QueryRow(ctx, query, cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, cat.Description, cat.ImageURLs, ownerId).Scan(&catId, &catCreatedAt); err != nil {
		tx.Rollback(ctx)
		return cat_entity.Cat{}, err
	}

	cat.Id = catId
	cat.IsDeleted = false
	cat.HasMatched = false
	cat.CreatedAt = catCreatedAt

	if err := tx.Commit(ctx); err != nil {
		return cat_entity.Cat{}, err
	}

	return cat, nil
}

func (repository *CatRepositoryImpl) Search(ctx context.Context, tx pgx.Tx, searchQuery cat_entity.CatSearch) ([]cat_entity.Cat, error) {
	var cats = []cat_entity.Cat{}
	query := `SELECT id, name, race, sex, age_in_month, image_urls, description, has_matched, created_at FROM cats
	WHERE
		($1 = '' OR id = $1) AND
		($2 = '' OR race = $2) AND
		($3 = '' OR sex = $3) AND
		($4 = '' OR name = $4)`

	if searchQuery.HasMatched != "" {
		if searchQuery.HasMatched == "true" || searchQuery.HasMatched == "false" {
			query += " AND has_matched = $5"
		}

	} else {
		query += " AND ($5 = '')"
	}
	if searchQuery.AgeInMonth != 0 {
		query += fmt.Sprintf(" AND age_in_month %s $6", searchQuery.AgeCondition)
	} else {
		query += " AND age_in_month != $6"
	}
	if searchQuery.Owned != "" {
		if searchQuery.Owned == "true" {
			query += " AND user_id = $7"
		} else {
			query += " AND user_id != $7"

		}
	} else {
		query += " AND user_id != $7"
	}
	query += " LIMIT $8 OFFSET $9;"

	err := pgxscan.Select(ctx, tx, &cats, query, searchQuery.Id, searchQuery.Race, searchQuery.Sex, searchQuery.Name, searchQuery.HasMatched, searchQuery.AgeInMonth, searchQuery.UserId, searchQuery.Limit, searchQuery.Offset)
	if err != nil {
		return cats, err
	}

	return cats, nil
}
