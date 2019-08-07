package handler

import (
	"net/http"

	"github.com/ironzhang/practice/x-cloud/x-account/account"
	"github.com/ironzhang/practice/x-cloud/x-account/types"
	"github.com/labstack/echo"
)

type Handler struct {
	account *account.Manager
}

func (h *Handler) Register(c echo.Context) error {
	var req types.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := h.account.Register(req.Type, req.Name, req.Password); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) Login(c echo.Context) error {
	var req types.LoginRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	token, err := h.account.Login(req.Type, req.Name, req.Password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, token)
}

func New(account *account.Manager) *Handler {
	return &Handler{account: account}
}

func Register(e *echo.Echo, h *Handler) {
	e.POST("/accounts", h.Register)
	e.POST("/accounts/login", h.Login)
}
