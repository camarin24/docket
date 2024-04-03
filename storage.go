package docket

import "sync"

type StorageAdapter interface {
	Scan(app *App, wg *sync.WaitGroup)
}
