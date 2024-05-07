package user_repository

import (
	user_entity "cats-social/entity/user"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepositoryImpl struct {
	DBPool *pgxpool.Pool
}

func NewUserRepository(dbPool *pgxpool.Pool) UserRepository {
	return &userRepositoryImpl{
		DBPool: dbPool,
	}
}

func (repository *userRepositoryImpl) Register(ctx context.Context, user user_entity.User) (user_entity.User, error) {
	var userId string
	query := "INSERT INTO users (id, name, email, password) VALUES (gen_random_uuid(), $1, $2, $3) RETURNING id"
	if err := repository.DBPool.QueryRow(ctx, query, user.Name, user.Email, user.Password).Scan(&userId); err != nil {
		return user_entity.User{}, err
	}

	user.Id = userId
	return user, nil
}

func (repository *userRepositoryImpl) Login(ctx context.Context, user user_entity.User) (user_entity.User, error) {
	query := "SELECT id, name, email, password FROM users WHERE email = $1 LIMIT 1"
	row := repository.DBPool.QueryRow(ctx, query, user.Email)

	var loggedInUser user_entity.User
	err := row.Scan(&loggedInUser.Id, &loggedInUser.Name, &loggedInUser.Email, &loggedInUser.Password)
	if err != nil {
		return user_entity.User{}, err
	}

	return loggedInUser, nil
}
