package adapters

import (
	"io"
	"io/fs"
	"os"
	"path"
	"sync"
	"time"

	"github.com/camarin24/docket"
	"github.com/camarin24/docket/pkg/types"
	"github.com/camarin24/docket/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type LocalFileSystemConfig struct {
	Path                         string
	Key                          string
	MaxSizeForMetadataExtraction int64
}

type LocalFileSystem struct {
	LocalFileSystemConfig
}

func (ls *LocalFileSystem) Scan(app *docket.App, wg *sync.WaitGroup) {
	app.Logger().Info("Scanning local file system...")

	allScannedFiles := make([]string, 0)
	existingFiles := utils.GetOnlyDocumentsNames(app.Db().GetDocumentsNameByStorageKey(ls.Key))
	files := make([]types.Document, 0)

	totalExcludedFiles := 0

	fays := os.DirFS(path.Join(app.FileSystemPath(), ls.Path))
	err := fs.WalkDir(fays, ".", func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			app.Logger().Info("Scanning file", zap.String("file", p))
			allScannedFiles = append(allScannedFiles, p)
		}
		return nil
	})

	if err != nil {
		app.Logger().Error("Error scanning local file system", zap.Error(err))
		wg.Done()
		return
	}

	for _, sf := range allScannedFiles {
		stats, err := os.Stat(path.Join(app.FileSystemPath(), ls.Path, sf))
		size := int64(0)
		name := sf
		if err == nil {
			app.Logger().Info("File found", zap.String("file", sf))
			size = stats.Size()
			name = stats.Name()
			app.Logger().Info("File size", zap.Int64("size", size))
			app.Logger().Info("File name", zap.String("name", name))
		} else {
			app.Logger().Error("Error getting file stat", zap.Error(err))
		}

		if utils.In(existingFiles, name) {
			totalExcludedFiles++
		} else {
			document := types.Document{
				ID:           uuid.New(),
				Name:         name,
				StorageKey:   ls.Key,
				OriginalPath: path.Join(ls.Path, sf),
				Size:         size,
			}

			files = append(files, document)
		}
	}

	app.Db().CreateDocuments(files...)

	app.Logger().Info("Total files scanned", zap.Int("count", len(allScannedFiles)))
	app.Logger().Info("Total excluded files", zap.Int("count", totalExcludedFiles))

	for _, doc := range files {
		if doc.Size <= ls.MaxSizeForMetadataExtraction {
			stats, err := os.Stat(path.Join(app.FileSystemPath(), doc.OriginalPath))
			if err != nil {
				app.Logger().Error("Error getting file stat", zap.Error(err))
				continue
			}

			if !stats.Mode().IsRegular() {
				app.Logger().Info("Skipping non-regular file", zap.String("file", path.Join(app.FileSystemPath(), doc.OriginalPath)))
				continue
			}

			source, err := os.Open(path.Join(app.FileSystemPath(), doc.OriginalPath))
			if err != nil {
				app.Logger().Error("Error opening file", zap.Error(err))
				continue
			}

			defer source.Close()
			destination, err := os.Create(path.Join(app.StoragePath(), stats.Name()))
			if err != nil {
				app.Logger().Error("Error creating file", zap.Error(err))
				continue
			}

			defer destination.Close()

			_, err = io.Copy(destination, source)
			if err != nil {
				app.Logger().Error("Error copying file", zap.Error(err))
				continue
			}
		}
	}

	time.Sleep(time.Second * 20)
	wg.Done()
}

func NewLocalFileSystem(cfg LocalFileSystemConfig) *LocalFileSystem {
	return &LocalFileSystem{
		LocalFileSystemConfig: cfg,
	}
}
