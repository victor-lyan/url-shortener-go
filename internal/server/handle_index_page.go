package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func HandleIndexPage(c echo.Context) error {
	return c.Render(http.StatusOK, "index", struct{}{})
}
