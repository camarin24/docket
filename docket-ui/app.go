package main

import (
	"context"
	"github.com/camarin24/docket"
	"github.com/camarin24/docket/pkg/types"
)

// App struct
type App struct {
	ctx context.Context
	da  *docket.App
}

// NewApp creates a new App application struct
func NewApp(da *docket.App) *App {
	return &App{
		da: da,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetRecentDocuments() []types.Document {
	return a.da.Db().GetRecentDocuments()
}

func (a *App) QueryDocuments(query string) []types.Document {
	vectorDocs, err := a.da.VectorDb().RetrieveDocuments(query)
	dbDocs := a.da.Db().QueryDocuments(query)
	if err != nil {
		return make([]types.Document, 0)
	}

	dbDocs = append(dbDocs, *vectorDocs...)
	return dbDocs
}
