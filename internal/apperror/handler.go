package apperror

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Handler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, map[string]string{"error": appErr.Message})
		return
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		c.JSON(echoErr.Code, map[string]any{"error": echoErr.Message})
		return
	}

	c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}
