package main

import (
	"fmt"
	"time"

	"github.com/camarin24/docket"
	"github.com/camarin24/docket/internal/storages"
)

func main() {
	app := docket.New(docket.Config{
		ScanInterval: time.Second * 10,
		Storages: []docket.StorageAdapter{
			storages.NewLocalFileSystem(),
			storages.NewS3FileSystem(storages.S3FileSystemConfig{
				Key:        "agora.iluma.files",
				BucketName: "agora.iluma.files",
				BatchSize: 10,
			})},
	})
	app.StartScanner()
	time.Sleep(time.Hour)
	fmt.Println(app)
}
