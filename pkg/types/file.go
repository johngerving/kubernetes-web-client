package types

import (
	"os"
	"time"
)

type File struct {
	Name    string
	IsDir   bool
	ModTime time.Time
}

type Files []File

// NewFilesFromDirEntries takes an os.DirEntry slice.
// It returns a list of File struct instances from the os.DirEntry slice.
func NewFilesFromDirEntries(dirEntries []os.DirEntry) (Files, error) {
	files := make([]File, len(dirEntries))

	for i, entry := range dirEntries {
		fileInfo, err := entry.Info()
		if err != nil {
			return nil, err
		}

		files[i] = File{
			Name:    entry.Name(),
			IsDir:   fileInfo.IsDir(),
			ModTime: fileInfo.ModTime(),
		}
	}

	return files, nil
}
