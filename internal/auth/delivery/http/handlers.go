package http

import (
	"net/http"

	"github.com/iamaul/go-pokedex/config"
	"github.com/iamaul/go-pokedex/internal/auth"
	"github.com/iamaul/go-pokedex/internal/domain"
	httpErr "github.com/iamaul/go-pokedex/pkg/error"
	"github.com/iamaul/go-pokedex/pkg/logger"
	"github.com/iamaul/go-pokedex/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	cfg         *config.Config
	authUsecase auth.Usecase
	logger      logger.Logger
}

func NewAuthHandler(cfg *config.Config, authUsecase auth.Usecase, log logger.Logger) auth.DeliveryHandlers {
	return &AuthHandler{cfg: cfg, authUsecase: authUsecase, logger: log}
}

func (h *AuthHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := &domain.User{}
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		createdUser, err := h.authUsecase.UserRegistration(c.Request().Context(), user)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, createdUser)
	}
}

func (h *AuthHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		login := &domain.UserLogin{}
		if err := utils.ReadRequest(c, login); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		userWithToken, err := h.authUsecase.UserAuthentication(c.Request().Context(), &domain.User{
			Username: login.Username,
			Password: login.Password,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, userWithToken)
	}
}

func (h *AuthHandler) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		me, ok := c.Get("user").(*domain.User)
		if !ok {
			utils.LogResponseError(c, h.logger, httpErr.NewUnauthorizedError(httpErr.Unauthorized))
			return utils.ErrResponseWithLog(c, h.logger, httpErr.NewUnauthorizedError(httpErr.Unauthorized))
		}
		user := &domain.UserUpdate{}
		user.ID = me.ID
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		updatedUser, err := h.authUsecase.UserUpdate(c.Request().Context(), user)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}

func (h *AuthHandler) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		if err = h.authUsecase.UserDeletion(c.Request().Context(), userID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.NoContent(http.StatusOK)
	}
}

func (h *AuthHandler) ListUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		usersList, err := h.authUsecase.UserList(c.Request().Context(), paginationQuery)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, usersList)
	}
}

func (h *AuthHandler) DetailUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		user, err := h.authUsecase.GetByID(c.Request().Context(), userID)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

func (h *AuthHandler) Me() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*domain.User)
		if !ok {
			utils.LogResponseError(c, h.logger, httpErr.NewUnauthorizedError(httpErr.Unauthorized))
			return utils.ErrResponseWithLog(c, h.logger, httpErr.NewUnauthorizedError(httpErr.Unauthorized))
		}

		return c.JSON(http.StatusOK, user)
	}
}
