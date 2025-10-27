package utils

import (
	"hackathon/internal/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespondSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, common.BaseApiResponse[any]{
		Success: true,
		Data:    data,
	})
}

func RespondError(c *gin.Context, code int, msg string) {
	c.JSON(code, common.BaseApiResponse[any]{
		HttpRequestStatus: code,
		Success:           false,
		Message:           msg,
		Data:              nil,
	})
}
