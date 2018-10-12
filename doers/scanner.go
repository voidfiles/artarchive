package doers

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/voidfiles/artarchive/config"
	"github.com/voidfiles/artarchive/debug"
	"github.com/voidfiles/artarchive/indexer"
	"github.com/voidfiles/artarchive/logging"
	"github.com/voidfiles/artarchive/pipeline"
	"github.com/voidfiles/artarchive/storage"
)

func RunScanner(renderer string) error {
	appConfig := config.NewAppConfig()
	logger := logging.NewLogger(false, nil)
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	sss := s3.New(sess)

	db, err := sqlx.Connect(appConfig.Database.Type, appConfig.Database.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	itemStore := storage.MustNewItemStorage(db)
	slideScanner := storage.NewDBStorageProducer(logger, itemStore)
	// Dump things
	var render indexer.Renderer
	if renderer == "slideshow" {
		render = indexer.NewSlideshowRenderer(
			logger,
			sss,
			// https://s3.us-west-2.amazonaws.com/art.rumproarious.com/v2/slides/054fd38619b0b3752d79ca3d66182796.a98f339d491468f9372651f7f90e072ea3add478dea2075c8ba9181b038eef38/data.json
			fmt.Sprintf("https://s3.us-west-2.amazonaws.com/%s", appConfig.Bucket),
			appConfig.Bucket)
	} else if renderer == "blog" {
		render = indexer.NewBlogRenderer(
			logger,
			sss,
			5,
			// https://s3.us-west-2.amazonaws.com/art.rumproarious.com/v2/slides/054fd38619b0b3752d79ca3d66182796.a98f339d491468f9372651f7f90e072ea3add478dea2075c8ba9181b038eef38/data.json
			fmt.Sprintf("https://s3.us-west-2.amazonaws.com/%s", appConfig.Bucket),
			appConfig.Bucket)
	}

	// limitConsumer := debug.NewDebugLimitSlideConsumer(logger, 100)
	slideStore := storage.NewDBUpdateTransform(logger, itemStore)
	debugConsumer := debug.NewDebugSlideConsumer(logger)
	pipeline := pipeline.NewPipeline(slideScanner, debugConsumer, render, slideStore)
	pipeline.Run()
	return nil
}
