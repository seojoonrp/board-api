package handler

import (
	"net/http"
	"seojoonrp/board-api/internal/apperror"
	"seojoonrp/board-api/internal/domain"
	"seojoonrp/board-api/internal/service"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) UserHandler {
	return UserHandler{svc: svc}
}

func (h *UserHandler) Login(c echo.Context) error {
	var req domain.LoginRequest
	if err := c.Bind(&req); err != nil {
		return apperror.NewBadRequest("invalid request body")
	}

	resp, err := h.svc.Login(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp) // TODO: Separate Created/OK
}
