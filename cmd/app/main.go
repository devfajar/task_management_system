package main

import (
	"log"
	"net/http"

	"github.com/devfajar/task-management-system/configs"
	"github.com/devfajar/task-management-system/database"
	httpdelivery "github.com/devfajar/task-management-system/internal/delivery/http/handler"
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

	userRepo := repository.NewUserRepository(pool)
	userUC := usecases.NewUserUsecase(userRepo)

	r := gin.Default()
	userHandler := &httpdelivery.UserHandler{UC: userUC}
	userHandler.RegisterRoutes(r)

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: r}
	log.Printf("HTTP listening on %s (mode=%s)", cfg.HTTPAddr, cfg.GinMode)
	log.Fatal(srv.ListenAndServe())
}
