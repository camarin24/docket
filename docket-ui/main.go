package main

import (
	"embed"
	"github.com/camarin24/docket"
	"github.com/camarin24/docket/adapters"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"os"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/docket.jpeg
var icon []byte

func main() {
	docketApp := docket.New(docket.Config{
		ScanInterval: time.Minute,
		Storages: []docket.FileSystemAdapter{
			adapters.NewLocalFileSystem(),
			adapters.NewS3FileSystem(adapters.S3FileSystemConfig{
				Key:        "agora.iluma.files",
				BucketName: "agora.iluma.files",
				BatchSize:  10,
			}),
			adapters.NewS3FileSystem(adapters.S3FileSystemConfig{
				Key:                          "docket-files",
				BucketName:                   "docket-files",
				MaxSizeForMetadataExtraction: 1e+7,
				BatchSize:                    10,
			})},
		DbAdapter: &docket.PostgresDb{
			DbHost:     "localhost",
			DbName:     "docket",
			DbUser:     "postgres",
			DbPassword: "postgres",
		},
		VectorDatabase: adapters.NewWeaviateVectorDb(adapters.WeaviateConfig{
			Config: weaviate.Config{
				Host:   "localhost:8080",
				Scheme: "http",
				Headers: map[string]string{
					"X-OpenAI-Api-Key": os.Getenv("OPENAI_API_KEY"),
				},
			},
		}),
		MetadataExtractor: adapters.NewTika(adapters.TikaConfig{
			ServerEndpoint: "http://localhost:9998",
		}),
	})
	docketApp.StartScanner()
	// Create an instance of the app structure
	app := NewApp(docketApp)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Docket",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour:         &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:                app.startup,
		Frameless:                runtime.GOOS != "darwin",
		EnableDefaultContextMenu: true,
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   "Docket App",
				Message: "",
				Icon:    icon,
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
