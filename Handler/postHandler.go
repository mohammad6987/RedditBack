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
