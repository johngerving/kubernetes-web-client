package app

import (
	"context"
	"log/slog"

	"github.com/johngerving/kubernetes-web-client/pkg/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// registerRoutes registers both page routes and other routes with an Echo router and returns the router.
func (a *App) registerRoutes() (*echo.Echo, error) {
	e := echo.New()

	// Serve static files
	e.Static("/static", "static")

	// Set up HTTP logger
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				a.logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				a.logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "private, max-age=5")

			if err := next(c); err != nil {
				c.Error(err)
			}

			return nil
		}
	})

	// Register page routes - these serve HTML
	err := a.registerPageRoutes(e)
	if err != nil {
		return nil, err
	}

	// Register non-page routes
	e.GET("/files/contents", handler.FileContentsGET(a.config.uploadDir))
	e.GET("/files/contents/:file", handler.FileContentsGET(a.config.uploadDir))

	return e, nil
}

// registerPageRoutes registers routes for HTML pages
func (a *App) registerPageRoutes(e *echo.Echo) error {
	e.GET("/home", handler.IndexPageGET())
	e.GET("/home/:file", handler.IndexPageGET())
	e.GET("/upload", handler.UploadPageGET())

	return nil
}
