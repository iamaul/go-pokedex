package monster

import (
	"github.com/labstack/echo/v4"
)

type DeliveryHandlers interface {
	CreateMonsterType() echo.HandlerFunc
	UpdateMonsterType() echo.HandlerFunc
	DeleteMonsterType() echo.HandlerFunc
	ListMonsterType() echo.HandlerFunc
	DetailMonsterType() echo.HandlerFunc
	CreateMonster() echo.HandlerFunc
	UpdateMonster() echo.HandlerFunc
	DeleteMonster() echo.HandlerFunc
	ListMonster() echo.HandlerFunc
	DetailMonster() echo.HandlerFunc
	AddMonsterType() echo.HandlerFunc
}
