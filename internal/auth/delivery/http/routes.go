package http

import (
	"github.com/iamaul/go-pokedex/internal/auth"
	"github.com/iamaul/go-pokedex/internal/middleware"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(authGroup *echo.Group, h auth.DeliveryHandlers, mw *middleware.MiddlewareManager) {
	authGroup.POST("/register", h.Register())
	authGroup.GET("/users", h.GetUserList())
	authGroup.DELETE("/:id", h.DeleteUser())
}
