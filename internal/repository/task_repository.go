package repository

import (
	"context"
	"time"

	"github.com/devfajar/task-management-system/internal/db"
	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/devfajar/task-management-system/pkg/helper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTaskParams struct {
	ProjectID                uuid.UUID
	BoardID                  uuid.UUID
	ColumnID                 uuid.UUID
	Title                    string
	Description              *string
	Status                   *string
	Priority                 *int32
	AssigneeID               *uuid.UUID
	DueDate                  *time.Time
	Position                 *float64
	StoryPoints              *float64
	OriginalEstimateSeconds  *int32
	RemainingEstimateSeconds *int32
}

type UpdateTaskFieldsParams struct {
	ID          uuid.UUID
	Title       *string
	Description *string
	Status      *string
	Priority    *int32
	AssigneeID  *uuid.UUID
	DueDate     *time.Time
	StoryPoints *float64
}

type TaskRepository interface {
	NextPosition(ctx context.Context, columnID uuid.UUID) (float64, error)
	CreateInColumn(ctx context.Context, p CreateTaskParams) (entities.Task, error)
	MoveToColumn(ctx context.Context, taskID uuid.UUID, toColumn uuid.UUID, newPos float64) error
	Reorder(ctx context.Context, taskID uuid.UUID, newPos float64) error
	UpdateFields(ctx context.Context, p UpdateTaskFieldsParams) (entities.Task, error)
	UpdateEstimates(ctx context.Context, id uuid.UUID, orig *int32, remain *int32) (entities.Task, error)
	ListByColumn(ctx context.Context, columnID uuid.UUID, limit, offset int32) ([]entities.Task, error)
}

type taskRepo struct{ q *db.Queries }

func (t *taskRepo) NextPosition(ctx context.Context, columnID uuid.UUID) (float64, error) {
	return t.q.NextTaskPosition(ctx, columnID)
}

func (t *taskRepo) CreateInColumn(ctx context.Context, p CreateTaskParams) (entities.Task, error) {
	var ass pgtype.UUID
	if p.AssigneeID != nil {
		ass = pgtype.UUID{Bytes: *p.AssigneeID, Valid: true}
	}
	var due pgtype.Timestamptz
	if p.DueDate != nil {
		due = pgtype.Timestamptz{Time: *p.DueDate, Valid: true}
	}

	sp := (*float64)(nil)
	if p.StoryPoints != nil {
		sp = p.StoryPoints
	}

	pos := 0.0
	if p.Position == nil {
		np, err := t.q.NextTaskPosition(ctx, p.ColumnID)
		if err != nil {
			return entities.Task{}, err
		}
		pos = np
	} else {
		pos = *p.Position
	}

	row, err := t.q.CreateTaskInColumn(ctx, db.CreateTaskInColumnParams{
		ProjectID:                p.ProjectID,
		Title:                    p.Title,
		Description:              helper.ToText(p.Description),
		Status:                   helper.ToText(p.Status),
		Priority:                 helper.ToInt4(p.Priority),
		AssigneeID:               ass,
		DueDate:                  due,
		BoardID:                  p.BoardID,
		ColumnID:                 p.ColumnID,
		Position:                 pos,
		StoryPoints:              helper.ToFloat8(sp),
		OriginalEstimateSeconds:  helper.ToInt4(p.OriginalEstimateSeconds),
		RemainingEstimateSeconds: helper.ToInt4(p.RemainingEstimateSeconds),
	})
	if err != nil {
		return entities.Task{}, err
	}

	return mapTaskFromCreate(row), nil
}

func (t *taskRepo) MoveToColumn(ctx context.Context, taskID uuid.UUID, toColumn uuid.UUID, newPos float64) error {
	return t.q.MoveTaskToColumn(ctx, db.MoveTaskToColumnParams{
		ID:       taskID,
		ColumnID: helper.ToUUID(toColumn),
		Position: helper.ToFloat8(&newPos),
	})
}

func (t *taskRepo) Reorder(ctx context.Context, taskID uuid.UUID, newPos float64) error {
	return t.q.ReorderTaskPosition(ctx, db.ReorderTaskPositionParams{
		ID:       taskID,
		Position: helper.ToFloat8(&newPos),
	})
}

func (t *taskRepo) UpdateFields(ctx context.Context, p UpdateTaskFieldsParams) (entities.Task, error) {
	var ass pgtype.UUID
	if p.AssigneeID != nil {
		ass = pgtype.UUID{Bytes: *p.AssigneeID, Valid: true}
	}
	var due pgtype.Timestamptz
	if p.DueDate != nil {
		due = pgtype.Timestamptz{Time: *p.DueDate, Valid: true}
	}

	row, err := t.q.UpdateTaskSprintFields(ctx, db.UpdateTaskSprintFieldsParams{
		ID:          p.ID,
		Title:       helper.ToText(p.Title),
		Description: helper.ToText(p.Description),
		Status:      helper.ToText(p.Status),
		Priority:    helper.ToInt4(p.Priority),
		AssigneeID:  ass,
		DueDate:     due,
		StoryPoints: helper.ToFloat8(p.StoryPoints),
	})
	if err != nil {
		return entities.Task{}, err
	}
	return mapTaskFromUpdate(row), nil
}

func (t *taskRepo) UpdateEstimates(ctx context.Context, id uuid.UUID, orig *int32, remain *int32) (entities.Task, error) {
	row, err := t.q.UpdateTaskEstimates(ctx, db.UpdateTaskEstimatesParams{
		ID:                       id,
		OriginalEstimateSeconds:  helper.ToInt4(orig),
		RemainingEstimateSeconds: helper.ToInt4(remain),
	})
	if err != nil {
		return entities.Task{}, err
	}
	return mapTaskFromEstimate(row), nil
}

func (t *taskRepo) ListByColumn(ctx context.Context, columnID uuid.UUID, limit, offset int32) ([]entities.Task, error) {
	rows, err := t.q.ListTasksByColumn(ctx, db.ListTasksByColumnParams{
		ColumnID: helper.ToUUID(columnID),
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, err
	}
	out := make([]entities.Task, 0, len(rows))
	for _, x := range rows {
		out = append(out, mapTaskFromList(x))
	}
	return out, nil
}

func NewTaskRepository(pool *pgxpool.Pool) TaskRepository { return &taskRepo{q: db.New(pool)} }

// Mapper
func mapTaskFromCreate(x db.CreateTaskInColumnRow) entities.Task {
	return entities.Task{
		ID:                       x.ID.String(),
		ProjectID:                x.ProjectID.String(),
		BoardID:                  helper.UuidToPtr(x.BoardID),
		ColumnID:                 helper.UuidToPtr(x.ColumnID),
		Title:                    x.Title,
		Description:              x.Description,
		Status:                   x.Status,
		Priority:                 x.Priority,
		AssigneeID:               helper.UuidToPtr(x.AssigneeID),
		DueDate:                  helper.TimestamptzToPtr(x.DueDate),
		Position:                 helper.Float64PtrFromPg(x.Position),
		StoryPoints:              helper.Float64PtrFromPg(x.StoryPoints),
		OriginalEstimateSeconds:  x.OriginalEstimateSeconds,
		RemainingEstimateSeconds: x.RemainingEstimateSeconds,
		CreatedAt:                x.CreatedAt,
		UpdatedAt:                x.UpdatedAt,
	}
}

func mapTaskFromUpdate(x db.UpdateTaskSprintFieldsRow) entities.Task {
	return entities.Task{
		ID:                       x.ID.String(),
		ProjectID:                x.ProjectID.String(),
		BoardID:                  helper.UuidToPtr(x.BoardID),
		ColumnID:                 helper.UuidToPtr(x.ColumnID),
		Title:                    x.Title,
		Description:              x.Description,
		Status:                   x.Status,
		Priority:                 x.Priority,
		AssigneeID:               helper.UuidToPtr(x.AssigneeID),
		DueDate:                  helper.TimestamptzToPtr(x.DueDate),
		Position:                 helper.Float64PtrFromPg(x.Position),
		StoryPoints:              helper.Float64PtrFromPg(x.StoryPoints),
		OriginalEstimateSeconds:  x.OriginalEstimateSeconds,
		RemainingEstimateSeconds: x.RemainingEstimateSeconds,
		CreatedAt:                x.CreatedAt,
		UpdatedAt:                x.UpdatedAt,
	}
}

func mapTaskFromEstimate(x db.UpdateTaskEstimatesRow) entities.Task {
	return entities.Task{
		ID:                       x.ID.String(),
		ProjectID:                x.ProjectID.String(),
		BoardID:                  helper.UuidToPtr(x.BoardID),
		ColumnID:                 helper.UuidToPtr(x.ColumnID),
		Title:                    x.Title,
		Description:              x.Description,
		Status:                   x.Status,
		Priority:                 x.Priority,
		AssigneeID:               helper.UuidToPtr(x.AssigneeID),
		DueDate:                  helper.TimestamptzToPtr(x.DueDate),
		Position:                 helper.Float64PtrFromPg(x.Position),
		StoryPoints:              helper.Float64PtrFromPg(x.StoryPoints),
		OriginalEstimateSeconds:  x.OriginalEstimateSeconds,
		RemainingEstimateSeconds: x.RemainingEstimateSeconds,
		CreatedAt:                x.CreatedAt,
		UpdatedAt:                x.UpdatedAt,
	}
}

func mapTaskFromList(x db.ListTasksByColumnRow) entities.Task {
	return entities.Task{
		ID:                       x.ID.String(),
		ProjectID:                x.ProjectID.String(),
		BoardID:                  helper.UuidToPtr(x.BoardID),
		ColumnID:                 helper.UuidToPtr(x.ColumnID),
		Title:                    x.Title,
		Description:              x.Description,
		Status:                   x.Status,
		Priority:                 x.Priority,
		AssigneeID:               helper.UuidToPtr(x.AssigneeID),
		DueDate:                  helper.TimestamptzToPtr(x.DueDate),
		Position:                 helper.Float64PtrFromPg(x.Position),
		StoryPoints:              helper.Float64PtrFromPg(x.StoryPoints),
		OriginalEstimateSeconds:  x.OriginalEstimateSeconds,
		RemainingEstimateSeconds: x.RemainingEstimateSeconds,
		CreatedAt:                x.CreatedAt,
		UpdatedAt:                x.UpdatedAt,
	}
}
