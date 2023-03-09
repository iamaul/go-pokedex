package auth

import (
	"github.com/labstack/echo/v4"
)

type DeliveryHandlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	UpdateUser() echo.HandlerFunc
	DeleteUser() echo.HandlerFunc
	ListUser() echo.HandlerFunc
	DetailUser() echo.HandlerFunc
	// CatchMonster() echo.HandlerFunc
	Me() echo.HandlerFunc
}
