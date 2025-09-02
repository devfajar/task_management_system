package usecases

import (
	"context"
	"errors"

	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/devfajar/task-management-system/internal/repository"
	"github.com/devfajar/task-management-system/pkg/helper"
	"github.com/google/uuid"
)

var (
	ErrEmailExists = repository.ErrEmailExists
	ErrBadInput    = errors.New("bad input")
)

type UserUsecase interface {
	Create(ctx context.Context, username, email, password string) (entities.User, error)
	GetById(ctx context.Context, id uuid.UUID) (entities.User, error)
	List(ctx context.Context, limit, offset int32) ([]entities.User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, username, email *string, isActive *bool) (entities.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error
	Delete(ctx context.Context, id uuid.UUID) error      // soft
	ForceDelete(ctx context.Context, id uuid.UUID) error // hard (admin-only)
	Restore(ctx context.Context, id uuid.UUID) error
}

type userUC struct{ repo repository.UserRepository }

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase { return &userUC{repo: userRepo} }

// Create a new user
func (u *userUC) Create(ctx context.Context, username, email, password string) (entities.User, error) {
	if username == "" || email == "" || password == "" {
		return entities.User{}, ErrBadInput
	}
	hash, err := helper.Hash(password)
	if err != nil {
		return entities.User{}, err
	}
	return u.repo.Create(ctx, entities.User{
		Id: uuid.New(), Username: username, Email: email, Password: hash,
	})
}

// GetById Get user by id
func (u *userUC) GetById(ctx context.Context, id uuid.UUID) (entities.User, error) {
	return u.repo.GetByID(ctx, id)
}

// List all users
func (u *userUC) List(ctx context.Context, limit, offset int32) ([]entities.User, error) {
	if limit <= 0 {
		limit = 20
	}
	return u.repo.List(ctx, limit, offset)
}

// UpdateProfile Update user profile
func (u *userUC) UpdateProfile(ctx context.Context, id uuid.UUID, username, email *string, isActive *bool) (entities.User, error) {
	return u.repo.UpdateProfile(ctx, id, username, email, isActive)
}

// UpdatePassword Update user password
func (u *userUC) UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	if newPassword == "" {
		return ErrBadInput
	}
	hash, err := helper.Hash(newPassword)
	if err != nil {
		return err
	}
	return u.repo.UpdatePassword(ctx, id, hash)
}

// Delete user
func (u *userUC) Delete(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}

// ForceDelete user
func (u *userUC) ForceDelete(ctx context.Context, id uuid.UUID) error {
	return u.repo.ForceDelete(ctx, id)
}

func (u *userUC) Restore(ctx context.Context, id uuid.UUID) error {
	return u.repo.Restore(ctx, id)
}
