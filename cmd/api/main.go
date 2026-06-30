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
	"github.com/satyamkanungo-dev/blog-rest-api/internal/database"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/server"
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

	// eg:
	router := gin.Default()
	router.GET("api/v1", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Blog API 👍",
		})
	})

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
