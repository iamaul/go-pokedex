package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag/example/basic/docs"

	authHttp "github.com/iamaul/go-pokedex/internal/auth/delivery/http"
	authRepository "github.com/iamaul/go-pokedex/internal/auth/repository"
	authUseCase "github.com/iamaul/go-pokedex/internal/auth/usecase"
	monsterHttp "github.com/iamaul/go-pokedex/internal/monster/delivery/http"
	monsterRepository "github.com/iamaul/go-pokedex/internal/monster/repository"
	monsterUseCase "github.com/iamaul/go-pokedex/internal/monster/usecase"

	apiMiddlewares "github.com/iamaul/go-pokedex/internal/middleware"
	"github.com/iamaul/go-pokedex/pkg/csrf"
	"github.com/iamaul/go-pokedex/pkg/utils"
)

// Map Server Handlers
func (s *Server) MapRouteHandlers(e *echo.Echo) error {
	// Repositories
	authRepo := authRepository.NewAuthRepo(s.db)
	monsterTypeRepo := monsterRepository.NewMonsterTypeRepo(s.db)
	monsterRepo := monsterRepository.NewMonsterRepo(s.db)

	// Usecases
	authUsecase := authUseCase.NewAuthUsecase(s.cfg, authRepo, monsterRepo, s.logger)
	monsterTypeUsecase := monsterUseCase.NewMonsterTypeUsecase(s.cfg, monsterTypeRepo, s.logger)
	monsterUsecase := monsterUseCase.NewMonsterUsecase(s.cfg, monsterRepo, monsterTypeRepo, s.logger)

	// Handlers
	authHandler := authHttp.NewAuthHandler(s.cfg, authUsecase, s.logger)
	monsterTypeHandler := monsterHttp.NewMonsterHandler(s.cfg, monsterTypeUsecase, monsterUsecase, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(authUsecase, s.cfg, []string{"*"}, s.logger)

	e.Use(mw.RequestLoggerMiddleware)

	docs.SwaggerInfo.Title = "go-pokedex API Docs"
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	if s.cfg.Server.Debug {
		e.Use(mw.DebugMiddleware)
	}

	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	authGroup := v1.Group("/auth")
	monsterGroup := v1.Group("/monster")

	authHttp.AuthRoutes(authGroup, authHandler, authUsecase, s.cfg, mw)
	monsterHttp.MonsterRoutes(monsterGroup, monsterTypeHandler, authUsecase, s.cfg, mw)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check requestId: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
