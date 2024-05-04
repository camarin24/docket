package adapters

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/camarin24/docket"
	"github.com/camarin24/docket/pkg/types"
	"github.com/camarin24/docket/pkg/utils"
	"go.uber.org/zap"
)

type S3FileSystemConfig struct {
	Key                          string
	Workers                      int
	BatchSize                    int32
	BucketName                   string
	Region                       string
	Path                         string
	MaxSizeForMetadataExtraction int64
}

type S3FileSystem struct {
	S3FileSystemConfig
	client *s3.Client
}

const (
	DefaultBatchSize                    = 1000
	DefaultWorkers                      = 10
	DefaultRegion                       = "us-east-1"
	DefaultMaxSizeForMetadataExtraction = 1000000 * 2
)

func (s *S3FileSystem) Scan(app *docket.App, wg *sync.WaitGroup) {
	fmt.Println("Scanning S3 file system...")

	os.Mkdir(app.StoragePath(), os.ModePerm)

	existingFiles := utils.GetOnlyDocumentsNames(app.Db().GetDocumentsNameByStorageKey(s.Key))

	allScannedFiles := make([]s3Types.Object, 0)
	files := make([]types.Document, 0)
	totalExcludedFiles := 0

	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket:  aws.String(s.BucketName),
		MaxKeys: aws.Int32(s.BatchSize),
	}, func(o *s3.ListObjectsV2PaginatorOptions) {
		o.Limit = s.BatchSize
	})

	pageNum := 0
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.TODO())
		if err != nil {
			app.Logger().Error("Couldn't list objects in bucket", zap.String("bucket", s.BucketName), zap.Error(err))
			wg.Done()
			return
		}

		allScannedFiles = append(allScannedFiles, output.Contents...)
		pageNum++
	}

	for _, sf := range allScannedFiles {
		if utils.In(existingFiles, *sf.Key) {
			totalExcludedFiles++
		} else {
			document := types.Document{
				ID:         uuid.New(),
				Name:       *sf.Key,
				StorageKey: s.Key,
				// TODO: Add prefix
				OriginalPath: *sf.Key,
				Size:         *sf.Size,
			}
			files = append(files, document)
		}
	}

	app.Db().CreateDocuments(files...)

	app.Logger().Info("Total files scanned", zap.Int("count", len(allScannedFiles)))
	app.Logger().Info("Total excluded files", zap.Int("count", totalExcludedFiles))

	// TODO: Add batch processing here

	for _, doc := range files {
		if doc.Size <= s.MaxSizeForMetadataExtraction {
			s.ExtractFileMetadata(app, doc)
		}
	}

	app.Logger().Info("Finished metadata extraction")
	wg.Done()
}

func (s *S3FileSystem) ExtractFileMetadata(app *docket.App, doc types.Document) {
	resp, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &s.BucketName,
		Key:    &doc.Name,
	})

	if err != nil {
		app.Logger().Error("Error downloading the file from S3 ", zap.String("file", doc.Name), zap.Error(err))
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		app.Logger().Error("Error reading the file ", zap.String("file", doc.Name), zap.Error(err))
		return
	}

	err = os.WriteFile(path.Join(app.StoragePath(), doc.Name), body, 0644)
	if err != nil {
		app.Logger().Error("Error writing the file ", zap.String("file", doc.Name), zap.Error(err))
		return
	}
}

func NewS3FileSystem(cfg ...S3FileSystemConfig) *S3FileSystem {
	config := S3FileSystemConfig{}

	if len(cfg) > 0 {
		config = cfg[0]
	}

	if config.Workers == 0 {
		config.Workers = DefaultWorkers
	}

	if config.BatchSize == 0 {
		config.BatchSize = DefaultBatchSize
	}

	if config.Region == "" {
		config.Region = DefaultRegion
	}

	if config.MaxSizeForMetadataExtraction == 0 {
		config.MaxSizeForMetadataExtraction = DefaultMaxSizeForMetadataExtraction
	}

	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	return &S3FileSystem{
		S3FileSystemConfig: config,
		client:             s3.NewFromConfig(awsCfg),
	}
}
