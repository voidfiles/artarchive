package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ScanerS3Interface interface {
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
	ListObjectsV2Pages(input *s3.ListObjectsV2Input, fn func(*s3.ListObjectsV2Output, bool) bool) error
}

type SlideScanner struct {
	ss       ScanerS3Interface
	outBound chan Slide
	bucket   string
	prefix   string
	random   int
}

func NewSlideScanner(ss ScanerS3Interface, outBound chan Slide, random int, bucket, prefix string) *SlideScanner {
	return &SlideScanner{
		ss:       ss,
		outBound: outBound,
		bucket:   bucket,
		prefix:   prefix,
		random:   random,
	}
}

func (ss *SlideScanner) resolveAndSendSlide(key *string) {
	resp, err := ss.ss.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(ss.bucket),
		Key:    key,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	newSlide := Slide{}

	err = json.Unmarshal(data, &newSlide)

	ss.outBound <- newSlide
}

func (ss *SlideScanner) Run() {
	params := &s3.ListObjectsV2Input{
		Bucket:  aws.String(ss.bucket),
		Prefix:  aws.String(ss.prefix),
		MaxKeys: aws.Int64(1000),
	}

	result := func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		innerSlideChan := make(chan *string, 8)
		var wg sync.WaitGroup
		for _ = range []int{1, 2, 3, 4, 5, 6, 7, 8} {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for key := range innerSlideChan {
					ss.resolveAndSendSlide(key)
				}
			}()
		}
		for _, object := range page.Contents {
			if ss.random > 0 {
				if rand.Intn(ss.random) < 10 {
					innerSlideChan <- object.Key
				}
			} else {
				innerSlideChan <- object.Key
			}
		}
		close(innerSlideChan)
		wg.Wait()
		if lastPage {
			close(ss.outBound)
		}
		return true
	}
	err := ss.ss.ListObjectsV2Pages(params, result)
	if err != nil {
		log.Fatal(err)
	}
}
