package utils

import "github.com/camarin24/docket/pkg/types"

func GetOnlyDocumentsNames(d []types.Document) []string {
	names := make([]string, 0)
	for _, d := range d {
		names = append(names, d.Name)
	}
	return names
}
