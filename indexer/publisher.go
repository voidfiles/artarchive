package indexer

import (
	"bytes"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog"
)

type PublishItem struct {
	key     string
	content string
}

type Publisher struct {
	ss          IndexerS3Interface
	bucket      string
	in          chan PublishItem
	logger      zerolog.Logger
	wg          sync.WaitGroup
	concurrency int
}

func NewPublisher(logger zerolog.Logger, ss IndexerS3Interface, bucket string, concurrency int) *Publisher {
	in := make(chan PublishItem)
	return &Publisher{
		ss:          ss,
		bucket:      bucket,
		in:          in,
		logger:      logger,
		wg:          sync.WaitGroup{},
		concurrency: concurrency,
	}
}

func (p *Publisher) Start() {
	p.logger.Info().Int("concurrency", p.concurrency).Msg("Starting consume")
	workers := make([]string, p.concurrency)
	for i, _ := range workers {
		p.logger.Info().Msgf("Starting worker %v", i)
		p.wg.Add(1)
		go func(i int) {
			defer p.wg.Done()
			p.logger.Info().Msgf("worker %v consuming", i)
			p.Consume()
			p.logger.Info().Msgf("worker %v stopped consuming", i)
		}(i)
	}
}

func (p *Publisher) Add(key, content string) {
	p.logger.Info().Str("key", key).Msg("Adding content to queue")
	p.in <- PublishItem{key: key, content: content}
}

func (p *Publisher) Wait() {
	p.logger.Info().Msg("Closing and waiting for channel")
	close(p.in)
	p.wg.Wait()
	p.logger.Info().Msg("Wait group stopped")
}

func (p *Publisher) Consume() {
	for pi := range p.in {
		p.logger.Info().Str("key", pi.key).Msg("Publishing")
		input := &s3.PutObjectInput{
			Bucket:      aws.String(p.bucket),
			Key:         aws.String(pi.key),
			Body:        bytes.NewReader([]byte(pi.content)),
			ContentType: aws.String("text/html"),
			ACL:         aws.String("public-read"),
		}
		_, err := p.ss.PutObject(input)
		if err != nil {
			log.Fatal(err)
		}

	}
}
