package repository

import (
	"context"

	"github.com/devfajar/task-management-system/internal/db"
	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/devfajar/task-management-system/pkg/helper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ColumnRepository interface {
	Create(ctx context.Context, boardID uuid.UUID, name string, wipLimit *int32, position *float64) (entities.Column, error)
	ListByBoard(ctx context.Context, boardID uuid.UUID) ([]entities.Column, error)
	Update(ctx context.Context, id uuid.UUID, name *string, wipLimit *int32, position *float64) (entities.Column, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	NextPosition(ctx context.Context, boardID uuid.UUID) (float64, error)
}

type columnRepo struct{ q *db.Queries }

func (c *columnRepo) Create(ctx context.Context, boardID uuid.UUID, name string, wipLimit *int32, position *float64) (entities.Column, error) {
	pos := 0.0
	if position == nil {
		p, err := c.q.NextColumnPosition(ctx, boardID)
		if err != nil {
			return entities.Column{}, err
		}
		pos = float64(p)
	} else {
		pos = *position
	}

	column, err := c.q.CreateColumn(ctx, db.CreateColumnParams{
		BoardID:  boardID,
		Name:     name,
		WipLimit: helper.ToInt4(wipLimit),
		Position: pos,
	})
	if err != nil {
		return entities.Column{}, err
	}
	return entities.Column{
		ID:        column.ID.String(),
		BoardID:   column.BoardID.String(),
		Name:      column.Name,
		WIPLimit:  helper.Int4PtrToPtr(&column.WipLimit),
		Position:  column.Position,
		CreatedAt: column.CreatedAt,
		UpdatedAt: column.UpdatedAt,
	}, nil
}

func (c *columnRepo) ListByBoard(ctx context.Context, boardID uuid.UUID) ([]entities.Column, error) {
	columns, err := c.q.ListColumnsByBoard(ctx, boardID)
	if err != nil {
		return nil, err
	}
	out := make([]entities.Column, 0, len(columns))
	for _, column := range columns {
		out = append(out, entities.Column{
			ID:        column.ID.String(),
			BoardID:   column.BoardID.String(),
			Name:      column.Name,
			WIPLimit:  helper.Int4PtrToPtr(&column.WipLimit),
			Position:  column.Position,
			CreatedAt: column.CreatedAt,
			UpdatedAt: column.UpdatedAt,
		})
	}
	return out, nil
}

func (c *columnRepo) Update(ctx context.Context, id uuid.UUID, name *string, wipLimit *int32, position *float64) (entities.Column, error) {
	row, err := c.q.UpdateColumn(ctx, db.UpdateColumnParams{
		ID:       id,
		Name:     helper.ToText(name),
		WipLimit: helper.ToInt4(wipLimit),
		Position: helper.ToFloat8(position),
	})
	if err != nil {
		return entities.Column{}, err
	}
	return entities.Column{
		ID:        row.ID.String(),
		BoardID:   row.BoardID.String(),
		Name:      row.Name,
		WIPLimit:  helper.Int4PtrToPtr(&row.WipLimit),
		Position:  row.Position,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (c *columnRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return c.q.SoftDeleteColumn(ctx, id)
}

func (c *columnRepo) NextPosition(ctx context.Context, boardID uuid.UUID) (float64, error) {
	return c.q.NextColumnPosition(ctx, boardID)
}

func NewColumnRepository(pool *pgxpool.Pool) ColumnRepository { return &columnRepo{q: db.New(pool)} }
