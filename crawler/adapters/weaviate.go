package adapters

import (
	"context"
	"encoding/json"
	"github.com/camarin24/docket/pkg/types"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/data/replication"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"strings"
)

type WeaviateVectorDb struct {
	client    *weaviate.Client
	className string
}

var DefaultWeaviateCollectionName = "DocketDocuments"

type WeaviateConfig struct {
	weaviate.Config
}

func (w *WeaviateVectorDb) IngestDocument(document types.Document) error {

	if len(document.Content) > 8191 {
		document.Content = document.Content[:5000]
	}
	_, err := w.client.Data().Creator().
		WithClassName(w.className).
		WithProperties(document).
		WithID(document.ID.String()).
		WithConsistencyLevel(replication.ConsistencyLevel.ALL).
		Do(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (w *WeaviateVectorDb) UpdateDocument(document types.Document) {

}

func (w *WeaviateVectorDb) RetrieveDocuments(query string) (*[]types.Document, error) {
	fields := []graphql.Field{
		{Name: "name"},
		{Name: "storageKey"},
		{Name: "originalPath"},
		{Name: "size"},
		{Name: "content"},
		{Name: "metaData"},
		{Name: "iD"},
		{
			Name: "_additional", Fields: []graphql.Field{
				{Name: "distance"},
			},
		},
	}

	nearTextObject := w.client.GraphQL().NearTextArgBuilder().WithConcepts(strings.Split(query, " ")).WithDistance(1)

	//hybridTextObject := w.client.GraphQL().HybridArgumentBuilder().WithQuery(query)

	//TODO: Add limit as config value
	resp, err := w.client.GraphQL().Get().
		WithClassName(w.className).
		//WithHybrid(hybridTextObject).
		WithNearText(nearTextObject).
		WithFields(fields...).
		WithLimit(10).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	documents := make([]types.Document, 0)

	// TODO: DocketDocuments is the class name so it can´t be changed
	type ResponseDocument struct {
		DocketDocuments []types.Document
	}

	for _, d := range resp.Data {
		var document ResponseDocument
		bytes, err := json.Marshal(d)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, &document)
		if err != nil {
			return nil, err
		}
		documents = append(documents, document.DocketDocuments...)
	}

	return &documents, nil
}

func (w *WeaviateVectorDb) createClass() {

	schema, err := w.client.Schema().ClassGetter().WithClassName(w.className).Do(context.Background())

	// TODO: Add class creation settings as config
	if err != nil || schema == nil {

		err = w.client.Schema().ClassCreator().WithClass(&models.Class{
			Class:        w.className,
			Description:  "A class to storage all scanned documents",
			ModuleConfig: map[string]string{"model": "text-embedding-3-small", "dimensions": "1536", "type": "text"},
			Vectorizer:   "text2vec-openai",
		}).Do(context.Background())

		if err != nil {
			panic(err)
		}
	}
}

func NewWeaviateVectorDb(config WeaviateConfig) *WeaviateVectorDb {
	client, err := weaviate.NewClient(config.Config)
	if err != nil {
		panic(err)
	}

	w := &WeaviateVectorDb{
		client:    client,
		className: DefaultWeaviateCollectionName,
	}
	w.createClass()
	return w
}
