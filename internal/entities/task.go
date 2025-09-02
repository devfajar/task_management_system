package entities

import (
	"time"
)

type Task struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	BoardID     *string    `json:"board_id,omitempty"`
	ColumnID    *string    `json:"column_id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    int32      `json:"priority"`
	AssigneeID  *string    `json:"assignee_id,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Position    *float64   `json:"position,omitempty"`
	StoryPoints *float64   `json:"story_points,omitempty"`

	OriginalEstimateSeconds  int32 `json:"original_estimate_seconds"`
	RemainingEstimateSeconds int32 `json:"remaining_estimate_seconds"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
