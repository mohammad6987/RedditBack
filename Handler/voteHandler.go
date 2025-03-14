package handler

import (
	"fmt"
	"io"
	"net/http"
	"redditBack/service"

	"github.com/gin-gonic/gin"
)

type VoteHandler struct {
	voteService service.VoteService
}

func NewVoteHandler(voteService service.VoteService) VoteHandler {
	return VoteHandler{voteService: voteService}
}

func (h *VoteHandler) VotePost(c *gin.Context) {

	usernameVal := c.Value("user_id")
	username, ok := usernameVal.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username passed from context"})
		return
	}

	var req struct {
		PostID    uint `json:"postID" binding:"required"`
		VoteValue int  `json:"voteValue" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	body, _ := io.ReadAll(c.Request.Body)
	fmt.Printf("Raw Request Body: %s\n", string(body))
	fmt.Printf("req :%s %s\n", req.PostID, req.VoteValue)

	err := h.voteService.VotePost(c.Request.Context(), uint(req.PostID), username, req.VoteValue)
	if err != nil {
		fmt.Print(err.Error())
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
