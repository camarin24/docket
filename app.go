package docket

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

type App struct {
	mutex  sync.Mutex
	config Config
	db     *Db
	logger *zap.Logger
}

const (
	DefaultScanInterval = time.Minute
)

type Config struct {
	Storages        []StorageAdapter
	VectorDatabases []VectorDbAdapter
	ScanInterval    time.Duration
	DbAdapter       DbAdapter
}

const Version = "0.1"

func (app *App) Logger() *zap.Logger {
	return app.logger
}

func (app *App) Db() *Db {
	// TODO: Figure out if we need to use mutex here
	return app.db
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

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	app.logger = logger
	// TODO: Figure out if there is a best way to perform the migrations
	app.db = app.config.DbAdapter.NewDb()
	app.db.Migrate()

	return app
}
