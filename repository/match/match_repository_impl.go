package match_repository

import (
	match_entity "cats-social/entity/match"
	exc "cats-social/exceptions"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type matchRepositoryImpl struct {
	DBPool *pgxpool.Pool
}

func NewMatchRepository(dbPool *pgxpool.Pool) MatchRepository {
	return &matchRepositoryImpl{
		DBPool: dbPool,
	}
}

func (repository *matchRepositoryImpl) Create(ctx context.Context, match match_entity.Match, userId string) (match_entity.Match, error) {
	if err := repository.checkCatExists(ctx, match.CatIssuerId, match.CatReceiverId); err != nil {
		return match_entity.Match{}, err
	}
	if err := repository.validateMatchCatCriteria(ctx, match.CatIssuerId, match.CatReceiverId, userId); err != nil {
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
	if err := repository.DBPool.QueryRow(ctx, query, match.Message, match.CatIssuerId, match.CatReceiverId, string(userId)).Scan(&matchId, &createdAt); err != nil {
		if err == pgx.ErrNoRows {
			return match_entity.Match{}, exc.BadRequestException("Invalid user id")
		}
		return match_entity.Match{}, err
	}

	match.Id = matchId
	match.CreatedAt = createdAt.Format(time.RFC3339)
	return match, nil
}

func (repository *matchRepositoryImpl) Get(ctx context.Context, userId string) ([]match_entity.MatchGetDataResponse, error) {
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
		AND m.status = 'requested'
	ORDER BY m.created_at DESC
	`
	rows, err := repository.DBPool.Query(ctx, query, string(userId))
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

func (repository *matchRepositoryImpl) Delete(ctx context.Context, match match_entity.Match, userId string) error {
	err := repository.checkMatchDeletionEligibility(ctx, match.Id, userId)
	if err != nil {
		return err
	}

	query := `DELETE FROM matches WHERE id = $1`
	if _, err = repository.DBPool.Exec(ctx, query, match.Id); err != nil {
		return err
	}
	return nil
}

func (repository *matchRepositoryImpl) Approve(ctx context.Context, match match_entity.Match, userId string) error {
	// check match id is exist
	var catIssuerId, catReceiverId, status string
	query := "SELECT cat_issuer_id, cat_receiver_id, status FROM matches WHERE id = $1 LIMIT 1"
	if err := repository.DBPool.QueryRow(ctx, query, string(match.Id)).Scan(&catIssuerId, &catReceiverId, &status); err != nil {
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
	if err := repository.DBPool.QueryRow(ctx, checkOwnerReceiverCatQ, string(catReceiverId)).Scan(&ownerReceiverId); err != nil {
		exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}
	if ownerReceiverId != userId {
		return exc.UnauthorizedException("You cannot approved that cat you are not belong")
	}

	// update status match to approved
	approveQuery := `UPDATE matches SET status = $1 WHERE id = $2`
	if _, err := repository.DBPool.Exec(ctx, approveQuery, "approved", string(match.Id)); err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error when update match: %s", err))
	}

	// update has_matched both cats to true
	updateCatQuery := `UPDATE cats SET has_matched = $1 WHERE id IN ($2, $3)`
	if _, err := repository.DBPool.Exec(ctx, updateCatQuery, true, catIssuerId, catReceiverId); err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error when update cat: %s", err))
	}

	// delete if any remain match requested on both cats
	if err := repository.deleteRemainingMatchCat(ctx, catIssuerId, catReceiverId); err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error when deleting remain match: %s", err))
	}
	return nil

}

func (repository *matchRepositoryImpl) Reject(ctx context.Context, match match_entity.Match, userId string) error {
	// check match id is exist
	var catIssuerId, catReceiverId, status, ownerReceiverId string
	query := "SELECT cat_issuer_id, cat_receiver_id, status, c.user_id ownerReceiverId FROM matches m JOIN cats c ON c.id = m.cat_receiver_id WHERE m.id = $1 LIMIT 1"
	if err := repository.DBPool.QueryRow(ctx, query, string(match.Id)).Scan(&catIssuerId, &catReceiverId, &status, &ownerReceiverId); err != nil {
		if err == pgx.ErrNoRows {
			return exc.NotFoundException("Match id is not found")
		}
		return exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}
	// check match id is valid or not based on status
	if status != "requested" {
		return exc.BadRequestException("Match id is no longer valid")
	}

	if ownerReceiverId != userId {
		return exc.UnauthorizedException("You cannot reject that cat you are not belong")
	}

	// update status match to rejected
	rejectQuery := `UPDATE matches SET status = $1 WHERE id = $2`
	if _, err := repository.DBPool.Exec(ctx, rejectQuery, "rejected", string(match.Id)); err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error when update match: %s", err))
	}
	return nil

}

/********************* HELPER METHODS *******************************/

func (repository *matchRepositoryImpl) checkMatchDeletionEligibility(ctx context.Context, matchId string, userId string) error {
	var status string
	var matchIssuerId string
	query := `SELECT m.status, c.user_id FROM matches m 
	JOIN cats c ON m.cat_issuer_id = c.id 
	WHERE m.id = $1 and c.user_id = $2;`
	err := repository.DBPool.QueryRow(ctx, query, matchId, userId).Scan(&status, &matchIssuerId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return exc.NotFoundException("MatchId not found")
		} else {
			return exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
		}
	}

	// check match status
	if status != "requested" {
		return exc.BadRequestException("matchId is already approved / reject")
	}

	return nil
}

func (repository *matchRepositoryImpl) deleteRemainingMatchCat(ctx context.Context, catIssuerId string, catReceiverId string) error {
	query := `DELETE FROM matches
		WHERE (
			cat_issuer_id = $1 OR cat_receiver_id = $2 OR
			cat_issuer_id = $2 OR cat_receiver_id = $1
		) AND status = 'requested'
	`
	if _, err := repository.DBPool.Exec(ctx, query, catIssuerId, catReceiverId); err != nil {
		return err
	}
	return nil
}

func (repository *matchRepositoryImpl) checkCatExists(ctx context.Context, catIssuerId string, catReceiverId string) error {
	query := `SELECT id FROM cats WHERE id = $1`
	catIssuer, err := repository.DBPool.Exec(ctx, query, string(catIssuerId))
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

	catReceiver, err := repository.DBPool.Exec(ctx, query, string(catReceiverId))
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

func (repository *matchRepositoryImpl) validateMatchCatCriteria(ctx context.Context, catIssuerId string, catReceiverId string, userId string) error {
	// check match request already exist or not
	checkRequestMatchQ := `SELECT EXISTS (SELECT 1 FROM matches m 
			WHERE (m.cat_issuer_id = $1 AND m.cat_receiver_id = $2 OR m.cat_issuer_id = $2 AND m.cat_receiver_id = $1)
			AND status = $3
		)`
	var isAlreadyRequestMatch bool
	if err := repository.DBPool.QueryRow(ctx, checkRequestMatchQ, string(catIssuerId), string(catReceiverId), "requested").Scan(&isAlreadyRequestMatch); err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}
	if isAlreadyRequestMatch {
		return exc.BadRequestException("Your cat has already request match to this cat")
	}

	query := `SELECT sex, has_matched, user_id FROM cats WHERE id = $1`
	var catIssuerSex, catIssuerUserId string
	var catIssuerHasMatched bool
	if err := repository.DBPool.QueryRow(ctx, query, string(catIssuerId)).Scan(&catIssuerSex, &catIssuerHasMatched, &catIssuerUserId); err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}

	var catReceiverSex, catReceiverUserId string
	var catReceiverHasMatched bool
	if err := repository.DBPool.QueryRow(ctx, query, string(catReceiverId)).Scan(&catReceiverSex, &catReceiverHasMatched, &catReceiverUserId); err != nil {
		return exc.InternalServerException(fmt.Sprintf("Internal server error: %s", err))
	}

	switch {
	// check cat owner
	case catIssuerUserId != userId:
		return exc.NotFoundException("You cannot match a cat that you do not own!")
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
