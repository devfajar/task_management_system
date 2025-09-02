package usecases

import (
	"context"
	"errors"

	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/devfajar/task-management-system/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrInvalidMove = errors.New("invalid move target")
)

type TaskUsecase interface {
	Create(ctx context.Context, p repository.CreateTaskParams) (entities.Task, error)
	Move(ctx context.Context, taskID, toColumn uuid.UUID, newPos *float64) error
	Reorder(ctx context.Context, taskID uuid.UUID, leftPos *float64, rightPos *float64) error
	UpdateFields(ctx context.Context, p repository.UpdateTaskFieldsParams) (entities.Task, error)
	UpdateEstimates(ctx context.Context, id uuid.UUID, orig *int32, remain *int32) (entities.Task, error)
	ListByColumn(ctx context.Context, columnID uuid.UUID, limit, offset int32) ([]entities.Task, error)
}

type taskUC struct{ repo repository.TaskRepository }

func (t *taskUC) Create(ctx context.Context, p repository.CreateTaskParams) (entities.Task, error) {
	return t.repo.CreateInColumn(ctx, p)
}

func (t *taskUC) Move(ctx context.Context, taskID, toColumn uuid.UUID, newPos *float64) error {
	pos := 0.0
	if newPos == nil {
		p, err := t.repo.NextPosition(ctx, toColumn)
		if err != nil {
			return err
		}
		pos = p
	} else {
		pos = *newPos
	}
	return t.repo.MoveToColumn(ctx, taskID, toColumn, pos)
}

func (t *taskUC) Reorder(ctx context.Context, taskID uuid.UUID, leftPos *float64, rightPos *float64) error {
	switch {
	case leftPos != nil && rightPos != nil:
		pos := (*leftPos + *rightPos) / 2.0
		return t.repo.Reorder(ctx, taskID, pos)
	case leftPos != nil && rightPos == nil:
		return t.repo.Reorder(ctx, taskID, *leftPos+1024)
	case leftPos == nil && rightPos != nil:
		return t.repo.Reorder(ctx, taskID, *rightPos-1)
	default:
		return ErrInvalidMove
	}
}

func (t *taskUC) UpdateFields(ctx context.Context, p repository.UpdateTaskFieldsParams) (entities.Task, error) {
	return t.repo.UpdateFields(ctx, p)
}

func (t *taskUC) UpdateEstimates(ctx context.Context, id uuid.UUID, orig *int32, remain *int32) (entities.Task, error) {
	return t.repo.UpdateEstimates(ctx, id, orig, remain)
}

func (t *taskUC) ListByColumn(ctx context.Context, columnID uuid.UUID, limit, offset int32) ([]entities.Task, error) {
	return t.repo.ListByColumn(ctx, columnID, limit, offset)
}

func NewTaskUsecase(repo repository.TaskRepository) TaskUsecase { return &taskUC{repo: repo} }
