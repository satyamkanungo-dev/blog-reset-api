package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/middleware"
)

type IUserController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Get(ctx *gin.Context)
	GetToken(ctx *gin.Context)
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
	RegisterRoutes(r *gin.RouterGroup, middleware middleware.IMiddleware)
}
