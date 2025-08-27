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

	// Repository
	userRepo := repository.NewUserRepository(pool)
	roleRepo := repository.NewRoleRepository(pool)
	permissionRepo := repository.NewPermissionRepository(pool)

	// UseCase
	userUC := usecases.NewUserUsecase(userRepo)
	roleUC := usecases.NewRoleUsecase(roleRepo)
	permUC := usecases.NewPermissionUsecase(permissionRepo, roleRepo)

	// Middleware
	auth := middleware.Auth{Perms: permUC}

	r := gin.Default()
	r.Use(auth.DevAuth()) // Dev Only

	// User Handler
	userHandler := &httpdelivery.UserHandler{UC: userUC}
	roleHandler := &httpdelivery.RoleHandler{UC: roleUC}
	permHandler := &httpdelivery.PermissionHandler{UC: permUC}

	// Routes
	userHandler.RegisterRoutes(r)
	roleHandler.RegisterRoutes(r)
	permHandler.RegisterRoutes(r)

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: r}
	log.Printf("HTTP listening on %s (mode=%s)", cfg.HTTPAddr, cfg.GinMode)
	log.Fatal(srv.ListenAndServe())
}
