package repository

import (
	"context"

	"github.com/devfajar/task-management-system/internal/db"
	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BoardRepository interface {
	Create(ctx context.Context, projectID uuid.UUID, name, desc string) (entities.Board, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]entities.Board, error)
	Get(ctx context.Context, id uuid.UUID) (entities.Board, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	ForceDelete(ctx context.Context, id uuid.UUID) error
}

type boardRepo struct{ q *db.Queries }

func (b *boardRepo) Create(ctx context.Context, projectID uuid.UUID, name, desc string) (entities.Board, error) {
	board, err := b.q.CreateBoard(ctx, db.CreateBoardParams{
		ProjectID:   projectID,
		Name:        name,
		Description: desc,
	})
	if err != nil {
		return entities.Board{}, err
	}
	return entities.Board{
		ID:          board.ID.String(),
		ProjectID:   board.ProjectID.String(),
		Name:        board.Name,
		Description: board.Description,
		CreatedAt:   board.CreatedAt,
		UpdatedAt:   board.UpdatedAt,
	}, nil
}

func (b *boardRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]entities.Board, error) {
	boards, err := b.q.ListBoardsByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	out := make([]entities.Board, 0, len(boards))
	for _, board := range boards {
		out = append(out, entities.Board{
			ID:          board.ID.String(),
			ProjectID:   board.ProjectID.String(),
			Name:        board.Name,
			Description: board.Description,
			CreatedAt:   board.CreatedAt,
			UpdatedAt:   board.UpdatedAt,
		})
	}
	return out, nil
}

func (b *boardRepo) Get(ctx context.Context, id uuid.UUID) (entities.Board, error) {
	board, err := b.q.GetBoardByID(ctx, id)
	if err != nil {
		return entities.Board{}, err
	}
	return entities.Board{
		ID:          board.ID.String(),
		ProjectID:   board.ProjectID.String(),
		Name:        board.Name,
		Description: board.Description,
		CreatedAt:   board.CreatedAt,
		UpdatedAt:   board.UpdatedAt,
	}, nil
}

func (b *boardRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return b.q.SoftDeleteBoard(ctx, id)
}

func (b *boardRepo) ForceDelete(ctx context.Context, id uuid.UUID) error {
	return b.q.ForceDeleteBoard(ctx, id)
}

func NewBoardRepository(pool *pgxpool.Pool) BoardRepository { return &boardRepo{q: db.New(pool)} }
