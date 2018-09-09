package artarchive

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/mmcdole/gofeed"
	"github.com/voidfiles/artarchive/logging"
	"golang.org/x/net/html"
)

func FetchFeed(feed_url string) gofeed.Feed {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feed_url)
	if err != nil {
		panic(err)
	}
	return *feed
}

func FetchFeeds(feed_urls []string) []gofeed.Feed {
	feeds := make([]gofeed.Feed, len(feed_urls))
	for i, url := range feed_urls {
		log.Printf("Downloading %s", url)
		feed := FetchFeed(url)
		feeds[i] = feed
	}

	return feeds
}

type FeedFetcher struct {
	slideChan chan Slide
	feeds     []string
}

func NewFeedFetcher(feeds []string, slideChan chan Slide) *FeedFetcher {
	return &FeedFetcher{
		slideChan: slideChan,
		feeds:     feeds,
	}
}

func FindImageUrls(content string) []string {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}

	queryDoc := goquery.NewDocumentFromNode(doc)
	imageURLs := make([]string, 0)

	queryDoc.Find("img").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		imageURL, exists := s.Attr("src")
		if exists {
			imageURLs = append(imageURLs, imageURL)
		}
	})

	return imageURLs
}

func BestContent(contentA, contentB string) string {
	content := contentA
	if len(contentB) > len(contentA) {
		content = contentB
	}

	return content
}

func ResolveImageURLs(baseLink string, imageURLs []string) []string {
	resolvedImageURLs := make([]string, 0)
	baseURL, err := url.Parse(baseLink)
	if err != nil {
		panic(err)
	}

	for _, imageURL := range imageURLs {
		parsedURL, err := url.Parse(imageURL)
		if err != nil {
			log.Printf("Failed to parse %v", imageURL)
			continue
		}

		resolvedURL := baseURL.ResolveReference(parsedURL)
		resolvedImageURLs = append(resolvedImageURLs, resolvedURL.String())
	}

	return resolvedImageURLs
}

func SlidesFromFeeditem(item *gofeed.Item, feed gofeed.Feed) []Slide {
	content := BestContent(item.Content, item.Description)

	imageURLs := FindImageUrls(content)
	resolvedImageURLs := ResolveImageURLs(item.Link, imageURLs)

	guidHash := sha256.Sum256([]byte(item.GUID))

	itemSlides := make([]Slide, len(resolvedImageURLs))
	site := Site{
		Title: feed.Title,
		URL:   feed.Link,
	}
	for i, imageURL := range resolvedImageURLs {
		page := Page{
			Title:     item.Title,
			URL:       item.Link,
			Published: item.PublishedParsed,
			GUIDHash:  fmt.Sprintf("%x", guidHash),
		}
		urlHash := sha256.Sum256([]byte(imageURL))

		itemSlides[i] = Slide{
			Site:           site,
			Page:           page,
			Content:        content,
			SourceImageURL: imageURL,
			GUIDHash:       fmt.Sprintf("%s.%x", page.GUIDHash, urlHash),
		}
	}
	return itemSlides
}

func (ff *FeedFetcher) Run() {
	for _, feed := range FetchFeeds(ff.feeds) {
		for _, item := range feed.Items {
			for _, slide := range SlidesFromFeeditem(item, feed) {
				ff.slideChan <- slide
			}
		}
	}
	close(ff.slideChan)
}

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

type Step interface {
	Configure(chan Slide)
}

func FeedRunner() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	sss := s3.New(sess)
	logger := logging.NewLogger(false, nil)

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
	debugConsumer := NewDebugSlideConsumer(logger, uploadToConsumer)
	pipeline := NewPipeline(rssFetcher, debugConsumer, resolveTransform, imageArchiver, slideUploader)
	pipeline.Run()
}
