package docket

import (
	"time"

	"go.uber.org/zap"
)

type App struct {
	config Config
	db     *Db
	logger *zap.Logger
}

const (
	DefaultScanInterval = time.Minute
)

type Config struct {
	storagePath  string
	flsPath      string
	Storages     []FileSystemAdapter
	ScanInterval time.Duration
	DbAdapter    DbAdapter
}

const Version = "0.1"

func (app *App) Logger() *zap.Logger {
	return app.logger
}

func (app *App) Db() *Db {
	// TODO: Figure out if we need to use mutex here
	return app.db
}

func (app *App) StoragePath() string {
	return app.config.storagePath
}

func (app *App) FileSystemPath() string {
	return app.config.flsPath
}

func New(config ...Config) *App {
	app := &App{
		config: Config{},
	}

	if len(config) > 0 {
		app.config = config[0]
	}

	if app.config.ScanInterval == 0 {
		app.config.ScanInterval = DefaultScanInterval
	}

	if app.config.DbAdapter == nil {
		app.config.DbAdapter = &SqliteDb{}
	}

	app.config.storagePath = "docket-files"
	app.config.flsPath = "/usr/src/app/lfs"
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	app.logger = logger
	// TODO: Figure out if there is a best way to perform the migrations
	app.db = app.config.DbAdapter.NewDb()
	app.db.Migrate()

	return app
}
