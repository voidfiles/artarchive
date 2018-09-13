package doers

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/voidfiles/artarchive/config"
	"github.com/voidfiles/artarchive/indexer"
	"github.com/voidfiles/artarchive/pipeline"
	"github.com/voidfiles/artarchive/scanner"
)

func RunScanner() error {
	appConfig := config.NewAppConfig()
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	sss := s3.New(sess)

	slideScanner := scanner.NewSlideScanner(sss, 1000, appConfig.Bucket, appConfig.Version)

	// Dump things
	indexer := indexer.NewIndexer(
		sss,
		fmt.Sprintf("http://%s", appConfig.Bucket),
		appConfig.Bucket,
		"art")
	pipeline := pipeline.NewPipeline(slideScanner, indexer)
	pipeline.Run()
	return nil
}
