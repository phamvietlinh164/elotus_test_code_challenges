package upload

import (
	"net/http"

	"hackathon/internal/domain/upload"
	"hackathon/internal/middlewares"
	"hackathon/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UploadHandler struct {
	svc upload.Service
}

func NewUploadHandler(svc upload.Service) *UploadHandler {
	return &UploadHandler{svc: svc}
}

func RegisterUploadRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	uploadRepo := upload.NewRepository(db)
	svc := upload.NewService(uploadRepo)
	h := NewUploadHandler(svc)
	rg.POST("/upload", middlewares.AuthMiddleware(), h.Upload)
}

func (h *UploadHandler) Upload(c *gin.Context) {

	userIDStr, existed := c.Get("userID")
	if !existed {
		utils.RespondError(c, http.StatusBadRequest, "missing file field 'data'")
		return
	}
	userID := userIDStr.(uint)

	file, header, err := c.Request.FormFile("data")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "missing file field 'data'")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	meta, err := h.svc.UploadImage(userID, file, header, contentType, ip, userAgent)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondSuccess(c, gin.H{
		"message": "upload success",
		"file":    meta,
	})
}
