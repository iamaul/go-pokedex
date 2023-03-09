package http

import (
	"github.com/iamaul/go-pokedex/config"
	"github.com/iamaul/go-pokedex/internal/auth"
	"github.com/iamaul/go-pokedex/internal/middleware"
	"github.com/iamaul/go-pokedex/internal/monster"
	"github.com/labstack/echo/v4"
)

func MonsterRoutes(monsterGroup *echo.Group, h monster.DeliveryHandlers, au auth.Usecase, cfg *config.Config, mw *middleware.MiddlewareManager) {
	monsterGroup.POST("/type", h.AddMonsterType(), mw.AuthJWTMiddleware(au, cfg), mw.RoleBasedAuthMiddleware([]string{"admin"}))
	monsterGroup.PUT("/type/:id", h.UpdateMonsterType(), mw.AuthJWTMiddleware(au, cfg), mw.RoleBasedAuthMiddleware([]string{"admin"}))
	monsterGroup.DELETE("/type/:id", h.DeleteMonsterType(), mw.AuthJWTMiddleware(au, cfg), mw.RoleBasedAuthMiddleware([]string{"admin"}))
	monsterGroup.GET("/types", h.ListMonsterType(), mw.AuthJWTMiddleware(au, cfg), mw.RoleBasedAuthMiddleware([]string{"admin"}))
	monsterGroup.GET("/type/:id", h.DetailMonsterType(), mw.AuthJWTMiddleware(au, cfg), mw.RoleBasedAuthMiddleware([]string{"admin"}))
}
