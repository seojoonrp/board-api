package server

import (
	"seojoonrp/board-api/internal/handler"
	"seojoonrp/board-api/internal/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	postHandler handler.PostHandler,
	userHandler handler.UserHandler,
	jwtSecret string,
) {
	e.GET("/posts", postHandler.GetAll)
	e.GET("/posts/:id", postHandler.Get)
	e.POST("/users", userHandler.Login)

	auth := e.Group("", middleware.JWTAuth(jwtSecret))
	auth.POST("/posts", postHandler.Create)
}
