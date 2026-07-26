package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/Error"
	apirequest "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_request"
	apiresponse "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_response"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/service"
)

type BlogController struct {
	BlogService service.IBlogService
}

func NewBlogController(service service.IBlogService) *BlogController {
	return &BlogController{BlogService: service}
}

func (bc *BlogController) RegisterRoutes(r gin.IRouter) {
	blogs := r.Group("/blogs")
	{
		blogs.POST("", bc.Create)
		blogs.GET("", bc.GetAll)
		blogs.DELETE("", bc.DeleteMultiple)
		blogs.GET("/:id", bc.Get)
		blogs.PUT("/:id", bc.Update)
		blogs.DELETE("/:id", bc.Delete)
	}
}

func (bc *BlogController) Create(ctx *gin.Context) {
	var input apirequest.BlogRequest
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// TODO: work todo
	var userId string

	blog, err := bc.BlogService.Create(&input, userId)
	if err != nil {
		if errors.Is(err, Error.ErrAllFieldRequired) {
			ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: err.Error(),
			})
		}

		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
	}

	ctx.JSON(http.StatusCreated, apiresponse.APIResponse{
		Code:   http.StatusCreated,
		Status: "success",
		Data:   blog,
	})
}

func (bc *BlogController) Get(ctx *gin.Context) {
	// TODO:
	// get a user from middleware
	var userId string

	id := ctx.Param("id")

	blog, err := bc.BlogService.Get(id, userId)
	if err != nil {
		if errors.Is(err, Error.ErrMissingIdentifiers) {
			ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		if errors.Is(err, Error.ErrBlogNotFound) {
			ctx.JSON(http.StatusNotFound, apiresponse.APIResponse{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, apiresponse.APIResponse{
		Code:   http.StatusOK,
		Status: "success",
		Data:   blog,
	})
}

func (bc *BlogController) GetAll(ctx *gin.Context) {
	var userId string
	cursor := ctx.Query("cursor")

	blogs, err := bc.BlogService.GetAll(userId, cursor)
	if err != nil {
		if errors.Is(err, Error.ErrMissingIdentifiers) {
			ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, apiresponse.APIResponse{
		Code:   http.StatusOK,
		Status: "success",
		Data:   blogs,
	})
}

func (bc *BlogController) Update(ctx *gin.Context) {
	var input apirequest.BlogRequest
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	var userId string

	id := ctx.Param("id")

	blog, err := bc.BlogService.Update(&input, id, userId)
	if err != nil {
		if errors.Is(err, Error.ErrMissingIdentifiers) || errors.Is(err, Error.ErrAllFieldRequired) {
			ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		if errors.Is(err, Error.ErrBlogNotFound) {
			ctx.JSON(http.StatusNotFound, apiresponse.APIResponse{
				Code:    http.StatusNotFound,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, apiresponse.APIResponse{
		Code:   http.StatusOK,
		Status: "success",
		Data:   blog,
	})
}

func (bc *BlogController) Delete(ctx *gin.Context) {
	var userId string

	id := ctx.Param("id")

	err := bc.BlogService.Delete(id, userId)
	if err != nil {
		if errors.Is(err, Error.ErrMissingIdentifiers) {
			ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusNoContent, apiresponse.APIResponse{
		Code:   http.StatusNoContent,
		Status: "success",
		Data:   nil,
	})
}

func (bc *BlogController) DeleteMultiple(ctx *gin.Context) {
	var input apirequest.DeleteBlogRequest
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
			Code:    http.StatusBadRequest,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	var userId string

	deletedblogIds, err := bc.BlogService.DeleteMultiple(&input, userId)
	if err != nil {
		if errors.Is(err, Error.ErrMissingIdentifiers) || errors.Is(err, Error.ErrAllFieldRequired) {
			ctx.JSON(http.StatusBadRequest, apiresponse.APIResponse{
				Code:    http.StatusBadRequest,
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, apiresponse.APIResponse{
			Code:    http.StatusInternalServerError,
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, apiresponse.APIResponse{
		Code:   http.StatusOK,
		Status: "success",
		Data:   deletedblogIds,
	})
}
