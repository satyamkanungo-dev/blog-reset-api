package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/config"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/controller"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/database"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/middleware"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/repository"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/server"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err = cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	// database
	pool, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	defer pool.Close()

	// repositories
	db := repository.NewRepository(pool)
	userRepo := repository.NewUserRepository(db)
	blogRepo := repository.NewBlogRepository(db)

	// services
	userServices := service.NewUserService(userRepo)
	blogServices := service.NewBlogService(blogRepo)
	authServices := service.NewAuthService(cfg.Secret)

	// controllers
	userController := controller.NewUserController(userServices, authServices)
	blogController := controller.NewBlogController(blogServices)
	mainController := controller.NewMainController(userController, blogController)

	// middleware
	authMiddleware := middleware.AuthMiddleware(authServices)

	// router
	router := gin.Default()
	v1 := router.Group("api/v1")
	{
		router.GET("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Blog API 👍",
			})
		})

		mainController.RegisterRoutes(v1, authMiddleware)
	}

	workDone := make(chan os.Signal, 1)

	svr := server.InitializeServer(router, cfg.Port)
	go func() {
		if err := svr.Serve(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	signal.Notify(workDone, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-workDone
	slog.Info("shutting down server...")

	if err := svr.ShutDown(); err != nil {
		slog.Error("graceful shutdown failed", "error:", err)
	}
}
