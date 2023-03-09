package monster

import (
	"github.com/labstack/echo/v4"
)

type DeliveryHandlers interface {
	AddMonsterType() echo.HandlerFunc
	UpdateMonsterType() echo.HandlerFunc
	DeleteMonsterType() echo.HandlerFunc
	ListMonsterType() echo.HandlerFunc
	DetailMonsterType() echo.HandlerFunc
}
