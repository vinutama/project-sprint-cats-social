package cat_repository

import (
	cat_entity "cats-social/entity/cat"
	exc "cats-social/exceptions"
	"context"
	"fmt"
	"strconv"
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
	query := `SELECT * FROM cats WHERE is_deleted = false`
	var whereClause []string
	params := []interface{}{}

	if searchQuery.Id != "" {
		whereClause = append(whereClause, fmt.Sprintf("id = $%s", strconv.Itoa(len(params)+1)))
		params = append(params, searchQuery.Id)
	}

	if searchQuery.Race != "" {
		whereClause = append(whereClause, fmt.Sprintf("race = $%s", strconv.Itoa(len(params)+1)))
		params = append(params, searchQuery.Race)
	}

	if searchQuery.Sex != "" {
		whereClause = append(whereClause, fmt.Sprintf("sex = $%s", strconv.Itoa(len(params)+1)))
		params = append(params, searchQuery.Sex)
	}

	if searchQuery.Name != "" {
		whereClause = append(whereClause, fmt.Sprintf("name ILIKE '%%' || $%s || '%%'", strconv.Itoa(len(params)+1)))
		params = append(params, searchQuery.Name)
	}

	if searchQuery.HasMatched != "" {
		whereClause = append(whereClause, fmt.Sprintf("has_matched = $%s", strconv.Itoa(len(params)+1)))
		hasMatched, err := strconv.ParseBool(searchQuery.HasMatched)
		if err != nil {
			return []cat_entity.Cat{}, err
		}
		params = append(params, hasMatched)
	}

	if searchQuery.AgeInMonth > 0 {
		whereClause = append(whereClause, fmt.Sprintf("(CASE WHEN $%s >= 0 THEN age_in_month %s $%s ELSE TRUE END)", strconv.Itoa(len(params)+1), searchQuery.AgeCondition, strconv.Itoa(len(params)+1)))
		params = append(params, searchQuery.AgeInMonth)
	}

	if searchQuery.Owned != "" {
		whereClause = append(whereClause, fmt.Sprintf("(CASE WHEN $%s = 'true' THEN user_id = $%s WHEN $%s = 'false' THEN user_id != $%s ELSE TRUE END)", strconv.Itoa(len(params)+1), strconv.Itoa(len(params)+2), strconv.Itoa(len(params)+1), strconv.Itoa(len(params)+2)))
		params = append(params, searchQuery.Owned)
		params = append(params, searchQuery.UserId)
	}

	if len(whereClause) > 0 {
		query += " AND " + strings.Join(whereClause, " AND ")
	}

	query += " ORDER BY created_at DESC"

	if searchQuery.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%s OFFSET $%s", strconv.Itoa(len(params)+1), strconv.Itoa(len(params)+2))
		params = append(params, searchQuery.Limit)
		params = append(params, searchQuery.Offset)
	}

	rows, err := tx.Query(ctx, query, params...)
	if err != nil {
		return []cat_entity.Cat{}, err
	}
	defer rows.Close()

	cats, err := pgx.CollectRows(rows, pgx.RowToStructByName[cat_entity.Cat])
	if err != nil {
		return []cat_entity.Cat{}, err
	}

	return cats, nil
}
func (repository *CatRepositoryImpl) Edit(ctx context.Context, tx pgx.Tx, cat cat_entity.Cat, catId string, ownerId string) (cat_entity.Cat, error) {
	//  first we check whether the cat is already match(someone request to match) or not(no catfishing, pun intended)
	var matchCatSex string
	match_query := `Select cats.sex from matches join cats on cats.id=$1 where cat_issuer_id=$1 or cat_receiver_id=$1 and status!='rejected';`
	if err := tx.QueryRow(ctx, match_query, catId).Scan(&matchCatSex); err != nil {
		if err != pgx.ErrNoRows {
			return cat_entity.Cat{}, exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
		}
	}

	if matchCatSex != "" {
		return cat_entity.Cat{}, exc.BadRequestException("Cannot change gender when cat already match")
	}

	var catId2 string

	query := `Update cats set name=$1, race=$2, sex=$3,
	age_in_month=$4, description=$5, image_urls=$6 where id=$7 and user_id=$8
	returning id`
	if err := tx.QueryRow(ctx, query, cat.Name, cat.Race, cat.Sex, cat.AgeInMonth, cat.Description, cat.ImageURLs, catId, ownerId).Scan(&catId2); err != nil {
		tx.Rollback(ctx)
		return cat_entity.Cat{}, exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}

	cat.Id = catId
	cat.IsDeleted = false
	cat.HasMatched = false

	if err := tx.Commit(ctx); err != nil {
		return cat_entity.Cat{}, exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}

	return cat, nil
}

func (repository *CatRepositoryImpl) Delete(ctx context.Context, tx pgx.Tx, catId string, ownerId string) (cat_entity.Cat, error) {
	query := `update cats set is_deleted=true where id=$1 and user_id=$2 returning id`
	if err := tx.QueryRow(ctx, query, catId, ownerId).Scan(&catId); err != nil {
		return cat_entity.Cat{}, err
	}

	cat := cat_entity.Cat{}

	cat.Id = catId

	if err := tx.Commit(ctx); err != nil {
		return cat_entity.Cat{}, err
	}

	return cat, nil
}
