package docket

import "sync"

type FileSystemAdapter interface {
	Scan(app *App, wg *sync.WaitGroup)
}
