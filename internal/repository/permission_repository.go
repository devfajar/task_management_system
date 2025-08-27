package repository

import (
	"context"

	"github.com/devfajar/task-management-system/internal/db"
	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PermissionRepository interface {
	CreatePermission(ctx context.Context, key, desc string) (entities.Permission, error)
	GetPermissionByKey(ctx context.Context, key string) (entities.Permission, error)
	ListPermissions(ctx context.Context) ([]entities.Permission, error)
	DeletePermission(ctx context.Context, id uuid.UUID) error

	GrantPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error
	RevokePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error

	HasPermission(ctx context.Context, userID uuid.UUID, permKey string) (bool, error)
	ListUserPermissionKeys(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type permissionRepo struct{ q *db.Queries }

func NewPermissionRepository(pool *pgxpool.Pool) PermissionRepository {
	return &permissionRepo{q: db.New(pool)}
}

func (p permissionRepo) CreatePermission(ctx context.Context, key, desc string) (entities.Permission, error) {
	permission, err := p.q.CreatePermission(ctx, db.CreatePermissionParams{Key: key, Description: desc})
	if err != nil {
		return entities.Permission{}, err
	}
	return entities.Permission{
		ID:          permission.ID,
		Key:         permission.Key,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
	}, nil
}

func (p permissionRepo) GetPermissionByKey(ctx context.Context, key string) (entities.Permission, error) {
	permission, err := p.q.GetPermissionByKey(ctx, key)
	if err != nil {
		return entities.Permission{}, err
	}
	return entities.Permission{
		ID:          permission.ID,
		Key:         permission.Key,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
	}, nil
}

func (p permissionRepo) ListPermissions(ctx context.Context) ([]entities.Permission, error) {
	permissions, err := p.q.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]entities.Permission, 0, len(permissions))
	for _, permission := range permissions {
		out = append(out, entities.Permission{
			ID:          permission.ID,
			Key:         permission.Key,
			Description: permission.Description,
		})
	}
	return out, nil
}

func (p permissionRepo) DeletePermission(ctx context.Context, id uuid.UUID) error {
	return p.q.DeletePermission(ctx, id)
}

func (p permissionRepo) GrantPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	return p.q.GrantPermissionToRole(ctx, db.GrantPermissionToRoleParams{RoleID: roleID, PermissionID: permissionID})
}

func (p permissionRepo) RevokePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	return p.q.RevokePermissionFromRole(ctx, db.RevokePermissionFromRoleParams{
		RoleID:       roleID,
		PermissionID: permissionID,
	})
}

func (p permissionRepo) HasPermission(ctx context.Context, userID uuid.UUID, permKey string) (bool, error) {
	return p.q.HasPermission(ctx, db.HasPermissionParams{
		UserID: userID,
		Key:    permKey,
	})
}

func (p permissionRepo) ListUserPermissionKeys(ctx context.Context, userID uuid.UUID) ([]string, error) {
	permissions, err := p.q.ListUserPermissionKeys(ctx, userID)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		keys = append(keys, permission)
	}
	return keys, nil
}
