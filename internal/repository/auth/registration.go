package auth

import (
	"book_loop/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	saveRefreshToken = `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)` //nolint:gosec
	getByHash        = `SELECT id, user_id, token_hash, expires_at FROM refresh_tokens WHERE token_hash = $1`
	revokeToken      = `UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1`
)

type repo struct {
	pool *pgxpool.Pool
}

type Repo interface {
	SaveRefresh(ctx context.Context, userId int64, tokenHash string, expiresAt time.Time) error
	GetByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string) error
}

func New(pool *pgxpool.Pool) Repo {
	return &repo{pool: pool}
}

func (r *repo) SaveRefresh(ctx context.Context, userId int64, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, saveRefreshToken, userId, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}
	return nil
}
func (r *repo) GetByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	row := &model.RefreshToken{}
	err := r.pool.QueryRow(ctx, getByHash, tokenHash).Scan(&row.ID, &row.UserID, &row.TokenHash, &row.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return row, model.ErrNotFound
		}
	}
	return row, nil
}

func (r *repo) Revoke(ctx context.Context, tokenHash string) error {
	tag, err := r.pool.Exec(ctx, revokeToken, tokenHash)
	if err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
