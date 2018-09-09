package doers

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/voidfiles/artarchive/indexer"
	"github.com/voidfiles/artarchive/pipeline"
	"github.com/voidfiles/artarchive/scanner"
)

func RunScanner() error {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	sss := s3.New(sess)

	slideScanner := scanner.NewSlideScanner(sss, 1000, "art.rumproarious.com", "v2")

	// Dump things
	indexer := indexer.NewIndexer(sss, "http://art.rumproarious.com", "art.rumproarious.com", "art")
	pipeline := pipeline.NewPipeline(slideScanner, indexer)
	pipeline.Run()
	return nil
}
