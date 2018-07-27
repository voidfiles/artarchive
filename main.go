package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type FeedToSlideProducer struct {
	slideChan chan Slide
	feeds     []string
}

func NewFeedToSlideProducer(feeds []string, slideChan chan Slide) *FeedToSlideProducer {
	return &FeedToSlideProducer{
		slideChan: slideChan,
		feeds:     feeds,
	}
}

func (ff *FeedToSlideProducer) Run() {
	for _, feed := range FetchFeeds(ff.feeds) {
		for _, item := range feed.Items {
			for _, slide := range SlidesFromFeeditem(item, feed) {
				ff.slideChan <- slide
			}
		}
	}
	close(ff.slideChan)
}

func runner() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	sss := s3.New(sess)

	feedToResolve := make(chan Slide, 0)
	resolveToImageArchive := make(chan Slide, 0)
	archiveToUpload := make(chan Slide, 0)
	uploadToConsumer := make(chan Slide, 0)

	// Fetch from feeds
	rssFetcher := NewFeedToSlideProducer([]string{
		"https://feedbin.com/starred/Cepxc9l63Bbn0RKef9J3MQ.xml",
		"https://feedbin.com/starred/3d5um7AVLzNCL-mMtMxKeg.xml",
	}, feedToResolve)

	// Resolve found slides with known versions
	slideStorage := NewSlideStorage(sss, "art.rumproarious.com", "v2")
	resolveTransform := NewSlideResolverTransform(feedToResolve, resolveToImageArchive, slideStorage)

	// Archive new images
	s3Upload := s3manager.NewUploader(sess)
	imageUploader := MustNewImageUploader(s3Upload, sss, "images", "art.rumproarious.com")
	imageArchiver := NewSlideImageUploader(resolveToImageArchive, archiveToUpload, imageUploader)

	// Upload slides
	slideUploader := NewSlideUploader(archiveToUpload, uploadToConsumer, slideStorage)

	// Dump things
	debugConsumer := NewDebugSlideConsumer(uploadToConsumer)
	pipeline := NewPipeline(rssFetcher, debugConsumer, resolveTransform, imageArchiver, slideUploader)
	pipeline.Run()
}

func HandleRequest() (string, error) {
	runner()
	return fmt.Sprintf("We made it"), nil
}

func scanner() error {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	sss := s3.New(sess)

	feedToResolve := make(chan Slide, 16)

	slideScanner := NewSlideScanner(sss, feedToResolve, 1000, "art.rumproarious.com", "v2")

	// Dump things
	indexer := NewIndexer(sss, feedToResolve, "http://art.rumproarious.com", "art.rumproarious.com", "art")
	pipeline := NewPipeline(slideScanner, indexer)
	pipeline.Run()
	return nil
}

func main() {
	scan := flag.Bool("scan", false, "Scan all the slides")
	flag.Parse()
	log.SetOutput(os.Stdout)
	if lambdaIs := os.Getenv("LAMBDA"); lambdaIs != "" {
		if lambdaIs == "finder" {
			lambda.Start(func() (string, error) {
				runner()
				return fmt.Sprintf("We made it"), nil
			})
		} else if lambdaIs == "indexer" {
			lambda.Start(func() (string, error) {
				scanner()
				return fmt.Sprintf("We made it"), nil
			})
		}
	}

	if *scan {
		scanner()
	} else {
		runner()
	}

}
