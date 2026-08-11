package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/Error"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/middleware"
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

func (bc *BlogController) RegisterRoutes(r *gin.RouterGroup, middleware middleware.IMiddleware) {
	blogs := r.Group("/blogs").Use(middleware.AuthMiddleware())
	{
		blogs.POST("", bc.Create)
		blogs.GET("", bc.GetAll)
		blogs.DELETE("", bc.DeleteMultiple)
		blogs.GET("/:id", bc.Get)
		blogs.PUT("/:id", bc.Update)
		blogs.DELETE("/:id", bc.Delete)
	}
}

// Create 			godoc
// @Summary			a new blog
// @Description		Create a new blog
// @Tags 			blogs
// @Accept 			json
// @Produce 		json
// @Security     	BearerAuth
// @Param			input	body 		apirequest.BlogRequest	true 	"blog details"
// @Success 		201 	{object} 	apiresponse.APIResponse{data=models.Blog} "blog created successfully"
// @Failure			400		{object}	apiresponse.APIResponse	"Invalid input, missing required fields"
// @Failure			401 	{object}	apiresponse.APIResponse			"Unauthorized"
// @Failure 		500 	{object} 	apiresponse.APIResponse	"Internal server error"
// @Router 			/blogs  [Post]
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

	userId, _ := getUserIdFromMiddleware(ctx)

	blog, err := bc.BlogService.Create(&input, userId)
	if err != nil {
		if errors.Is(err, Error.ErrAllFieldRequired) {
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

	ctx.JSON(http.StatusCreated, apiresponse.APIResponse{
		Code:   http.StatusCreated,
		Status: "success",
		Data:   blog,
	})
}

// Get 				godoc
// @Summary			Get a blog
// @Description		Get a blog by id
// @Tags 			blogs
// @Accept 			json
// @Produce 		json
// @Security     	BearerAuth
// @Param 			id  	path	string	true 	"blog id"
// @Success 		200 	{object} 	apiresponse.APIResponse{data=models.Blog}
// @Failure			400		{object}	apiresponse.APIResponse	"Invalid input, missing required fields"
// @Failure 		401		{object}	apiresponse.APIResponse "Unauthorized"
// @Failure 		404		{object}	apiresponse.APIResponse	"blog not found"
// @Failure 		500 	{object} 	apiresponse.APIResponse	"Internal server error"
// @Router 			/blogs/{id}  [Get]
func (bc *BlogController) Get(ctx *gin.Context) {
	userId, _ := getUserIdFromMiddleware(ctx)

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

// GetAll 			godoc
// @Summary			Get all blogs
// @Description		Get all blogs
// @Tags 			blogs
// @Accept 			json
// @Produce 		json
// @Security     	BearerAuth
// @Param        	cursor  query     string  false  "Pagination cursor from the previous response"
// @Success 		200 	{object} 	apiresponse.APIResponse{data=models.Blog}
// @Failure			400		{object}	apiresponse.APIResponse	"Invalid input, missing required fields"
// @Failure 		401		{object}	apiresponse.APIResponse "Unauthorized"
// @Failure 		500 	{object} 	apiresponse.APIResponse	"Internal server error"
// @Router 			/blogs  [Get]
func (bc *BlogController) GetAll(ctx *gin.Context) {
	userId, _ := getUserIdFromMiddleware(ctx)
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

// Update 			godoc
// @Summary 		Update users details
// @Description		Get a updated user
// @Tags 			blogs
// @Accept			json
// @Produce 		json
// @Security	 	BearerAuth
// @Param			input 	body 		apirequest.UpdateUserRequest true "Update user details"
// @Success 		200		{object}	apiresponse.APIResponse{data=models.User}	 "User details updated successfully"
// @Failure			400		{object}	apiresponse.APIResponse	"Invalid input, missing required fields, or password too short"
// @Failure 		401		{object}	apiresponse.APIResponse "Unauthorized"
// @Failure			404		{object}	apiresponse.APIResponse	"Invalid blog"
// @Failure 		500 	{object} 	apiresponse.APIResponse	"Internal server error"
// @Router			/blogs/{id} [Put]
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

	userId, _ := getUserIdFromMiddleware(ctx)

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

// Delete 			godoc
// @Summary			Delete a blog
// @Description		Delete a blog by id
// @Tags 			blogs
// @Accept 			json
// @Produce 		json
// @Security     	BearerAuth
// @Param 			id  	path	string	true 	"blog id"
// @Success 		204 	{object} 	apiresponse.APIResponse	"Delete successfully"
// @Failure			400		{object}	apiresponse.APIResponse	"Invalid input, missing required fields"
// @Failure 		401		{object}	apiresponse.APIResponse "Unauthorized"
// @Failure 		404		{object}	apiresponse.APIResponse	"blog not found"
// @Failure 		500 	{object} 	apiresponse.APIResponse	"Internal server error"
// @Router 			/blogs/{id}  [Delete]
func (bc *BlogController) Delete(ctx *gin.Context) {
	userId, _ := getUserIdFromMiddleware(ctx)

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

// DeleteMultiple 	godoc
// @Summary			Delete multiple blogs
// @Description		Deletes an array of blog IDs belonging to the authenticated user. Returns lists of succeeded and failed deletions.
// @Tags 			blogs
// @Accept 			json
// @Produce 		json
// @Security     	BearerAuth
// @Param 			input  	body		apirequest.DeleteBlogRequest true 	"List of Todo IDs to delete"
// @Success 		207 	{object} 	apiresponse.APIResponse{data=apiresponse.BulkDeleteResponse}	"Multi-Status response indicating successful and failed deletions""
// @Failure			400		{object}	apiresponse.APIResponse	"Invalid input, missing required fields"
// @Failure 		401		{object}	apiresponse.APIResponse "Unauthorized"
// @Failure 		500 	{object} 	apiresponse.APIResponse	"Internal server error"
// @Router 			/blogs  [Delete]
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

	userId, _ := getUserIdFromMiddleware(ctx)

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

	ctx.JSON(http.StatusMultiStatus, apiresponse.APIResponse{
		Code:   http.StatusMultiStatus,
		Status: "success",
		Data:   deletedblogIds,
	})
}
