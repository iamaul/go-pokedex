package http

import (
	"net/http"

	"github.com/iamaul/go-pokedex/config"
	"github.com/iamaul/go-pokedex/internal/domain"
	"github.com/iamaul/go-pokedex/internal/monster"
	httpErr "github.com/iamaul/go-pokedex/pkg/error"
	"github.com/iamaul/go-pokedex/pkg/logger"
	"github.com/iamaul/go-pokedex/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MonsterHandler struct {
	cfg                *config.Config
	monsterTypeUsecase monster.MonsterTypeUsecase
	logger             logger.Logger
}

func NewMonsterHandler(cfg *config.Config, monsterTypeUsecase monster.MonsterTypeUsecase, log logger.Logger) monster.DeliveryHandlers {
	return &MonsterHandler{cfg: cfg, monsterTypeUsecase: monsterTypeUsecase, logger: log}
}

func (h *MonsterHandler) AddMonsterType() echo.HandlerFunc {
	return func(c echo.Context) error {
		monsterType := &domain.MonsterType{}
		if err := utils.ReadRequest(c, monsterType); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		createdMonsterType, err := h.monsterTypeUsecase.MonsterTypeCreate(c.Request().Context(), monsterType)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdMonsterType)
	}
}

func (h *MonsterHandler) UpdateMonsterType() echo.HandlerFunc {
	return func(c echo.Context) error {
		monsterTypeID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		monsterType := &domain.MonsterTypeUpdate{}
		monsterType.ID = monsterTypeID
		if err := utils.ReadRequest(c, monsterType); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		updatedMonsterType, err := h.monsterTypeUsecase.MonsterTypeUpdate(c.Request().Context(), monsterType)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedMonsterType)
	}
}

func (h *MonsterHandler) DeleteMonsterType() echo.HandlerFunc {
	return func(c echo.Context) error {
		monsterTypeID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		if err = h.monsterTypeUsecase.MonsterTypeDeletion(c.Request().Context(), monsterTypeID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *MonsterHandler) ListMonsterType() echo.HandlerFunc {
	return func(c echo.Context) error {
		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		monsterTypeList, err := h.monsterTypeUsecase.GetMonsterTypeList(c.Request().Context(), paginationQuery)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, monsterTypeList)
	}
}

func (h *MonsterHandler) DetailMonsterType() echo.HandlerFunc {
	return func(c echo.Context) error {
		monsterTypeID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		monsterType, err := h.monsterTypeUsecase.GetByID(c.Request().Context(), monsterTypeID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, monsterType)
	}
}
