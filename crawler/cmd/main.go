package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/camarin24/docket/adapters"

	"github.com/camarin24/docket"
)

func main() {
	app := docket.New(docket.Config{
		ScanInterval: time.Minute * 60,
		Storages: []docket.FileSystemAdapter{
			adapters.NewLocalFileSystem(adapters.LocalFileSystemConfig{
				Path: "Docs",
				Key: "local.files",
				MaxSizeForMetadataExtraction: 1000000000 * 2,
			}),
			// adapters.NewS3FileSystem(adapters.S3FileSystemConfig{
			// 	Key:        "agora.iluma.files",
			// 	BucketName: "agora.iluma.files",
			// 	BatchSize:  10,
			// })},
		},
		DbAdapter: &docket.PostgresDb{
			DbHost:     "postgres",
			DbName:     "docket",
			DbUser:     "postgres",
			DbPassword: "postgres",
		},
	})
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	app.Logger().Info("Starting crawler...")
	app.StartScanner()
	app.Logger().Info("Crawler started")
	<-done
}
