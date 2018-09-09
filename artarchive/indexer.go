package artarchive

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Interface holds methods we use to interact with s3
type IndexerS3Interface interface {
	PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error)
}

type Indexer struct {
	ss       IndexerS3Interface
	slidesIn chan Slide
	public   string
	bucket   string
	prefix   string
}

func NewIndexer(ss IndexerS3Interface, slidesIn chan Slide, public, bucket, prefix string) *Indexer {
	return &Indexer{
		ss:       ss,
		slidesIn: slidesIn,
		public:   public,
		bucket:   bucket,
		prefix:   prefix,
	}
}

var mainDoc = `<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/reveal.js/3.6.0/css/reveal.css">
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/reveal.js/3.6.0/css/theme/white.css">
	</head>
	<body>
		<div class="reveal">
			<div class="slides">
        %s
			</div>
		</div>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/reveal.js/3.6.0/js/reveal.js"></script>
		<script>
			Reveal.initialize({shuffle: true, controls: false, loop: true, autoSlide: 60000});
		</script>
	</body>
</html>`

func (i *Indexer) slideToHTML(slide Slide) string {
	content := ""
	for _, artistInfo := range slide.ArtistsInfo {
		content += fmt.Sprintf("<h1>%s</h1>\n", artistInfo.Name)
	}

	link := buildKey("v2", slide)
	content += fmt.Sprintf("<a href='http://art.rumproarious.com/slide-editor/?data=%s'>edit</a>", link)
	return fmt.Sprintf(`<section data-background-image="%s/images/%s" data-background-size="contain">%s</section>`, i.public, slide.ArchivedImage.Filename, content)
}

func (i *Indexer) uploadDoc(html string) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(i.bucket),
		Key:         aws.String("art/slideshow.html"),
		Body:        bytes.NewReader([]byte(html)),
		ContentType: aws.String("text/html"),
		ACL:         aws.String("public-read"),
	}
	_, err := i.ss.PutObject(input)
	if err != nil {
		log.Fatal(err)
	}
}

func (i *Indexer) Run() {
	templatedSlides := make([]string, 0)
	for slide := range i.slidesIn {
		if slide.ArchivedImage != nil {
			templatedSlides = append(templatedSlides, i.slideToHTML(slide))
		}
	}
	randomeSlides := make([]string, 100)
	for i := 0; i < 100; i++ {
		index := rand.Intn(len(templatedSlides))
		randomeSlides[i] = templatedSlides[index]
	}
	templatedSlideString := strings.Join(randomeSlides, "\n")

	fullHTML := fmt.Sprintf(mainDoc, templatedSlideString)
	i.uploadDoc(fullHTML)
}
