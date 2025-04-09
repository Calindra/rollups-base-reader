package transaction

import (
	"github.com/labstack/echo/v4"
)

// Register the Paio API to echo
func Register(e *echo.Echo, paioServerAPI ServerInterface) {
	RegisterHandlers(e, paioServerAPI)
}
