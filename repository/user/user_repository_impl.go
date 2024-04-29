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
	query := "INSERT INTO users (id, name, email, password) VALUES (gen_random_uuid(), $1, $2, $3) RETURNING id"
	if err := tx.QueryRow(ctx, query, user.Name, user.Email, user.Password).Scan(&userId); err != nil {
		return user_entity.User{}, err
	}

	user.Id = userId
	if err := tx.Commit(ctx); err != nil {
		return user_entity.User{}, err
	}
	return user, nil
}

func (repository *UserRepositoryImpl) Login(ctx context.Context, tx pgx.Tx, user user_entity.User) (user_entity.User, error) {
	query := "SELECT id, name, email, password FROM users WHERE email = $1 LIMIT 1"
	row := tx.QueryRow(ctx, query, user.Email)

	var loggedInUser user_entity.User
	err := row.Scan(&loggedInUser.Id, &loggedInUser.Name, &loggedInUser.Email, &loggedInUser.Password)
	if err != nil {
		return user_entity.User{}, err
	}

	return loggedInUser, nil
}
