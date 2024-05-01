package match_repository

import (
	match_entity "cats-social/entity/match"
	exc "cats-social/exceptions"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type matchRepositoryImpl struct {
}

func NewMatchRepository() MatchRepository {
	return &matchRepositoryImpl{}
}

func (repository *matchRepositoryImpl) Create(ctx context.Context, tx pgx.Tx, match match_entity.Match, userId string) error {
	if err := checkCatExists(ctx, tx, match.CatIssuerId, match.CatReceiverId); err != nil {
		return err
	}
	if err := validateMatchCatCriteria(ctx, tx, match.CatIssuerId, match.CatReceiverId, userId); err != nil {
		return err
	}

	var matchId string
	query := `INSERT INTO matches (id, message, cat_issuer_id, cat_receiver_id)
	SELECT 
		gen_random_uuid(), $1, $2, $3
	WHERE EXISTS (
		SELECT 1 FROM users WHERE id = $4
	)
	RETURNING id;
	`
	if err := tx.QueryRow(ctx, query, match.Message, match.CatIssuerId, match.CatReceiverId, string(userId)).Scan(&matchId); err != nil {
		tx.Rollback(ctx)
		if err == pgx.ErrNoRows {
			return exc.BadRequestException("Invalid user id")
		}
		return err
	}

	match.Id = matchId
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (repository *matchRepositoryImpl) Approve(ctx context.Context, tx pgx.Tx, match match_entity.Match, userId string) error {

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
