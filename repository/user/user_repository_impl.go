package user_repository

import (
	user_entity "cats-social/entity/user"
	"context"

	"github.com/jackc/pgx/v5"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Register(ctx context.Context, tx pgx.Tx, user user_entity.User) (user_entity.User, error) {
	var userId string
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3)"
	if err := tx.QueryRow(ctx, query, user.Name, user.Email, user.Password).Scan(&userId); err != nil {
		return user_entity.User{}, err
	}

	user.Id = userId
	return user, nil

}
