package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/GiaBao0510/project-templates/pkg/response"
)

// Recovery middleware bắt panic và trả về lỗi 500 thay vì crash server.
func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic recovered", zap.Any("error", err))
				response.Error(c, http.StatusInternalServerError, "Internal Server Error", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}
