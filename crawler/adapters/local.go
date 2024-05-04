package adapters

import (
	"fmt"
	"sync"
	"time"

	"github.com/camarin24/docket"
)

type LocalFileSystem struct {
}

func (ls *LocalFileSystem) Scan(app *docket.App, wg *sync.WaitGroup) {
	fmt.Println("Scanning local file system...")
	time.Sleep(time.Second * 20)
	wg.Done()
}

func NewLocalFileSystem() *LocalFileSystem {
	return &LocalFileSystem{}
}
