package match_repository

import (
	match_entity "cats-social/entity/match"
	exc "cats-social/exceptions"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type matchRepositoryImpl struct {
}

func NewMatchRepository() MatchRepository {
	return &matchRepositoryImpl{}
}

func (repository *matchRepositoryImpl) Create(ctx context.Context, tx pgx.Tx, match match_entity.Match, userId string) (match_entity.Match, error) {
	if err := checkCatExists(ctx, tx, match.CatIssuerId, match.CatReceiverId); err != nil {
		return match_entity.Match{}, err
	}
	if err := validateMatchCatCriteria(ctx, tx, match.CatIssuerId, match.CatReceiverId, userId); err != nil {
		return match_entity.Match{}, err
	}

	var matchId string
	var createdAt time.Time
	query := `INSERT INTO matches (id, message, cat_issuer_id, cat_receiver_id)
	SELECT 
		gen_random_uuid(), $1, $2, $3
	WHERE EXISTS (
		SELECT 1 FROM users WHERE id = $4
	)
	RETURNING id, created_at;
	`
	if err := tx.QueryRow(ctx, query, match.Message, match.CatIssuerId, match.CatReceiverId, string(userId)).Scan(&matchId, &createdAt); err != nil {
		tx.Rollback(ctx)
		if err == pgx.ErrNoRows {
			return match_entity.Match{}, exc.BadRequestException("Invalid user id")
		}
		return match_entity.Match{}, err
	}

	match.Id = matchId
	match.CreatedAt = createdAt.Format(time.RFC3339)
	if err := tx.Commit(ctx); err != nil {
		return match_entity.Match{}, err
	}
	return match, nil
}

func (repository *matchRepositoryImpl) Approve(ctx context.Context, tx pgx.Tx, match match_entity.Match, userId string) error {
	// check match id is exist
	var catIssuerId, catReceiverId, status string
	query := "SELECT cat_issuer_id, cat_receiver_id, status FROM matches WHERE id = $1 LIMIT 1"
	if err := tx.QueryRow(ctx, query, string(match.Id)).Scan(&catIssuerId, &catReceiverId, &status); err != nil {
		if err == pgx.ErrNoRows {
			return exc.NotFoundException("Match id is not found")
		}
		return exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}
	// check match id is valid or not based on status
	if status != "requested" {
		return exc.BadRequestException("Match id is no longer valid")
	}

	// check owner cat id
	var ownerReceiverId string
	checkOwnerReceiverCatQ := `SELECT user_id FROM cats WHERE id = $1`
	if err := tx.QueryRow(ctx, checkOwnerReceiverCatQ, string(catReceiverId)).Scan(&ownerReceiverId); err != nil {
		exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}
	if ownerReceiverId != userId {
		return exc.UnauthorizedException("You cannot approved that cat you are not belong")
	}

	// update status match to approved
	approveQuery := `UPDATE matches SET status = $1 WHERE id = $2`
	if _, err := tx.Exec(ctx, approveQuery, "approved", string(match.Id)); err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error when update match: %s", err))
	}

	// update has_matched both cats to true
	updateCatQuery := `UPDATE cats SET has_matched = $1 WHERE id IN ($2, $3)`
	if _, err := tx.Exec(ctx, updateCatQuery, true, catIssuerId, catReceiverId); err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error when update cat: %s", err))
	}

	// search remaining other match both cats
	remainMatchCatIssuerIds, err := getRemainingMatchCat(ctx, tx, catIssuerId)
	if err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error when get remain match issuer: %s", err))
	}
	remainMatchCatReceiverIds, err := getRemainingMatchCat(ctx, tx, catIssuerId)
	if err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error when get remain match receiver: %s", err))
	}
	remainMatchIds := append(remainMatchCatIssuerIds, remainMatchCatReceiverIds...)

	// delete if any remain match requested on both cats
	if len(remainMatchIds) > 0 {
		if err := deleteRemainingMatchCat(ctx, tx, remainMatchIds); err != nil {
			return exc.InternalServerException(fmt.Sprintf("Internal server error when deleting remain match: %s", err))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}

	return nil

}

func (repository *matchRepositoryImpl) Get(ctx context.Context, tx pgx.Tx, userId string) ([]match_entity.MatchGetDataResponse, error) {
	query := `SELECT m.id, 
	json_build_object('name', u.name, 'email', u.email, 'createdAt', to_char(u.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')) issuedBy, 
	json_build_object('id', c1.id, 'name', c1.name, 'race', c1.race, 'sex', c1.sex, 'ageInMonth', c1.age_in_month, 'imageUrls', c1.image_urls, 'description', c1.description, 'hasMatched', c1.has_matched, 'createdAt', to_char(c1.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')) matchCatDetail,
	json_build_object('id', c2.id, 'name', c2.name, 'race', c2.race, 'sex', c2.sex, 'ageInMonth', c2.age_in_month, 'imageUrls', c2.image_urls, 'description', c2.description, 'hasMatched', c2.has_matched, 'createdAt', to_char(c2.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')) userCatDetail,
	m.message,
	to_char(m.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') createdAt
	FROM matches m
		JOIN cats c1 ON c1.id = m.cat_receiver_id
		JOIN cats c2 ON c2.id = m.cat_issuer_id
		JOIN users u ON u.id = c2.user_id
	WHERE c1.user_id = $1 OR c2.user_id = $1
	ORDER BY m.created_at DESC
	`
	rows, err := tx.Query(ctx, query, string(userId))
	if err != nil {
		return []match_entity.MatchGetDataResponse{}, err
	}
	defer rows.Close()

	matches, err := pgx.CollectRows(rows, pgx.RowToStructByName[match_entity.MatchGetDataResponse])

	if err != nil {
		return []match_entity.MatchGetDataResponse{}, err
	}

	return matches, nil
}

