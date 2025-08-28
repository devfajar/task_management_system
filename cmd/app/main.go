package main

import (
	"log"
	"net/http"

	"github.com/devfajar/task-management-system/configs"
	"github.com/devfajar/task-management-system/database"
	httpdelivery "github.com/devfajar/task-management-system/internal/delivery/http/handler"
	"github.com/devfajar/task-management-system/internal/delivery/http/middleware"
	"github.com/devfajar/task-management-system/internal/repository"
	"github.com/devfajar/task-management-system/internal/usecases"
	"github.com/devfajar/task-management-system/pkg/helper"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	gin.SetMode(cfg.GinMode)

	if cfg.MigrateOnStart {
		if err := database.Up(cfg.DatabaseURL); err != nil {
			log.Fatalf("migrate up: %v", err)
		}
	}

	pool, err := helper.NewPGXPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("pgxpool: %v", err)
	}
	defer pool.Close()

	// repos
	userRepo := repository.NewUserRepository(pool)
	roleRepo := repository.NewRoleRepository(pool)
	permRepo := repository.NewPermissionRepository(pool)
	authRepo := repository.NewAuthRepository(pool)

	// usecases
	userUC := usecases.NewUserUsecase(userRepo)
	roleUC := usecases.NewRoleUsecase(roleRepo)
	permUC := usecases.NewPermissionUsecase(permRepo, roleRepo)
	authUC := usecases.NewAuthUsecase(usecases.AuthConfig{
		JWTSecret:       []byte(cfg.JWTSecret),
		JWTIssuer:       cfg.JWTIssuer,
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	}, authRepo, roleRepo, permRepo)

	// handlers
	authH := &httpdelivery.AuthHandler{UC: authUC}
	userH := &httpdelivery.UserHandler{UC: userUC, Auth: middleware.Auth{}}
	roleH := &httpdelivery.RoleHandler{UC: roleUC}
	permH := &httpdelivery.PermissionHandler{UC: permUC}

	r := gin.Default()

	// root group /api/v1
	v1 := r.Group("/api/v1")

	// public (no JWT)
	authH.RegisterRoutes(v1) // /api/v1/auth/*

	// protected (JWT)
	v1Auth := v1.Group("")
	v1Auth.Use(middleware.JWT(middleware.JWTConfig{Secret: []byte(cfg.JWTSecret)}))

	// domain routes (protected)
	userH.RegisterRoutes(v1Auth) // /api/v1/users/*
	roleH.RegisterRoutes(v1Auth) // /api/v1/roles/*
	permH.RegisterRoutes(v1Auth) // /api/v1/permissions/*

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: r}
	log.Printf("HTTP listening on %s (mode=%s)", cfg.HTTPAddr, cfg.GinMode)
	log.Fatal(srv.ListenAndServe())
}
