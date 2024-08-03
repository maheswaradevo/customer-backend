package pkg

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Init(router *echo.Echo, db *gorm.DB, logger *zap.Logger) {
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	// app := router.Group("api/v1")
	// {

	// }
}
