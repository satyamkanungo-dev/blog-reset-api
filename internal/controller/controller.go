package controller

import "github.com/gin-gonic/gin"

type IUserController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Get(ctx *gin.Context)
}

type IBlogController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Get(ctx *gin.Context)
	Delete(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	DeleteMultiple(ctx *gin.Context)
}

type IController interface {
	RegisterRoutes(r gin.IRouter, authMiddleware gin.HandlerFunc)
}
