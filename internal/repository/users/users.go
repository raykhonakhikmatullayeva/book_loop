package users

import (
	"book_loop/internal/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	createUser     = `INSERT INTO users (login, password_hash, role) VALUES ($1, $2, $3)`
	getUserByLogin = `SELECT id, login, password_hash, role, created_at FROM users WHERE login = $1`
	getUserById = `SELECT id, login, password_hash, role, created_at FROM users WHERE id = $1`
)

type Repo interface {
	CreateUser(ctx context.Context, login, password, role string) error
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
	GetUserByID(ctx context.Context, id int64) (model.User, error)
}
type repo struct {
	repo *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Repo {
	return &repo{repo: pool}
}

func (r *repo) CreateUser(ctx context.Context, login, password, role string) error {
	_, err := r.repo.Exec(ctx, createUser, login, password, role)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *repo) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User
	err := r.repo.QueryRow(ctx, getUserByLogin, login).Scan(&user.Id, &user.Login,
		&user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		return user, fmt.Errorf("get user by login: %w", err)
	}
	return user, nil
}

func (r *repo) GetUserByID(ctx context.Context, id int64) (model.User, error) {
	var user model.User
	err := r.repo.QueryRow(ctx, getUserById, id).Scan(&user.Id, &user.Login, &user.PasswordHash,
		&user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, model.ErrNotFound
		}
	}
	return user, nil
}