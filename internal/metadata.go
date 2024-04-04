package internal

import (
	"fmt"

	"github.com/barasher/go-exiftool"
)

func GetFileMatada(path string) (*map[string]interface{}, error) {
	et, err := exiftool.NewExiftool()
	if err != nil {
		return nil, err
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata(path)

	metadata := make(map[string]interface{}, 0)

	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			fmt.Printf("Error concerning %v: %v\n", fileInfo.File, fileInfo.Err)
			continue
		}

		for k, v := range fileInfo.Fields {
			metadata[k] = v
		}
	}

	return &metadata, nil
}
