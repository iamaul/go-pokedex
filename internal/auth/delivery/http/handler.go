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

// Register godoc
// @Summary Register new user
// @Description register new user, returns user and token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} domain.User
// @Router /auth/new [post]
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

// Login godoc
// @Summary user authentication
// @Description returns user and token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} domain.User
// @Router /auth [post]
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

// UpdateUser godoc
// @Summary Update user
// @Description update existing user
// @Tags Auth
// @Accept json
// @Param id path int true "id"
// @Produce json
// @Success 200 {object} domain.User
// @Router /auth/{id} [put]
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

// DeleteUser godoc
// @Summary Delete user
// @Description delete existing user
// @Tags Auth
// @Accept json
// @Param id path int true "id"
// @Produce json
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpErr.RestError
// @Router /auth/{id} [delete]
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

// ListUser godoc
// @Summary Get user list
// @Description list of users
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} domain.User
// @Router /auth/user/list [get]
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

// DetailUser godoc
// @Summary Detail user
// @Description get user detail
// @Tags Auth
// @Accept json
// @Param id path int true "id"
// @Produce json
// @Success 200 {object} domain.User
// @Router /auth/{id} [get]
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

// CatchMonster godoc
// @Summary Catch monster
// @Description capture a monster
// @Tags Auth
// @Accept json
// @Param id path int true "id"
// @Produce json
// @Success 200 {object} domain.UserMonsterBody
// @Router /auth/{id} [post]
func (h *AuthHandler) CatchMonster() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		monster := &domain.UserMonsterBody{}
		if err := utils.ReadRequest(c, monster); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		if err := h.authUsecase.UserCatchMonster(c.Request().Context(), userID, monster.MonsterID); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErr.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, monster)
	}
}

// Me godoc
// @Summary Get user by id
// @Description Get current user by id
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} domain.User
// @Failure 500 {object} httpErr.RestError
// @Router /auth/me [get]
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
