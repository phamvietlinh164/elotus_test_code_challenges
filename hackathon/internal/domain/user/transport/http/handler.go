package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"hackathon/internal/domain/user"
	"hackathon/internal/utils"
)

type AuthHandler struct {
	svc user.Service
}

func NewAuthHandler(svc user.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func RegisterAuthRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	userRepo := user.NewRepository(db)
	svc := user.NewService(userRepo)
	h := NewAuthHandler(svc)

	auth := rg.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/revoke", h.Revoke) // ✅ public revoke route
	}
}

type registerReq struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.svc.Register(req.Username, req.Password); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondSuccess(c, gin.H{"message": "user created"})
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request")
		return
	}
	token, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, err.Error())
		return
	}
	utils.RespondSuccess(c, gin.H{"token": token})
}

type revokeReq struct {
	Token string `json:"token" binding:"required"`
}

// ✅ Public revoke endpoint
func (h *AuthHandler) Revoke(c *gin.Context) {
	var req revokeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request")
		return
	}

	claims, err := utils.GetClaimsFromJWT(req.Token)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "invalid token")
		return
	}

	if err := h.svc.RevokeTokens(claims.UserID); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to revoke token")
		return
	}

	utils.RespondSuccess(c, gin.H{"message": "token revoked"})
}
