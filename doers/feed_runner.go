package doers

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jmoiron/sqlx"
	"github.com/voidfiles/artarchive/config"
	"github.com/voidfiles/artarchive/debug"
	"github.com/voidfiles/artarchive/feeds"
	"github.com/voidfiles/artarchive/images"
	"github.com/voidfiles/artarchive/logging"
	"github.com/voidfiles/artarchive/pipeline"
	"github.com/voidfiles/artarchive/slides"
	"github.com/voidfiles/artarchive/storage"
)

func FeedRunner() {
	appConfig := config.NewAppConfig()
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	sss := s3.New(sess)
	logger := logging.NewLogger(false, nil)
	db, err := sqlx.Connect(appConfig.Database.Type, appConfig.Database.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch from feeds
	rssFetcher := feeds.NewFeedToSlideProducer(logger, []string{
		"https://feedbin.com/starred/Cepxc9l63Bbn0RKef9J3MQ.xml",
		"https://feedbin.com/starred/3d5um7AVLzNCL-mMtMxKeg.xml",
	})

	// Store items into a db only if we have not seen them
	itemStore := storage.MustNewItemStorage(db)
	slideStore := storage.NewDBStorageTransform(logger, itemStore)
	dropStore := storage.NewDBStorageDropTransform(logger, itemStore)

	// Resolve found slides with known versions
	slideStorage := slides.NewSlideStorage(sss, appConfig.Bucket, appConfig.Version)
	resolveTransform := slides.NewSlideResolverTransform(slideStorage)

	// Archive new images
	s3Upload := s3manager.NewUploader(sess)
	imageUploader := images.MustNewImageUploader(s3Upload, sss, appConfig.ImagePath, appConfig.Bucket)
	imageArchiver := images.NewSlideImageUploader(imageUploader)

	// Upload slides
	slideUploader := slides.NewSlideUploader(slideStorage)

	// Dump things
	debugConsumer := debug.NewDebugSlideConsumer(logger)
	pipeline := pipeline.NewPipeline(rssFetcher, debugConsumer, dropStore, resolveTransform, imageArchiver, slideUploader, slideStore)
	pipeline.Run()
}
