package repository

import (
	"context"
	"time"

	"github.com/devfajar/task-management-system/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	GetUserCredentialsByEmail(ctx context.Context, email string) (uid uuid.UUID, hash string, isActive bool, deleted bool, err error)
	SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time, ua, ip string) error
	FindRefreshToken(ctx context.Context, tokenHash string) (id uuid.UUID, userID uuid.UUID, expiresAt time.Time, revokedAt *time.Time, err error)
	RevokeRefreshToken(ctx context.Context, id uuid.UUID) error
}

type authRepo struct{ q *db.Queries }

func NewAuthRepository(pool *pgxpool.Pool) AuthRepository { return &authRepo{q: db.New(pool)} }

func (auth authRepo) GetUserCredentialsByEmail(ctx context.Context, email string) (uuid.UUID, string, bool, bool, error) {
	data, err := auth.q.GetUserCredentialsByEmail(ctx, email)
	if err != nil {
		return uuid.Nil, "", false, false, err
	}
	var deleted bool
	if dz, ok := any(data.DeletedAt).(pgtype.Timestamptz); ok {
		deleted = dz.Valid
	}
	return data.ID, data.Password, data.IsActive, deleted, nil
}

func (auth *authRepo) SaveRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time, ua, ip string) error {
	_, err := auth.q.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		UserAgent: toNullableText(ua),
		IpAddress: toNullableText(ip),
	})
	return err
}

func (auth authRepo) FindRefreshToken(ctx context.Context, tokenHash string) (id uuid.UUID, userID uuid.UUID, expiresAt time.Time, revokedAt *time.Time, err error) {
	data, err := auth.q.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return uuid.Nil, uuid.Nil, time.Time{}, nil, err
	}
	var revoked *time.Time
	if data.RevokedAt.Valid {
		revoked = &data.RevokedAt.Time
	}
	return data.ID, data.UserID, data.ExpiresAt, revoked, nil
}

func (auth authRepo) RevokeRefreshToken(ctx context.Context, id uuid.UUID) error {
	return auth.q.RevokeRefreshToken(ctx, id)
}

// Private
func toNullableText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}
