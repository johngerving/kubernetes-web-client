package handler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/johngerving/kubernetes-web-client/pkg/templates"
	"github.com/johngerving/kubernetes-web-client/pkg/types"
	"github.com/johngerving/kubernetes-web-client/pkg/views"
	"github.com/labstack/echo/v4"
)

// GET /home/:file
func IndexPageGET() echo.HandlerFunc {
	return func(c echo.Context) error {
		file := c.Param("file")
		contentParam := ""
		if file != "" {
			contentParam = "/" + file
		}
		return views.Index(contentParam).Render(context.Background(), c.Response().Writer)
	}
}

// GET /files/contents/:file
func FileContentsGET(uploadDir string) echo.HandlerFunc {
	return func(c echo.Context) error {
		file := c.Param("file")

		contents, err := getFileContents(uploadDir, file)
		if err != nil {
			return err
		}

		return templates.FileList(contents).Render(context.Background(), c.Response().Writer)
	}
}

// getFileContents returns a list of Files in a given directory and an
// error if unsuccessful.
func getFileContents(baseDir, file string) (types.Files, error) {
	fmt.Println(filepath.Join(baseDir, file))
	dirEntries, err := os.ReadDir(filepath.Join(baseDir, file))

	if err != nil {
		return nil, err
	}

	files, err := types.NewFilesFromDirEntries(dirEntries)
	if err != nil {
		return nil, err
	}

	return files, nil
}
