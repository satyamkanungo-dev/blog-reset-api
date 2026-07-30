package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/middleware"
)

type MainController struct {
	controllers []IController
}

func NewMainController(controllers ...IController) *MainController {
	return &MainController{
		controllers: controllers,
	}
}

func (mc *MainController) RegisterRoutes(r gin.IRouter, middleware middleware.IAuthMiddleware) {
	for _, c := range mc.controllers {
		c.RegisterRoutes(r, middleware)
	}
}
