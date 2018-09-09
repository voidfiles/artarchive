package doers

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/voidfiles/artarchive/debug"
	"github.com/voidfiles/artarchive/feeds"
	"github.com/voidfiles/artarchive/images"
	"github.com/voidfiles/artarchive/logging"
	"github.com/voidfiles/artarchive/pipeline"
	"github.com/voidfiles/artarchive/slides"
)

func FeedRunner() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	sss := s3.New(sess)
	logger := logging.NewLogger(false, nil)

	// Fetch from feeds
	rssFetcher := feeds.NewFeedToSlideProducer(logger, []string{
		"https://feedbin.com/starred/Cepxc9l63Bbn0RKef9J3MQ.xml",
		"https://feedbin.com/starred/3d5um7AVLzNCL-mMtMxKeg.xml",
	})

	// // Resolve found slides with known versions
	slideStorage := slides.NewSlideStorage(sss, "art.rumproarious.com", "v2")
	resolveTransform := slides.NewSlideResolverTransform(slideStorage)
	//
	// Archive new images
	s3Upload := s3manager.NewUploader(sess)
	imageUploader := images.MustNewImageUploader(s3Upload, sss, "images", "art.rumproarious.com")
	imageArchiver := images.NewSlideImageUploader(imageUploader)

	// // Upload slides
	slideUploader := slides.NewSlideUploader(slideStorage)

	// // Dump things
	debugConsumer := debug.NewDebugSlideConsumer(logger)
	pipeline := pipeline.NewPipeline(rssFetcher, debugConsumer, resolveTransform, imageArchiver, slideUploader)
	pipeline.Run()
}
