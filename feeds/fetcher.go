package feeds

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog"
	"github.com/voidfiles/artarchive/slides"
	"golang.org/x/net/html"
)

type FeedToSlideProducer struct {
	binding slides.Binding
	feeds   []string
	logger  zerolog.Logger
}

func NewFeedToSlideProducer(logger zerolog.Logger, feeds []string) *FeedToSlideProducer {
	return &FeedToSlideProducer{
		feeds:  feeds,
		logger: logger,
	}
}

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

func SlidesFromFeeditem(item *gofeed.Item, feed gofeed.Feed) []slides.Slide {
	content := BestContent(item.Content, item.Description)

	imageURLs := FindImageUrls(content)
	resolvedImageURLs := ResolveImageURLs(item.Link, imageURLs)

	guidHash := sha256.Sum256([]byte(item.GUID))

	itemSlides := make([]slides.Slide, len(resolvedImageURLs))
	site := slides.Site{
		Title: feed.Title,
		URL:   feed.Link,
	}
	for i, imageURL := range resolvedImageURLs {
		page := slides.Page{
			Title:     item.Title,
			URL:       item.Link,
			Published: item.PublishedParsed,
			GUIDHash:  fmt.Sprintf("%x", guidHash),
		}
		urlHash := sha256.Sum256([]byte(imageURL))

		itemSlides[i] = slides.Slide{
			Site:           site,
			Page:           page,
			Content:        content,
			SourceImageURL: imageURL,
			GUIDHash:       fmt.Sprintf("%s.%x", page.GUIDHash, urlHash),
		}
	}
	return itemSlides
}

func (ff *FeedToSlideProducer) Configure(binding slides.Binding) {
	ff.binding = binding
}

func (ff *FeedToSlideProducer) Run() {
	for _, feed := range FetchFeeds(ff.feeds) {
		for _, item := range feed.Items {
			for _, slide := range SlidesFromFeeditem(item, feed) {
				ff.logger.Info().Msgf("fetcher: %v", slide.GUIDHash)
				ff.binding.Out <- slide
			}
		}
	}
	close(ff.binding.Out)
}
