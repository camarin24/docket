package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/camarin24/docket/pkg/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"gorm.io/datatypes"
	"os"
	"strings"
	"testing"
)

func getClient() *WeaviateVectorDb {
	return NewWeaviateVectorDb(WeaviateConfig{
		Config: weaviate.Config{
			Host:   "localhost:8080",
			Scheme: "http",
			Headers: map[string]string{
				"X-OpenAI-Api-Key": os.Getenv("OPENAI_API_KEY"),
			},
		},
	})
}

func TestWeaviateClient(t *testing.T) {
	client := getClient()

	schema, err := client.client.Schema().Getter().Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", schema)
	assert.NotNil(t, schema, "weaviate")
}

func TestIngestDocument(t *testing.T) {
	client := getClient()

	data := map[string]string{
		"Author": "camarin24", "CreateDate": "2023:12:04 11:20:13-05:00",
	}
	metadata, err := json.Marshal(data)
	assert.NoError(t, err)
	document := types.Document{
		ID:           uuid.New(),
		Name:         "test document.pdf",
		StorageKey:   "s3.key",
		OriginalPath: "s3.key/test.document.pdf",
		Size:         100000,
		MetaData:     datatypes.JSON(strings.Replace(string(metadata), "\n", "", -1)),
	}

	err = client.IngestDocument(document)
	assert.NoError(t, err)
}

func TestGetDocument(t *testing.T) {
	client := getClient()

	docs, err := client.RetrieveDocuments("vitaminas")

	for _, doc := range *docs {
		t.Logf("----------------%s %s", doc.ID, doc.Content)
	}
	assert.NoError(t, err)
}
