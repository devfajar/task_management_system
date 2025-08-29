package repository

import (
	"context"
	"encoding/json"

	"github.com/devfajar/task-management-system/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SecurityRepository interface {
	GetTokenVersion(ctx context.Context, userID uuid.UUID) (int32, error)
	BumpTokenVersion(ctx context.Context, userID uuid.UUID) error
	InsertAudit(ctx context.Context, userID *uuid.UUID, action, entity string, entityID *uuid.UUID, metadata map[string]any, ip, ua string) error
}

type securityRepo struct{ q *db.Queries }

func (s *securityRepo) GetTokenVersion(ctx context.Context, userID uuid.UUID) (int32, error) {
	return s.q.GetUserTokenVersion(ctx, userID)
}

func (s *securityRepo) BumpTokenVersion(ctx context.Context, userID uuid.UUID) error {
	return s.q.BumpUserTokenVersion(ctx, userID)
}

func (s *securityRepo) InsertAudit(ctx context.Context, userID *uuid.UUID, action, entity string, entityID *uuid.UUID, metadata map[string]any, ip, ua string) error {
	var uid pgtype.UUID
	if userID != nil {
		uid = pgtype.UUID{Bytes: *userID, Valid: true}
	}
	var eid pgtype.UUID
	if entityID != nil {
		eid = pgtype.UUID{Bytes: *entityID, Valid: true}
	}
	meta, _ := json.Marshal(metadata)
	_, err := s.q.InsertAuditLog(ctx, db.InsertAuditLogParams{
		UserID:    uid,
		Action:    action,
		Entity:    entity,
		EntityID:  eid,
		Metadata:  meta,
		IpAddress: toNullableText(ip),
		UserAgent: toNullableText(ua),
	})
	return err
}

func NewSecurityRepository(pool *pgxpool.Pool) SecurityRepository {
	return &securityRepo{q: db.New(pool)}
}
