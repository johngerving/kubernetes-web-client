package handler

import (
	"context"
	"os"

	"github.com/johngerving/kubernetes-web-client/pkg/templates"
	"github.com/johngerving/kubernetes-web-client/pkg/types"
	"github.com/johngerving/kubernetes-web-client/pkg/views"
	"github.com/labstack/echo/v4"
)

// GET /
func IndexPageGET() echo.HandlerFunc {
	return func(c echo.Context) error {
		return views.Index().Render(context.Background(), c.Response().Writer)
	}
}

// GET /files
func FilesGET(uploadDir string) echo.HandlerFunc {
	return func(c echo.Context) error {
		files, err := getFiles(uploadDir)
		if err != nil {
			return err
		}

		return templates.FileList(files).Render(context.Background(), c.Response().Writer)
	}
}

// getFiles returns a list of Files in a given directory and an
// error if unsuccessful.
func getFiles(dir string) (types.Files, error) {
	dirEntries, err := os.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	files, err := types.NewFilesFromDirEntries(dirEntries)
	if err != nil {
		return nil, err
	}

	return files, nil
}
