package usecases

import (
	"context"

	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/devfajar/task-management-system/internal/repository"
	"github.com/google/uuid"
)

type ColumnUsecase interface {
	Create(ctx context.Context, boardID uuid.UUID, name string, wipLimit *int32, position *float64) (entities.Column, error)
	ListByBoard(ctx context.Context, boardID uuid.UUID) ([]entities.Column, error)
	Update(ctx context.Context, id uuid.UUID, name *string, wipLimit *int32, position *float64) (entities.Column, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}

type columnUC struct{ repo repository.ColumnRepository }

func (c *columnUC) Create(ctx context.Context, boardID uuid.UUID, name string, wipLimit *int32, position *float64) (entities.Column, error) {
	return c.repo.Create(ctx, boardID, name, wipLimit, position)
}

func (c *columnUC) ListByBoard(ctx context.Context, boardID uuid.UUID) ([]entities.Column, error) {
	return c.repo.ListByBoard(ctx, boardID)
}

func (c *columnUC) Update(ctx context.Context, id uuid.UUID, name *string, wipLimit *int32, position *float64) (entities.Column, error) {
	return c.repo.Update(ctx, id, name, wipLimit, position)
}

func (c *columnUC) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return c.repo.SoftDelete(ctx, id)
}

func NewColumnUsecase(repo repository.ColumnRepository) ColumnUsecase { return &columnUC{repo: repo} }
