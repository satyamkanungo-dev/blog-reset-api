package controller

import "github.com/gin-gonic/gin"

func getUserIdFromMiddleware(ctx *gin.Context) (string, bool) {
	userIdInterface, exists := ctx.Get("userId")
	userIdStr := userIdInterface.(string)
	return userIdStr, exists
}
