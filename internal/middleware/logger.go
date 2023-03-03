package middleware

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/iamaul/go-pokedex/pkg/utils"
)

func (mw *MiddlewareManager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start).String()
		requestID := utils.GetRequestID(ctx)

		mw.logger.Infof("RequestID: %s, method: %s, URI: %s, status: %v, size: %v, time: %s",
			requestID, req.Method, req.URL, status, size, s,
		)
		return err
	}
}
