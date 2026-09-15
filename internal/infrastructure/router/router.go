package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/GiaBao0510/project-templates/internal/adapter/handler"
	"github.com/GiaBao0510/project-templates/internal/adapter/middleware"
)

// NewRouter thiết lập HTTP router với Gin framework.
// Router thuộc tầng infrastructure — nơi kết nối mọi thứ lại với nhau.
func NewRouter(log *zap.Logger, userHandler *handler.UserHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger(log))
	r.Use(middleware.Recovery(log))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", userHandler.Create)
			users.GET("", userHandler.GetAll)
			users.GET("/:id", userHandler.GetByID)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}
	}

	return r
}
