package artarchive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Site contains identifying info about what parent site an image was found
type Site struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Page contains identifying info about where an image was found
type Page struct {
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	Published *time.Time `json:"published"`
	GUIDHash  string
}

// ImageInfo is important information about a single image
type ImageInfo struct {
	URL         string `json:"url,omitempty"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Filename    string `json:"filename,omitempty"`
}

// ArtistInfo is important information about the artist who made the image
type ArtistInfo struct {
	Name         string   `json:"name,omitempty"`
	ArtsyURL     string   `json:"artsy_url,omitempty"`
	WikipediaURL string   `json:"wikipedia_url,omitempty"`
	WebsiteURL   string   `json:"website_url,omitempty"`
	FeedURL      string   `json:"feed_url,omitempty"`
	InstagramURL string   `json:"instagram_url,omitempty"`
	TwitterURL   string   `json:"twitter_url,omitempty"`
	Description  string   `json:"description,omitempty"`
	Feeds        []string `json:"feeds,omitempty"`
	Sites        []string `json:"sites,omitempty"`
}

// WorkInfo is important information about the work
type WorkInfo struct {
	Name string `json:"name,omitempty"`
}

// Slide bundles together important information about a find
type Slide struct {
	Site           Site          `json:"site,omitempty"`
	Page           Page          `json:"page,omitempty"`
	Content        string        `json:"content,omitempty"`
	GUIDHash       string        `json:"guid_hash,omitempty"`
	SourceImageURL string        `json:"source_image_url,omitempty"`
	ArchivedImage  *ImageInfo    `json:"archived_image,omitempty"`
	ArtistInfo     *ArtistInfo   `json:"artist_info,omitempty"`
	ArtistsInfo    []*ArtistInfo `json:"artists,omitempty"`
	WorkInfo       *WorkInfo     `json:"work_info,omitempty"`
}

// S3Interface holds methods we use to interact with s3
type S3Interface interface {
	HeadObject(input *s3.HeadObjectInput) (*s3.HeadObjectOutput, error)
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

// SlideStorage mediates our relationship with s3
type SlideStorage struct {
	sss    S3Interface
	bucket string
	prefix string
}

// NewSlideStorage returns a new instance of SlideStorage
func NewSlideStorage(sss S3Interface, bucket, prefix string) *SlideStorage {
	return &SlideStorage{
		sss:    sss,
		bucket: bucket,
		prefix: prefix,
	}
}

func buildKey(prefix string, slide Slide) string {
	if prefix != "" {
		return fmt.Sprintf("%s/slides/%s/data.json", prefix, slide.GUIDHash)
	}

	return fmt.Sprintf("slides/%s/data.json", slide.GUIDHash)
}

// Resolve will resolve a slide if that slide was previously uploaded
func (ss *SlideStorage) Resolve(slide Slide) Slide {
	resp, err := ss.sss.GetObject(&s3.GetObjectInput{
		Key:    aws.String(buildKey(ss.prefix, slide)),
		Bucket: aws.String(ss.bucket),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != s3.ErrCodeNoSuchKey {
				log.Printf("Error resolving slide: %v key: %s bucket: %s", err, buildKey(ss.prefix, slide), ss.bucket)
			}
		} else {
			log.Printf("Error resolving slide: %v key: %s bucket: %s", err, buildKey(ss.prefix, slide), ss.bucket)
		}

		return slide
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return slide
	}

	newSlide := Slide{}

	err = json.Unmarshal(data, &newSlide)

	if err != nil {
		return slide
	}

	return newSlide
}

// Resolve will resolve a slide if that slide was previously uploaded
func (ss *SlideStorage) Upload(slide Slide) Slide {
	data, err := json.Marshal(slide)

	if err != nil {
		log.Printf("Error uploading slide: %v key: %s bucket: %s", err, buildKey(ss.prefix, slide), ss.bucket)
		return slide
	}

	input := &s3.PutObjectInput{
		Key:         aws.String(buildKey(ss.prefix, slide)),
		Bucket:      aws.String(ss.bucket),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/json"),
		ACL:         aws.String("public-read"),
	}

	if _, err := ss.sss.PutObject(input); err != nil {
		log.Printf("%v key: %s bucket: %s", err, buildKey(ss.prefix, slide), ss.bucket)
		return slide
	}

	return slide
}

type SlideResolverTransform struct {
	slideChanIn  chan Slide
	slideChanOut chan Slide
	ss           *SlideStorage
}

func NewSlideResolverTransform(slideChanIn chan Slide, slideChanOut chan Slide, ss *SlideStorage) *SlideResolverTransform {
	return &SlideResolverTransform{
		slideChanIn:  slideChanIn,
		slideChanOut: slideChanOut,
		ss:           ss,
	}
}

func (sc *SlideResolverTransform) Run() {
	for slide := range sc.slideChanIn {
		sc.slideChanOut <- sc.ss.Resolve(slide)
	}
	close(sc.slideChanOut)
}

type SlideUploader struct {
	slideChanIn  chan Slide
	slideChanOut chan Slide
	ss           *SlideStorage
}

func NewSlideUploader(slideChanIn chan Slide, slideChanOut chan Slide, ss *SlideStorage) *SlideUploader {
	return &SlideUploader{
		slideChanIn:  slideChanIn,
		slideChanOut: slideChanOut,
		ss:           ss,
	}
}

func (sc *SlideUploader) Run() {
	for slide := range sc.slideChanIn {
		sc.slideChanOut <- sc.ss.Upload(slide)
	}
	close(sc.slideChanOut)
}
