package handler

import (
	"fmt"
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


// VotePost godoc
// @Summary Vote on a post
// @Description Vote (+1/-1) on a post
// @Tags votes
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param vote body handler.VoteHandler.VotePost.true.req true "Vote data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 403 {object} map[string]string "Cannot vote on own post"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /votes [post]
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
