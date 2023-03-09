package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"

	"github.com/iamaul/go-pokedex/config"
	"github.com/iamaul/go-pokedex/internal/auth"
	"github.com/iamaul/go-pokedex/internal/domain"
	httpErr "github.com/iamaul/go-pokedex/pkg/error"
	"github.com/iamaul/go-pokedex/pkg/utils"
)

// Authentication based JWT
func (mw *MiddlewareManager) AuthJWTMiddleware(authUsecase auth.Usecase, cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			bearerHeader := c.Request().Header.Get("Authorization")

			mw.logger.Infof("auth middleware bearerHeader %s", bearerHeader)

			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					mw.logger.Error("auth middleware", zap.String("headerParts", "len(headerParts) != 2"))
					return c.JSON(http.StatusUnauthorized, httpErr.NewUnauthorizedError(httpErr.Unauthorized))
				}

				tokenString := headerParts[1]

				if err := mw.validateJWTToken(tokenString, authUsecase, c, cfg); err != nil {
					mw.logger.Error("middleware validateJWTToken", zap.String("headerJWT", err.Error()))
					return c.JSON(http.StatusUnauthorized, httpErr.NewUnauthorizedError(httpErr.Unauthorized))
				}

				return next(c)
			}

			cookie, err := c.Cookie("jwt-token")
			if err != nil {
				mw.logger.Errorf("c.Cookie", err.Error())
				return c.JSON(http.StatusUnauthorized, httpErr.NewUnauthorizedError(httpErr.Unauthorized))
			}

			if err = mw.validateJWTToken(cookie.Value, authUsecase, c, cfg); err != nil {
				mw.logger.Errorf("validateJWTToken", err.Error())
				return c.JSON(http.StatusUnauthorized, httpErr.NewUnauthorizedError(httpErr.Unauthorized))
			}
			return next(c)
		}
	}
}

// Admin permission role
func (mw *MiddlewareManager) AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*domain.User)
		if !ok || *user.Role != "admin" {
			return c.JSON(http.StatusForbidden, httpErr.NewUnauthorizedError(httpErr.PermissionDenied))
		}
		return next(c)
	}
}

// Role based auth middleware using ctx user
func (mw *MiddlewareManager) AuthenticatedOrAdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*domain.User)
			if !ok {
				mw.logger.Errorf("Error c.Get(user) RequestID: %s, Error: %s,", utils.GetRequestID(c), "invalid user ctx")
				return c.JSON(http.StatusUnauthorized, httpErr.NewUnauthorizedError(httpErr.Unauthorized))
			}

			if *user.Role == "admin" {
				return next(c)
			}

			if user.ID.String() != c.Param("id") {
				mw.logger.Errorf("Error c.Get(user) RequestID: %s, UserID: %s, Error: %s,",
					utils.GetRequestID(c),
					user.ID.String(),
					"invalid user ctx",
				)
				return c.JSON(http.StatusForbidden, httpErr.NewForbiddenError(httpErr.Forbidden))
			}

			return next(c)
		}
	}
}

// Role based auth middleware using ctx user
func (mw *MiddlewareManager) RoleBasedAuthMiddleware(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*domain.User)
			if !ok {
				mw.logger.Errorf("Error c.Get(user) RequestID: %s, UserID: %s, ERROR: %s,",
					utils.GetRequestID(c),
					user.ID.String(),
					"invalid user ctx",
				)
				return c.JSON(http.StatusUnauthorized, httpErr.NewUnauthorizedError(httpErr.Unauthorized))
			}

			for _, role := range roles {
				if role == *user.Role {
					return next(c)
				}
			}

			mw.logger.Errorf("Error c.Get(user) RequestID: %s, UserID: %s, ERROR: %s,",
				utils.GetRequestID(c),
				user.ID.String(),
				"invalid user ctx",
			)

			return c.JSON(http.StatusForbidden, httpErr.NewForbiddenError(httpErr.PermissionDenied))
		}
	}
}

func (mw *MiddlewareManager) validateJWTToken(tokenString string, authUsecase auth.Usecase, c echo.Context, cfg *config.Config) error {
	if tokenString == "" {
		return httpErr.InvalidJWTToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(cfg.Server.JwtSecretKey)
		return secret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return httpErr.InvalidJWTToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return httpErr.InvalidJWTClaims
		}

		userUUID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return err
		}

		u, err := authUsecase.GetByID(c.Request().Context(), userUUID)
		if err != nil {
			return err
		}

		c.Set("user", u)

		ctx := context.WithValue(c.Request().Context(), utils.UserCtxKey{}, u)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}
