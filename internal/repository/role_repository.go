package repository

import (
	"context"

	"github.com/devfajar/task-management-system/internal/db"
	"github.com/devfajar/task-management-system/internal/entities"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository interface {
	CreateRole(ctx context.Context, key, name string) (entities.Role, error)
	GetRoleByKey(ctx context.Context, key string) (entities.Role, error)
	ListRoles(ctx context.Context) ([]entities.Role, error)
	DeleteRole(ctx context.Context, id uuid.UUID) error

	AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error
	RevokeRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error
	ListUserRoles(ctx context.Context, userID uuid.UUID) ([]entities.Role, error)
}

type roleRepo struct{ q *db.Queries }

func NewRoleRepository(pool *pgxpool.Pool) RoleRepository { return &roleRepo{q: db.New(pool)} }

func (r roleRepo) CreateRole(ctx context.Context, key, name string) (entities.Role, error) {
	// TODO need to make data transfer object
	role, err := r.q.CreateRole(ctx, db.CreateRoleParams{Key: key, Name: name})
	if err != nil {
		return entities.Role{}, err
	}
	return entities.Role{
		ID:        role.ID,
		Key:       role.Key,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
	}, nil
}

func (r roleRepo) GetRoleByKey(ctx context.Context, key string) (entities.Role, error) {
	role, err := r.q.GetRoleByKey(ctx, key)
	if err != nil {
		return entities.Role{}, err
	}
	return entities.Role{
		ID:        role.ID,
		Key:       role.Key,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
	}, nil
}

func (r roleRepo) ListRoles(ctx context.Context) ([]entities.Role, error) {
	roles, err := r.q.ListRoles(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]entities.Role, 0, len(roles))
	for _, role := range roles {
		out = append(out, entities.Role{
			ID:        role.ID,
			Key:       role.Key,
			Name:      role.Name,
			CreatedAt: role.CreatedAt,
		})
	}
	return out, nil
}

func (r roleRepo) DeleteRole(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteRole(ctx, id)
}

func (r roleRepo) AssignRoleToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	return r.q.AssignRoleToUser(ctx, db.AssignRoleToUserParams{UserID: userID, RoleID: roleID})
}

func (r roleRepo) RevokeRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
	return r.q.RevokeRoleFromUser(ctx, db.RevokeRoleFromUserParams{UserID: userID, RoleID: roleID})
}

func (r roleRepo) ListUserRoles(ctx context.Context, userID uuid.UUID) ([]entities.Role, error) {
	data, err := r.q.ListUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]entities.Role, 0, len(data))
	for _, role := range data {
		out = append(out, entities.Role{
			ID:        role.ID,
			Key:       role.Key,
			Name:      role.Name,
			CreatedAt: role.CreatedAt,
		})
	}
	return out, nil
}
