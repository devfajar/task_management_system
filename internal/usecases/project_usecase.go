package usecases

import (
	"context"

	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/devfajar/task-management-system/internal/repository"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	Create(ctx context.Context, name string, description *string, ownerID *uuid.UUID) (entities.Project, error)
	Get(ctx context.Context, id uuid.UUID) (entities.Project, error)
	List(ctx context.Context, limit, offset int32) ([]entities.Project, error)
	Update(ctx context.Context, id uuid.UUID, name *string, description *string) (entities.Project, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	ForceDelete(ctx context.Context, id uuid.UUID) error
}

type projectUC struct{ repo repository.ProjectRepository }

func (p *projectUC) Create(ctx context.Context, name string, description *string, ownerID *uuid.UUID) (entities.Project, error) {
	return p.repo.Create(ctx, name, description, ownerID)
}

func (p *projectUC) Get(ctx context.Context, id uuid.UUID) (entities.Project, error) {
	return p.repo.Get(ctx, id)
}

func (p *projectUC) List(ctx context.Context, limit, offset int32) ([]entities.Project, error) {
	return p.repo.List(ctx, limit, offset)
}

func (p *projectUC) Update(ctx context.Context, id uuid.UUID, name *string, description *string) (entities.Project, error) {
	return p.repo.Update(ctx, id, name, description)
}

func (p *projectUC) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return p.repo.SoftDelete(ctx, id)
}

func (p *projectUC) ForceDelete(ctx context.Context, id uuid.UUID) error {
	return p.repo.ForceDelete(ctx, id)
}

func NewProjectUsecase(repo repository.ProjectRepository) ProjectUsecase {
	return &projectUC{repo: repo}
}
