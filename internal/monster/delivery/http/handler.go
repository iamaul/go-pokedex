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
	monsterUsecase     monster.MonsterUsecase
	logger             logger.Logger
}

func NewMonsterHandler(cfg *config.Config, monsterTypeUsecase monster.MonsterTypeUsecase, monsterUsecase monster.MonsterUsecase, log logger.Logger) monster.DeliveryHandlers {
	return &MonsterHandler{cfg: cfg, monsterTypeUsecase: monsterTypeUsecase, monsterUsecase: monsterUsecase, logger: log}
}

// CreateMonsterType godoc
// @Summary Create a new monster type
// @Description returns monster type
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} domain.MonsterType
// @Router /monster/type [post]
func (h *MonsterHandler) CreateMonsterType() echo.HandlerFunc {
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

// UpdateMonsterType godoc
// @Summary Update monster type
// @Description update existing monster type
// @Tags Auth
// @Accept json
// @Param id path int true "id"
// @Produce json
// @Success 200 {object} domain.MonsterTypeUpdate
// @Router /monster/type/{id} [put]
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

// DeleteMonsterType godoc
// @Summary Delete monster type
// @Description delete existing monster type
// @Tags Auth
// @Accept json
// @Param id path int true "id"
// @Produce json
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpErr.RestError
// @Router /monster/type/{id} [put]
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

// ListMonsterType godoc
// @Summary Get monster type list
// @Description list of monster types
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} domain.MonsterTypeList
// @Router /monster/type/list [get]
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

// DetailMonsterType godoc
// @Summary Detail monster type
// @Description Get monster type detail
// @Tags Auth
// @Accept json
// @Param id path int true "id"
// @Produce json
// @Success 200 {object} domain.MonsterType
// @Router /monster/type/{id} [get]
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

// CreateMonster godoc
// @Summary Create a new monster
// @Description returns monster
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} domain.Monster
// @Router /monster [post]
func (h *MonsterHandler) CreateMonster() echo.HandlerFunc {
	return func(c echo.Context) error {
		monster := &domain.Monster{}
		if err := utils.ReadRequest(c, monster); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		createdMonster, err := h.monsterUsecase.MonsterCreate(c.Request().Context(), monster)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdMonster)
	}
}

// UpdateMonster godoc
// @Summary Update monster
// @Description update existing monster
// @Tags Auth
// @Accept json
// @Param id path int true "id"
// @Produce json
// @Success 200 {object} domain.MonsterUpdate
// @Router /monster/{id} [put]
func (h *MonsterHandler) UpdateMonster() echo.HandlerFunc {
	return func(c echo.Context) error {
		monsterID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		monster := &domain.MonsterUpdate{}
		monster.ID = monsterID
		if err := utils.ReadRequest(c, monster); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		updatedMonster, err := h.monsterUsecase.MonsterUpdate(c.Request().Context(), monster)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedMonster)
	}
}

// DeleteMonster godoc
// @Summary Delete monster
// @Description delete existing monster
// @Tags Auth
// @Accept json
// @Param id path int true "id"
// @Produce json
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpErr.RestError
// @Router /monster/{id} [put]
func (h *MonsterHandler) DeleteMonster() echo.HandlerFunc {
	return func(c echo.Context) error {
		monsterID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		if err = h.monsterUsecase.MonsterDeletion(c.Request().Context(), monsterID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

// ListMonster godoc
// @Summary Get monster list
// @Description list of monster
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} domain.MonsterList
// @Router /monster/list [get]
func (h *MonsterHandler) ListMonster() echo.HandlerFunc {
	return func(c echo.Context) error {
		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		monsterList, err := h.monsterUsecase.GetMonsterList(c.Request().Context(), paginationQuery)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, monsterList)
	}
}

// DetailMonster godoc
// @Summary Detail monster
// @Description Get monster detail
// @Tags Auth
// @Accept json
// @Param id path int true "id"
// @Produce json
// @Success 200 {object} domain.Monster
// @Router /monster/{id} [get]
func (h *MonsterHandler) DetailMonster() echo.HandlerFunc {
	return func(c echo.Context) error {
		monsterID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		monster, err := h.monsterUsecase.GetByID(c.Request().Context(), monsterID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, monster)
	}
}

// AddMonsterType godoc
// @Summary Add monster type
// @Description add monster type
// @Tags Auth
// @Accept json
// @Param id path int true "id"
// @Produce json
// @Success 200 {object} domain.MonsterTypeBody
// @Router /monster/{id} [post]
func (h *MonsterHandler) AddMonsterType() echo.HandlerFunc {
	return func(c echo.Context) error {
		monsterID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		monsterType := &domain.MonsterTypeBody{}
		if err := utils.ReadRequest(c, monsterType); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		if err := h.monsterUsecase.AttachMonsterType(c.Request().Context(), monsterID, monsterType.MonsterTypeID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, monsterType)
	}
}
