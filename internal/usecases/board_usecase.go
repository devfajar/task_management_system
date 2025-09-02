package usecases

import (
	"context"

	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/devfajar/task-management-system/internal/repository"
	"github.com/google/uuid"
)

type BoardUsecase interface {
	Create(ctx context.Context, projectID uuid.UUID, name, desc string) (entities.Board, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]entities.Board, error)
	Get(ctx context.Context, id uuid.UUID) (entities.Board, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	ForceDelete(ctx context.Context, id uuid.UUID) error
}
type boardUC struct{ repo repository.BoardRepository }

func (b *boardUC) Create(ctx context.Context, projectID uuid.UUID, name, desc string) (entities.Board, error) {
	return b.repo.Create(ctx, projectID, name, desc)
}

func (b *boardUC) ListByProject(ctx context.Context, projectID uuid.UUID) ([]entities.Board, error) {
	return b.repo.ListByProject(ctx, projectID)
}

func (b *boardUC) Get(ctx context.Context, id uuid.UUID) (entities.Board, error) {
	return b.repo.Get(ctx, id)
}

func (b *boardUC) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return b.repo.SoftDelete(ctx, id)
}

func (b *boardUC) ForceDelete(ctx context.Context, id uuid.UUID) error {
	return b.repo.ForceDelete(ctx, id)
}

func NewBoardUsecase(repo repository.BoardRepository) BoardUsecase { return &boardUC{repo: repo} }
