package repository

import (
	"context"

	"github.com/devfajar/task-management-system/internal/db"
	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/devfajar/task-management-system/pkg/helper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository interface {
	Create(ctx context.Context, name string, description *string, ownerID *uuid.UUID) (entities.Project, error)
	Get(ctx context.Context, id uuid.UUID) (entities.Project, error)
	List(ctx context.Context, limit, offset int32) ([]entities.Project, error)
	Update(ctx context.Context, id uuid.UUID, name *string, description *string) (entities.Project, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	ForceDelete(ctx context.Context, id uuid.UUID) error
}

type projectRepo struct{ q *db.Queries }

func (p *projectRepo) Create(ctx context.Context, name string, description *string, ownerID *uuid.UUID) (entities.Project, error) {
	var owner pgtype.UUID
	if ownerID != nil {
		owner = pgtype.UUID{Bytes: *ownerID, Valid: true}
	}
	row, err := p.q.CreateProject(ctx, db.CreateProjectParams{
		Name:        name,
		Description: description, // narg => *string
		OwnerID:     owner,       // narg::uuid => pgtype.UUID
	})
	if err != nil {
		return entities.Project{}, err
	}
	return mapProjectFromCreate(row), nil
}

func (p *projectRepo) Get(ctx context.Context, id uuid.UUID) (entities.Project, error) {
	row, err := p.q.GetProjectByID(ctx, id)
	if err != nil {
		return entities.Project{}, err
	}
	return mapProjectFromGet(row), nil
}

func (p *projectRepo) List(ctx context.Context, limit, offset int32) ([]entities.Project, error) {
	rows, err := p.q.ListProjects(ctx, db.ListProjectsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]entities.Project, 0, len(rows))
	for _, rw := range rows {
		out = append(out, mapProjectFromList(rw))
	}
	return out, nil
}

func (p *projectRepo) Update(ctx context.Context, id uuid.UUID, name *string, description *string) (entities.Project, error) {
	row, err := p.q.UpdateProject(ctx, db.UpdateProjectParams{
		ID:          id,
		Name:        helper.ToNullableTextPtr(name),        // narg => *string
		Description: helper.ToNullableTextPtr(description), // narg => *string
	})
	if err != nil {
		return entities.Project{}, err
	}
	return mapProjectFromUpdate(row), nil
}

func (p *projectRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return p.q.SoftDeleteProject(ctx, id)
}

func (p *projectRepo) ForceDelete(ctx context.Context, id uuid.UUID) error {
	return p.q.ForceDeleteProject(ctx, id)
}

func NewProjectRepository(pool *pgxpool.Pool) ProjectRepository { return &projectRepo{q: db.New(pool)} }

func mapProjectFromGet(row db.GetProjectByIDRow) entities.Project {
	return entities.Project{
		ID:          row.ID.String(),
		Name:        row.Name,
		Description: row.Description,
		OwnerID:     helper.OwnerToPtr(row.OwnerID),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func mapProjectFromCreate(row db.CreateProjectRow) entities.Project {
	return entities.Project{
		ID:          row.ID.String(),
		Name:        row.Name,
		Description: row.Description,
		OwnerID:     helper.OwnerToPtr(row.OwnerID),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func mapProjectFromUpdate(row db.UpdateProjectRow) entities.Project {
	return entities.Project{
		ID:          row.ID.String(),
		Name:        row.Name,
		Description: row.Description,
		OwnerID:     helper.OwnerToPtr(row.OwnerID),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func mapProjectFromList(row db.ListProjectsRow) entities.Project {
	return entities.Project{
		ID:          row.ID.String(),
		Name:        row.Name,
		Description: row.Description,
		OwnerID:     helper.OwnerToPtr(row.OwnerID),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
