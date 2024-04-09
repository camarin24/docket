package internal

type MetadataExtractor interface {
	GetFileMetadata(path string) (*map[string]string, error)
}
