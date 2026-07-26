package controller

import "github.com/gin-gonic/gin"

type MainController struct {
	controllers []IController
}

func NewMainController(controllers ...IController) *MainController {
	return &MainController{
		controllers: controllers,
	}
}

func (mc *MainController) RegisterRoutes(r *gin.Engine) {
	for _, c := range mc.controllers {
		c.RegisterRoutes(r)
	}
}
