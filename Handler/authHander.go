package handler

import (
	"net/http"
	"redditBack/model"
	"redditBack/service"
	"redditBack/utility"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return AuthHandler{authService: authService}
}


// SignUp godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body handler.AuthHandler.SignUp.true.req true "User registration data"
// @Success 201 {object} map[string]interface{} "Successfully created user"
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tempUser := &model.User{
		Username:     req.Username,
		PasswordHash: req.Password,
		Email:        req.Email,
	}

	err := h.authService.Register(c.Request.Context(), tempUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := utility.GenerateToken(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user":  req.Username,
	})
}


// Login godoc
// @Summary Authenticate user
// @Description Login with username and password to get JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body handler.AuthHandler.Login.true.req true "Login credentials"
// @Success 200 {object} map[string]interface{} "Successfully logged in"
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := utility.GenerateToken(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}


// SignOut godoc
// @Summary Logout user
// @Description Invalidate user's JWT token
// @Tags authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string "Successfully logged out"
// @Failure 400 {object} map[string]string "Missing authorization token"
// @Failure 500 {object} map[string]string "Failed to invalidate token"
// @Router /signout [post]
func (h *AuthHandler) SignOut(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No authorization token provided"})
		return
	}


	err := h.authService.InvalidateToken(c.Request.Context(), tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully signed out"})
}
