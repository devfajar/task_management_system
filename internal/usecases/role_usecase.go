package usecases

import (
	"context"

	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/devfajar/task-management-system/internal/repository"
	"github.com/google/uuid"
)

type RoleUsecase interface {
	CreateRole(ctx context.Context, key, name string) (entities.Role, error)
	ListRoles(ctx context.Context) ([]entities.Role, error)
	DeleteRole(ctx context.Context, id uuid.UUID) error
	Assign(ctx context.Context, userID uuid.UUID, roleKey string) error
	Revoke(ctx context.Context, userID uuid.UUID, roleKey string) error
}

type roleUC struct {
	roles repository.RoleRepository
}

func NewRoleUsecase(roles repository.RoleRepository) RoleUsecase { return &roleUC{roles: roles} }

func (r *roleUC) CreateRole(ctx context.Context, key, name string) (entities.Role, error) {
	return r.roles.CreateRole(ctx, key, name)
}

func (r *roleUC) ListRoles(ctx context.Context) ([]entities.Role, error) {
	return r.roles.ListRoles(ctx)
}

func (r *roleUC) DeleteRole(ctx context.Context, id uuid.UUID) error {
	return r.roles.DeleteRole(ctx, id)
}

func (r *roleUC) Assign(ctx context.Context, userID uuid.UUID, roleKey string) error {
	role, err := r.roles.GetRoleByKey(ctx, roleKey)
	if err != nil {
		return err
	}
	return r.roles.AssignRoleToUser(ctx, userID, role.ID)
}

func (r *roleUC) Revoke(ctx context.Context, userID uuid.UUID, roleKey string) error {
	role, err := r.roles.GetRoleByKey(ctx, roleKey)
	if err != nil {
		return err
	}
	return r.roles.RevokeRoleFromUser(ctx, userID, role.ID)
}
