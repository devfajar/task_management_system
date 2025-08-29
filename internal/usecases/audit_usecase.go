package usecases

import (
	"context"

	"github.com/devfajar/task-management-system/internal/repository"
	"github.com/google/uuid"
)

type AuditUsecase interface {
	Log(ctx context.Context, userID *uuid.UUID, action, entity string, entityID *uuid.UUID, metadata map[string]any, ip, ua string) error
}

type auditUC struct{ repo repository.SecurityRepository }

func NewAuditUsecase(repo repository.SecurityRepository) AuditUsecase { return &auditUC{repo: repo} }

func (a *auditUC) Log(ctx context.Context, userID *uuid.UUID, action, entity string, entityID *uuid.UUID, metadata map[string]any, ip, ua string) error {
	return a.repo.InsertAudit(ctx, userID, action, entity, entityID, metadata, ip, ua)
}
