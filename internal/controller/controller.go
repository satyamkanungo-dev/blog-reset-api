package controller

import "github.com/gin-gonic/gin"

type IUserController interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Get(ctx *gin.Context)
}
