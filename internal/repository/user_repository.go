package repository

import (
	"context"
	"errors"

	"github.com/devfajar/task-management-system/internal/db"
	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/devfajar/task-management-system/pkg/helper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEmailExists = errors.New("email already exists")

type UserRepository interface {
	Create(ctx context.Context, u entities.User) (entities.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (entities.User, error)
	List(ctx context.Context, limit, offset int32) ([]entities.User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, username, email *string, isActive *bool) (entities.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, newHash string) error
	Delete(ctx context.Context, id uuid.UUID) error // alias SoftDelete
	ForceDelete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
}

type userRepo struct {
	q *db.Queries
}

// NewUserRepository
func NewUserRepository(pool *pgxpool.Pool) UserRepository { return &userRepo{q: db.New(pool)} }

// Create new user
func (userRepo *userRepo) Create(ctx context.Context, u entities.User) (entities.User, error) {
	row, err := userRepo.q.CreateUser(ctx, db.CreateUserParams{
		ID:       u.Id,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return entities.User{}, ErrEmailExists
		}
		return entities.User{}, err
	}
	return entities.User{
		Id:        row.ID,
		Username:  row.Username,
		Email:     row.Email,
		IsActive:  row.IsActive,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

// Get user by id
func (userRepo userRepo) GetByID(ctx context.Context, id uuid.UUID) (entities.User, error) {
	row, err := userRepo.q.GetUserByID(ctx, id)
	if err != nil {
		return entities.User{}, err
	}

	return entities.User{
		Id:        row.ID,
		Username:  row.Username,
		Email:     row.Email,
		IsActive:  row.IsActive,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

// List all users
func (userRepo userRepo) List(ctx context.Context, limit, offset int32) ([]entities.User, error) {
	rows, err := userRepo.q.ListUsers(ctx, db.ListUsersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	out := make([]entities.User, 0, len(rows))
	for _, user := range rows {
		out = append(out, entities.User{
			Id:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
	return out, nil
}

// Update user profile
func (userRepo userRepo) UpdateProfile(ctx context.Context, id uuid.UUID, username, email *string, isActive *bool) (entities.User, error) {
	row, err := userRepo.q.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		ID:       id,
		Username: helper.ToText(username),
		Email:    helper.ToText(email),
		IsActive: helper.ToBool(isActive),
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return entities.User{}, ErrEmailExists
		}
		return entities.User{}, err
	}
	return entities.User{
		Id:        row.ID,
		Username:  row.Username,
		Email:     row.Email,
		IsActive:  row.IsActive,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

// Update user password
func (userRepo userRepo) UpdatePassword(ctx context.Context, id uuid.UUID, newHash string) error {
	return userRepo.q.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{ID: id, Password: newHash})
}

// Soft delete user
func (userRepo *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return userRepo.q.SoftDeleteUser(ctx, id)
}

// Hard delete user
func (userRepo *userRepo) ForceDelete(ctx context.Context, id uuid.UUID) error {
	return userRepo.q.DeleteUser(ctx, id)
}

// Restore user
func (userRepo *userRepo) Restore(ctx context.Context, id uuid.UUID) error {
	return userRepo.q.RestoreUser(ctx, id)
}
