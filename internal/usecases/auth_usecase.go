package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/devfajar/task-management-system/internal/repository"
	"github.com/devfajar/task-management-system/pkg/helper"
	"github.com/devfajar/task-management-system/pkg/utils"
	"github.com/google/uuid"
)

var (
	ErrAuthInvalid = errors.New("invalid email or password")
	ErrAuthBlocked = errors.New("account disabled or deleted")
	ErrRTInvalid   = errors.New("invalid refresh token")
	ErrRTExpired   = errors.New("refresh token expired")
	ErrRTRevoked   = errors.New("refresh token revoked")
)

type AuthConfig struct {
	JWTSecret       []byte
	JWTIssuer       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type AuthUsecase interface {
	Login(ctx context.Context, email, passwordPlain, userAgent, ip string) (accessToken, refreshToken string, err error)
	Refresh(ctx context.Context, refreshToken, userAgent, ip string) (newAccess, newRefresh string, err error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, userID uuid.UUID) error
}
type authUC struct {
	cfg   AuthConfig
	auth  repository.AuthRepository
	roles repository.RoleRepository
	perms repository.PermissionRepository
	sec   repository.SecurityRepository
}

func NewAuthUsecase(cfg AuthConfig, auth repository.AuthRepository, roles repository.RoleRepository, perms repository.PermissionRepository, sec repository.SecurityRepository) AuthUsecase {
	return &authUC{cfg: cfg, auth: auth, roles: roles, perms: perms, sec: sec}
}

func (u *authUC) Login(ctx context.Context, email, passwordPlain, ua, ip string) (string, string, error) {
	uid, hash, isActive, deleted, err := u.auth.GetUserCredentialsByEmail(ctx, email)
	if err != nil || !helper.Verify(hash, passwordPlain) {
		return "", "", ErrAuthInvalid
	}
	if !isActive || deleted {
		return "", "", ErrAuthBlocked
	}

	roleKeys, permKeys, err := u.loadRolePerm(ctx, uid)
	if err != nil {
		return "", "", err
	}

	tv, err := u.sec.GetTokenVersion(ctx, uid)
	if err != nil {
		return "", "", err
	}

	access, err := utils.SignHS256(u.cfg.JWTSecret, u.cfg.JWTIssuer, uid, roleKeys, permKeys, int(tv), u.cfg.AccessTokenTTL)
	if err != nil {
		return "", "", err
	}

	rt, err := helper.GenerateOpaque(32)
	if err != nil {
		return "", "", err
	}
	rth := helper.SHA256Hex(rt)
	if err := u.auth.SaveRefreshToken(ctx, uid, rth, time.Now().Add(u.cfg.RefreshTokenTTL), ua, ip); err != nil {
		return "", "", err
	}
	return access, rt, nil
}

func (u *authUC) Refresh(ctx context.Context, refreshToken, ua, ip string) (string, string, error) {
	if refreshToken == "" {
		return "", "", ErrRTInvalid
	}
	hash := helper.SHA256Hex(refreshToken)

	id, uid, exp, revokedAt, err := u.auth.FindRefreshToken(ctx, hash)
	if err != nil {
		return "", "", ErrRTInvalid
	}
	if revokedAt != nil {
		return "", "", ErrRTRevoked
	}
	if time.Now().After(exp) {
		return "", "", ErrRTExpired
	}

	// rotate old RT
	if err := u.auth.RevokeRefreshToken(ctx, id); err != nil {
		return "", "", err
	}

	roleKeys, permKeys, err := u.loadRolePerm(ctx, uid)
	if err != nil {
		return "", "", err
	}

	tv, err := u.sec.GetTokenVersion(ctx, uid)
	if err != nil {
		return "", "", err
	}

	newAccess, err := utils.SignHS256(u.cfg.JWTSecret, u.cfg.JWTIssuer, uid, roleKeys, permKeys, int(tv), u.cfg.AccessTokenTTL)
	if err != nil {
		return "", "", err
	}

	newRT, err := helper.GenerateOpaque(32)
	if err != nil {
		return "", "", err
	}
	newHash := helper.SHA256Hex(newRT)
	if err := u.auth.SaveRefreshToken(ctx, uid, newHash, time.Now().Add(u.cfg.RefreshTokenTTL), ua, ip); err != nil {
		return "", "", err
	}
	return newAccess, newRT, nil
}

func (u *authUC) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return ErrRTInvalid
	}
	hash := helper.SHA256Hex(refreshToken)
	id, _, _, revokedAt, err := u.auth.FindRefreshToken(ctx, hash)
	if err != nil {
		return ErrRTInvalid
	}
	if revokedAt != nil {
		return nil
	}
	return u.auth.RevokeRefreshToken(ctx, id)
}

func (u *authUC) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	return u.sec.BumpTokenVersion(ctx, userID)
}

// private function
func (u *authUC) loadRolePerm(ctx context.Context, uid uuid.UUID) (roleKeys, permKeys []string, err error) {
	// roles
	rs, err := u.roles.ListUserRoles(ctx, uid)
	if err != nil {
		return nil, nil, err
	}
	roleKeys = make([]string, 0, len(rs))
	for _, r := range rs {
		roleKeys = append(roleKeys, r.Key)
	}
	// perms
	permKeys, err = u.perms.ListUserPermissionKeys(ctx, uid)
	if err != nil {
		return nil, nil, err
	}
	return roleKeys, permKeys, nil
}
