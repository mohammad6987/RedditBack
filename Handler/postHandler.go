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
	fmt.Printf("username in post :%s", username)
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
