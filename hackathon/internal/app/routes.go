package app

import (
	uploadhttp "hackathon/internal/domain/upload/transport/http"
	"hackathon/internal/domain/user/transport/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerRoutes(router *gin.Engine, db *gorm.DB) {
	api := router.Group("/api")
	{
		// auth routes
		http.RegisterAuthRoutes(api, db)

		uploadhttp.RegisterUploadRoutes(api, db)
	}
}
