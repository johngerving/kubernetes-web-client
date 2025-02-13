package handler

import (
	"context"

	"github.com/johngerving/kubernetes-web-client/pkg/views"
	"github.com/labstack/echo/v4"
)

// GET /
func IndexPageGET() echo.HandlerFunc {
	return func(c echo.Context) error {
		return views.Chat().Render(context.Background(), c.Response().Writer)
	}
}
