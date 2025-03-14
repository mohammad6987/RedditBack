package handler

import (
	"net/http"
	"redditBack/service"

	"github.com/gin-gonic/gin"
)

type VoteHandler struct {
	voteService *service.VoteService
}

func NewVoteHandler(voteService *service.VoteService) *VoteHandler {
	return &VoteHandler{voteService: voteService}
}

func (h *VoteHandler) VotePost(c *gin.Context) {

	usernameVal := c.Value("user_id")
	username, ok := usernameVal.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username passed from context"})
		return
	}

	var req struct {
		postID    uint `json:"postID" binding:"required"`
		voteValue int  `json:"vote" binding:"required,oneof=1 0 -1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.voteService.VotePost(c.Request.Context(), uint(req.postID), username, req.voteValue)
	if err != nil {
		switch err.Error() {
		case "post not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		case "cannot vote on your own post":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process vote"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "vote processed successfully"})
}
