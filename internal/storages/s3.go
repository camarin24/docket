package storages

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/camarin24/docket"
	"github.com/camarin24/docket/pkg/utils"
	"go.uber.org/zap"
)

type S3FileSystemConfig struct {
	Key                          string
	Workers                      int
	BatchSize                    int
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

	existingFiles := utils.GetOnlyDocumentsNames(app.Db().GetDocumentsNameByStorageKey(s.Key))

	result, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:  aws.String(s.BucketName),
		MaxKeys: aws.Int32(int32(s.BatchSize)),
	})

	allScannedFiles := make([]s3Types.Object, 0)
	files := make([]s3Types.Object, 0)
	filesToExtractMetadata := make([]s3Types.Object, 0)

	totalExcludedFiles := 0

	if err != nil {
		app.Logger().Error("Couldn't list objects in bucket", zap.String("bucket", s.BucketName), zap.Error(err))
		return
	}

	if *result.IsTruncated {
		// TODO: Paginated response, deal with that
	} else {
		allScannedFiles = append(allScannedFiles, result.Contents...)
	}

	for _, sf := range allScannedFiles {
		if utils.In(existingFiles, *sf.Key) {
			totalExcludedFiles++
		} else {
			if *sf.Size <= s.MaxSizeForMetadataExtraction {
				filesToExtractMetadata = append(filesToExtractMetadata, sf)
			} else {
				files = append(files, sf)
			}
		}
	}

	

	//app.Logger().Sugar().Info(*result.NextContinuationToken)

	wg.Done()
}

func (s *S3FileSystem) ExtractFileMetadata(path string) {

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

	awsConfig, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	return &S3FileSystem{
		S3FileSystemConfig: config,
		client:             s3.NewFromConfig(awsConfig),
	}
}