/********************* HELPER METHODS *******************************/

func deleteRemainingMatchCat(ctx context.Context, tx pgx.Tx, matchIds []string) error {
	placeholders := make([]string, len(matchIds))
	values := make([]interface{}, len(matchIds))
	for i, matchId := range matchIds {
		placeholders[i] = "$" + strconv.Itoa(i+1)
		values[i] = matchId
	}

	query := fmt.Sprintf(`DELETE FROM matches WHERE id IN (%s)`, strings.Join(placeholders, ", "))
	if _, err := tx.Exec(ctx, query, values...); err != nil {
		return err
	}
	return nil
}

func getRemainingMatchCat(ctx context.Context, tx pgx.Tx, catId string) ([]string, error) {
	var matchCatIds []string
	remainingMatchCatsQ := `SELECT id FROM matches WHERE (cat_issuer_id = $1 OR cat_receiver_id = $1) AND status = 'requested'`
	rows, err := tx.Query(ctx, remainingMatchCatsQ, string(catId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		matchCatIds = append(matchCatIds, id)
	}

	return matchCatIds, nil
}

func checkCatExists(ctx context.Context, tx pgx.Tx, catIssuerId string, catReceiverId string) error {
	query := `SELECT id FROM cats WHERE id = $1`
	catIssuer, err := tx.Exec(ctx, query, string(catIssuerId))
	if err != nil {
		if err == pgx.ErrNoRows {
			return exc.NotFoundException("Cat Issuer does not exist!")
		} else {
			return exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
		}
	}
	if catIssuer.RowsAffected() == 0 {
		return exc.NotFoundException("Cat Issuer does not exist!")
	}

	catReceiver, err := tx.Exec(ctx, query, string(catReceiverId))
	if err != nil {
		if err == pgx.ErrNoRows {
			return exc.NotFoundException("Cat Receiver does not exist!")
		} else {
			return exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
		}
	}
	if catReceiver.RowsAffected() == 0 {
		return exc.NotFoundException("Cat Receiver does not exist!")
	}

	return nil
}

func validateMatchCatCriteria(ctx context.Context, tx pgx.Tx, catIssuerId string, catReceiverId string, userId string) error {
	// check match request already exist or not
	checkRequestMatchQ := `SELECT EXISTS (SELECT 1 FROM matches m WHERE m.cat_issuer_id = $1 AND m.cat_receiver_id = $2 AND status = $3)`
	var isAlreadyRequestMatch bool
	if err := tx.QueryRow(ctx, checkRequestMatchQ, string(catIssuerId), string(catReceiverId), "requested").Scan(&isAlreadyRequestMatch); err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}
	if isAlreadyRequestMatch {
		return exc.ConflictException("Your cat has already request match to this cat, please waiting for response from receiver!")
	}

	query := `SELECT sex, has_matched, user_id FROM cats WHERE id = $1`
	var catIssuerSex, catIssuerUserId string
	var catIssuerHasMatched bool
	if err := tx.QueryRow(ctx, query, string(catIssuerId)).Scan(&catIssuerSex, &catIssuerHasMatched, &catIssuerUserId); err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}

	var catReceiverSex, catReceiverUserId string
	var catReceiverHasMatched bool
	if err := tx.QueryRow(ctx, query, string(catReceiverId)).Scan(&catReceiverSex, &catReceiverHasMatched, &catReceiverUserId); err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}

	switch {
	// check cat owner
	case catIssuerUserId != userId:
		return exc.UnauthorizedException("You cannot match a cat that you do not own!")
	// check cat issuer and receiver cannot from the same owner
	case catIssuerUserId == catReceiverUserId:
		return exc.BadRequestException("Match cannot be made from the same cat's owner!")
	// check cat alreay match or not
	case catIssuerHasMatched:
		return exc.BadRequestException("Cat's issuer already matched, match another one!")
	case catReceiverHasMatched:
		return exc.BadRequestException("Cat's receiver already matched, match another one!")
	// check cat gender cannot be same
	case catIssuerSex == catReceiverSex:
		return exc.BadRequestException("Cat's gender cannot be same")
	default:
		return nil
	}
}
