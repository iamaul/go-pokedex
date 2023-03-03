package middleware

import (
	"github.com/iamaul/go-pokedex/config"
	"github.com/iamaul/go-pokedex/internal/auth"
	"github.com/iamaul/go-pokedex/pkg/logger"
)

type MiddlewareManager struct {
	authUsecase auth.Usecase
	cfg         *config.Config
	origins     []string
	logger      logger.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(authUsecase auth.Usecase, cfg *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{authUsecase: authUsecase, cfg: cfg, origins: origins, logger: logger}
}
