package handler

import (
	"context"

	"github.com/johngerving/kubernetes-web-client/pkg/views"
	"github.com/labstack/echo/v4"
)

// GET /upload
func UploadPageGET() echo.HandlerFunc {
	return func(c echo.Context) error {
		return views.Upload().Render(context.Background(), c.Response().Writer)
	}
}
