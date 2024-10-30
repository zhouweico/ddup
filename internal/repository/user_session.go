package repository

import (
	"context"
	"ddup-apis/internal/model"
)

func (r *Repository) CreateUserSession(ctx context.Context, session *model.UserSession) error {
	query := `
		INSERT INTO user_sessions (user_id, token, is_valid, expired_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query,
		session.UserID, session.Token, session.IsValid, session.ExpiredAt)
	return err
}

func (r *Repository) GetUserSessionByToken(ctx context.Context, token string) (*model.UserSession, error) {
	var session model.UserSession
	query := `SELECT user_id, token, is_valid, expired_at FROM user_sessions WHERE token = $1`
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&session.UserID, &session.Token, &session.IsValid, &session.ExpiredAt,
	)
	return &session, err
}

func (r *Repository) InvalidateUserSession(ctx context.Context, token string) error {
	query := `UPDATE user_sessions SET is_valid = false WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}
