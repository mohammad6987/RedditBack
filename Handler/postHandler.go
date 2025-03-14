package handler

import (
	"fmt"
	"net/http"
	"redditBack/model"
	"redditBack/service"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) PostHandler {
	return PostHandler{postService: postService}
}


// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post with title and content
// @Tags posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post body handler.PostHandler.CreatePost.true.req true "Post creation data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	usernameVal := c.Value("user_id")
	username, ok := usernameVal.(string)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username passed from context"})
		return
	}
	var req struct {
		Title   string `json:"Title" binding:"required,min=6"`
		Context string `json:"Context" binding:"required,min=12"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &model.Post{
		Title:   req.Title,
		Content: req.Title,
		UserID:  0,
	}

	err := h.postService.CreateNewPost(c.Request.Context(), post, username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})

}


// EditPost godoc
// @Summary Update a post
// @Description Update an existing post's title or content
// @Tags posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post body handler.PostHandler.EditPost.true.req true "Post update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 403 {object} map[string]string "Unauthorized to edit"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts [put]
func (h *PostHandler) EditPost(c *gin.Context) {
	usernameVal := c.Value("user_id")
	username, ok := usernameVal.(string)
	fmt.Printf("username in post :%s", username)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username passed from context"})
		return
	}
	var req struct {
		ID      uint   `json:"ID" binding:"required"`
		Title   string `json:"Title" binding:"omitempty,min=6"`
		Content string `json:"Content" binding:"omitempty,min=12"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPost := &model.Post{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
	}

	nextErr := h.postService.EditPost(c.Request.Context(), updatedPost, username)
	if nextErr != nil {
		switch nextErr.Error() {
		case "post not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		case "unauthorized to edit post":
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized to edit this post"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update post"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "post updated successfully",
		"post":    updatedPost,
	})

}


// RemovePost godoc
// @Summary Delete a post
// @Description Delete an existing post
// @Tags posts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param post body handler.PostHandler.RemovePost.true.req true "Post deletion data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 403 {object} map[string]string "Unauthorized to delete"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts [delete]
func (h *PostHandler) RemovePost(c *gin.Context) {
	usernameVal := c.Value("user_id")
	username, ok := usernameVal.(string)
	fmt.Printf("username in post :%s", username)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username passed from context"})
		return
	}
	var req struct {
		ID uint `json:"ID" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPost := &model.Post{
		ID: req.ID,
	}

	nextErr := h.postService.RemovePost(c.Request.Context(), updatedPost, username)
	if nextErr != nil {
		switch nextErr.Error() {
		case "post not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		case "unauthorized to edit post":
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not authorized to edit this post"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove post"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "post removed successfully",
	})

}
