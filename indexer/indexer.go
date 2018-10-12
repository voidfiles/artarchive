package indexer

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Interface holds methods we use to interact with s3
type IndexerS3Interface interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}
