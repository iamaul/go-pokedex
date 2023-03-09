package http

import (
	"github.com/iamaul/go-pokedex/config"
	"github.com/iamaul/go-pokedex/internal/auth"
	"github.com/iamaul/go-pokedex/internal/middleware"
	"github.com/labstack/echo/v4"
)

func AuthRoutes(authGroup *echo.Group, h auth.DeliveryHandlers, au auth.Usecase, cfg *config.Config, mw *middleware.MiddlewareManager) {
	authGroup.POST("/new", h.Register())
	authGroup.POST("", h.Login())
	authGroup.PUT("", h.UpdateUser(), mw.AuthJWTMiddleware(au, cfg))
	authGroup.DELETE("/:id", h.DeleteUser(), mw.AuthJWTMiddleware(au, cfg), mw.RoleBasedAuthMiddleware([]string{"admin"}))
	authGroup.GET("/users", h.ListUser())
	authGroup.GET("/:id", h.DetailUser(), mw.AuthJWTMiddleware(au, cfg))
	authGroup.GET("/me", h.Me(), mw.AuthJWTMiddleware(au, cfg))
}
