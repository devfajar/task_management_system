package usecases

import (
	"context"

	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/devfajar/task-management-system/internal/repository"
	"github.com/google/uuid"
)

type PermissionUsecase interface {
	CreatePermission(ctx context.Context, key, desc string) (entities.Permission, error)
	ListPermissions(ctx context.Context) ([]entities.Permission, error)
	DeletePermission(ctx context.Context, id uuid.UUID) error

	Grant(ctx context.Context, roleKey, permKey string) error
	Revoke(ctx context.Context, roleKey, permKey string) error

	Has(ctx context.Context, userID uuid.UUID, permKey string) (bool, error)
	ListUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type permUC struct {
	perms repository.PermissionRepository
	roles repository.RoleRepository
}

func NewPermissionUsecase(perms repository.PermissionRepository, roles repository.RoleRepository) PermissionUsecase {
	return &permUC{perms: perms, roles: roles}
}

func (p *permUC) CreatePermission(ctx context.Context, key, desc string) (entities.Permission, error) {
	return p.perms.CreatePermission(ctx, key, desc)
}

func (p *permUC) ListPermissions(ctx context.Context) ([]entities.Permission, error) {
	return p.perms.ListPermissions(ctx)
}

func (p *permUC) DeletePermission(ctx context.Context, id uuid.UUID) error {
	return p.perms.DeletePermission(ctx, id)
}

func (p *permUC) Grant(ctx context.Context, roleKey, permKey string) error {
	role, err := p.roles.GetRoleByKey(ctx, roleKey)
	if err != nil {
		return err
	}
	permission, err := p.perms.GetPermissionByKey(ctx, permKey)
	if err != nil {
		return err
	}
	return p.perms.GrantPermissionToRole(ctx, role.ID, permission.ID)
}

func (p *permUC) Revoke(ctx context.Context, roleKey, permKey string) error {
	role, err := p.roles.GetRoleByKey(ctx, roleKey)
	if err != nil {
		return err
	}
	permission, err := p.perms.GetPermissionByKey(ctx, permKey)
	if err != nil {
		return err
	}
	return p.perms.RevokePermissionFromRole(ctx, role.ID, permission.ID)
}

func (p *permUC) Has(ctx context.Context, userID uuid.UUID, permKey string) (bool, error) {
	return p.perms.HasPermission(ctx, userID, permKey)
}

func (p *permUC) ListUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return p.perms.ListUserPermissionKeys(ctx, userID)
}
