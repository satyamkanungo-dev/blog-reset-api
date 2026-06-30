package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Iserver interface {
	Serve() error
	ShutDown() error
}

type Server struct {
	server *http.Server
}

func InitializeServer(router *gin.Engine, portAdd string) *Server {
	return &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%s", portAdd),
			Handler:      router,
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Serve() error {
	return s.server.ListenAndServe()
}

func (s *Server) ShutDown() error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %W", err)
	}

	return nil
}
