package main

import (
	"context"
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
	"github.com/google/uuid"
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
	secRepo := repository.NewSecurityRepository(pool)
	projectRepo := repository.NewProjectRepository(pool)
	boardRepo := repository.NewBoardRepository(pool)
	columnRepo := repository.NewColumnRepository(pool)
	taskRepo := repository.NewTaskRepository(pool)

	// usecases
	userUC := usecases.NewUserUsecase(userRepo)
	roleUC := usecases.NewRoleUsecase(roleRepo)
	permUC := usecases.NewPermissionUsecase(permRepo, roleRepo)
	authUC := usecases.NewAuthUsecase(usecases.AuthConfig{
		JWTSecret:       []byte(cfg.JWTSecret),
		JWTIssuer:       cfg.JWTIssuer,
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	}, authRepo, roleRepo, permRepo, secRepo)

	auditUC := usecases.NewAuditUsecase(secRepo)
	projectUC := usecases.NewProjectUsecase(projectRepo)
	boardUC := usecases.NewBoardUsecase(boardRepo)
	columnUC := usecases.NewColumnUsecase(columnRepo)
	taskUC := usecases.NewTaskUsecase(taskRepo)

	// handlers
	authH := &httpdelivery.AuthHandler{UC: authUC}
	userH := &httpdelivery.UserHandler{UC: userUC, Auth: middleware.Auth{}, Audit: auditUC}
	roleH := &httpdelivery.RoleHandler{UC: roleUC}
	permH := &httpdelivery.PermissionHandler{UC: permUC}
	projectH := &httpdelivery.ProjectHandler{UC: projectUC, Auth: middleware.Auth{}, Audit: auditUC}
	boardH := &httpdelivery.BoardHandler{UC: boardUC, Auth: middleware.Auth{}}
	columnH := &httpdelivery.ColumnHandler{UC: columnUC, Auth: middleware.Auth{}}
	taskH := &httpdelivery.TaskHandler{UC: taskUC, Auth: middleware.Auth{}}

	r := gin.Default()

	// root group /api/v1
	v1 := r.Group("/api/v1")

	// public (no JWT)
	authH.RegisterRoutes(v1) // /api/v1/auth/*

	// protected (JWT)
	v1Auth := v1.Group("")
	v1Auth.Use(middleware.JWT(middleware.JWTConfig{
		Secret:            []byte(cfg.JWTSecret),
		CheckTokenVersion: true,
		GetTokenVersion: func(userID uuid.UUID) (int, error) {
			tv, err := secRepo.GetTokenVersion(context.Background(), userID)
			return int(tv), err
		},
	}))

	// domain routes (protected)
	userH.RegisterRoutes(v1Auth) // /api/v1/users/*
	roleH.RegisterRoutes(v1Auth) // /api/v1/roles/*
	permH.RegisterRoutes(v1Auth) // /api/v1/permissions/*

	projectH.RegisterRoutes(v1Auth)
	boardH.RegisterRoutes(v1Auth)  // /api/v1/boards/*
	columnH.RegisterRoutes(v1Auth) // /api/v1/boards/:board_id/columns & /api/v1/columns/:id
	taskH.RegisterRoutes(v1Auth)   // /api/v1/tasks/* & /api/v1/columns/:column_id/tasks

	v1Auth.POST("/auth/logout_all", authH.LogoutAll)

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: r}
	log.Printf("HTTP listening on %s (mode=%s)", cfg.HTTPAddr, cfg.GinMode)
	log.Fatal(srv.ListenAndServe())
}
