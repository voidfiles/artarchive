package indexer

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/rs/zerolog"
	"github.com/voidfiles/artarchive/slides"
)

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

type SlideshowRenderer struct {
	logger    zerolog.Logger
	ss        IndexerS3Interface
	binding   slides.Binding
	public    string
	bucket    string
	publisher *Publisher
}

func NewSlideshowRenderer(logger zerolog.Logger, ss IndexerS3Interface, public, bucket string) *SlideshowRenderer {
	publisher := NewPublisher(logger, ss, bucket, 1)
	publisher.Start()

	return &SlideshowRenderer{
		logger:    logger.With().Str("renderer", "slideshow").Logger(),
		ss:        ss,
		public:    public,
		bucket:    bucket,
		publisher: publisher,
	}
}

func (sr *SlideshowRenderer) Configure(binding slides.Binding) {
	sr.binding = binding
}

func (sr *SlideshowRenderer) slideToHTML(slide slides.Slide) string {
	content := ""
	for _, artistInfo := range slide.ArtistsInfo {
		content += fmt.Sprintf("<h1>%s</h1>\n", artistInfo.Name)
	}

	content += fmt.Sprintf("<a href='%s/slide-editor/index.html?key=%s'>edit</a>", sr.public, slide.GUIDHash)
	return fmt.Sprintf(`<section data-background-image="%s/images/%s" data-background-size="contain">%s</section>`, sr.public, slide.ArchivedImage.Filename, content)
}

func (sr *SlideshowRenderer) Run() {
	sr.logger.Info().Msg("Running slideshow renderer")
	templatedSlides := make([]string, 0)
	for slide := range sr.binding.In {
		if slide.ArchivedImage != nil {
			templatedSlides = append(templatedSlides, sr.slideToHTML(slide))
		}
		sr.binding.Out <- slide
	}

	sr.logger.Info().Msg("Choosing randome slides")
	randomeSlides := make([]string, 100)
	for i := 0; i < 100; i++ {
		index := rand.Intn(len(templatedSlides))
		randomeSlides[i] = templatedSlides[index]
	}

	sr.logger.Info().Msg("Templating slideshow")
	templatedSlideString := strings.Join(randomeSlides, "\n")

	fullHTML := fmt.Sprintf(mainDoc, templatedSlideString)

	sr.logger.Info().Msg("Sending template to slideshow")
	sr.publisher.Add("art/slideshow.html", fullHTML)
	sr.logger.Info().Msg("Waiting for publisher")
	sr.publisher.Wait()
	sr.logger.Info().Msg("Publisher is done")
	close(sr.binding.Out)
}
