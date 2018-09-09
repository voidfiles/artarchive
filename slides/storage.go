package slides

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

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

func BuildKey(prefix string, slide Slide) string {
	if prefix != "" {
		return fmt.Sprintf("%s/slides/%s/data.json", prefix, slide.GUIDHash)
	}

	return fmt.Sprintf("slides/%s/data.json", slide.GUIDHash)
}

// Resolve will resolve a slide if that slide was previously uploaded
func (ss *SlideStorage) Resolve(slide Slide) Slide {
	resp, err := ss.sss.GetObject(&s3.GetObjectInput{
		Key:    aws.String(BuildKey(ss.prefix, slide)),
		Bucket: aws.String(ss.bucket),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != s3.ErrCodeNoSuchKey {
				log.Printf("Error resolving slide: %v key: %s bucket: %s", err, BuildKey(ss.prefix, slide), ss.bucket)
			}
		} else {
			log.Printf("Error resolving slide: %v key: %s bucket: %s", err, BuildKey(ss.prefix, slide), ss.bucket)
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
		log.Printf("Error uploading slide: %v key: %s bucket: %s", err, BuildKey(ss.prefix, slide), ss.bucket)
		return slide
	}

	input := &s3.PutObjectInput{
		Key:         aws.String(BuildKey(ss.prefix, slide)),
		Bucket:      aws.String(ss.bucket),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/json"),
		ACL:         aws.String("public-read"),
	}

	if _, err := ss.sss.PutObject(input); err != nil {
		log.Printf("%v key: %s bucket: %s", err, BuildKey(ss.prefix, slide), ss.bucket)
		return slide
	}

	return slide
}
