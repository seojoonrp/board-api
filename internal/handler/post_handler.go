package handler

import (
	"net/http"
	"seojoonrp/board-api/internal/apperror"
	"seojoonrp/board-api/internal/dto"
	"seojoonrp/board-api/internal/middleware"
	"seojoonrp/board-api/internal/service"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	svc service.PostService
}

func NewPostHandler(svc service.PostService) PostHandler {
	return PostHandler{svc: svc}
}

func (h *PostHandler) Create(c echo.Context) error {
	userID, err := middleware.UserIDFrom(c)
	if err != nil {
		return err
	}

	var req dto.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return apperror.NewBadRequest("invalid request body")
	}

	post, err := h.svc.Create(c.Request().Context(), req, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()

	if author := c.QueryParam("author"); author != "" {
		posts, err := h.svc.GetByUserID(ctx, author)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, posts)
	}

	posts, err := h.svc.GetAll(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) Get(c echo.Context) error {
	id := c.Param("id")

	post, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, post)
}
