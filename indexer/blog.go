package indexer

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"html/template"
	"log"

	"github.com/rs/zerolog"
	"github.com/voidfiles/artarchive/slides"
)

type BlogRenderer struct {
	logger    zerolog.Logger
	ss        IndexerS3Interface
	binding   slides.Binding
	public    string
	bucket    string
	publisher *Publisher
}

func NewBlogRenderer(logger zerolog.Logger, ss IndexerS3Interface, concurrency int, public, bucket string) *BlogRenderer {
	publisher := NewPublisher(logger, ss, bucket, concurrency)
	publisher.Start()
	return &BlogRenderer{
		logger:    logger.With().Str("renderer", "blog").Logger(),
		ss:        ss,
		public:    public,
		bucket:    bucket,
		publisher: publisher,
	}
}

var blogDoc = `<html>
	<head>
		<!-- Required meta tags -->
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

		<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css" integrity="sha384-MCw98/SFnGE8fJT3GXwEOngsV7Zt27NXFoaoApmYm81iuXoPkFOJwJ8ERdknLPMO" crossorigin="anonymous">

	</head>
	<body>
		<div class="container">
			<div class="row">
		    <div class="col">
					<div class="card">
						{{ if (.ArchivedImage) and (.ArchivedImage.Filename) }}
					  <img class="card-img-top" src="https://s3.us-west-2.amazonaws.com/art.rumproarious.com/images/{{ .ArchivedImage.Filename }}" alt="Card image cap">
						{{ end }}
						<div class="card-body">
							{{ if (.WorkInfo) and (.WorkInfo.Name) }}
					    <h5 class="card-title">{{ .WorkInfo.Name }}<h5>
							{{ end }}
							{{ range $i, $artist := .ArtistsInfo }}
   							<p class="card-text">$artist.Name</p>
							{{ end }}
					  </div>
					</div>
		    </div>
		  </div>
		</div>
		{{ . }}
	</body>
</html>`

var indexDoc = `<html>
	<head>
		<!-- Required meta tags -->
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

		<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css" integrity="sha384-MCw98/SFnGE8fJT3GXwEOngsV7Zt27NXFoaoApmYm81iuXoPkFOJwJ8ERdknLPMO" crossorigin="anonymous">

	</head>
	<body>
		<div class="container">
			{{ range $i, $slide := . }}
				{{ if ($slide.ArchivedImage) and ($slide.ArchivedImage.Filename) }}
				<div class="row">
			    <div class="col">
						<div class="card">
						  <img class="card-img-top" src="https://s3.us-west-2.amazonaws.com/art.rumproarious.com/images/{{ $slide.ArchivedImage.Filename }}" alt="Card image cap">
						  </div>
						</div>
			    </div>
			  </div>
				{{ end }}
			{{ end }}
		</div>
	</body>
</html>`

func hashIt(in []byte) string {
	h := sha1.New()
	h.Write(in)
	bs := h.Sum(nil)
	return fmt.Sprintf("%x\n", bs)
}

func RenderBlogSlide(slide slides.Slide) []byte {
	t, err := template.New("blog").Parse(blogDoc) // Create a template.
	if err != nil {
		log.Printf("Ran into an error %v", err)
		return []byte{}
	}
	out := bytes.NewBufferString("")
	err = t.Execute(out, slide)
	if err != nil {
		log.Printf("Ran into an error %v", err)
		return []byte{}
	}

	return out.Bytes()
}

func RenderBlogSlideIndex(slides []slides.Slide) []byte {
	t, err := template.New("index").Parse(indexDoc) // Create a template.
	if err != nil {
		log.Printf("Ran into an error %v", err)
		return []byte{}
	}
	out := bytes.NewBufferString("")
	err = t.Execute(out, slides)
	if err != nil {
		log.Printf("Ran into an error %v", err)
		return []byte{}
	}

	return out.Bytes()
}

func (br *BlogRenderer) Configure(binding slides.Binding) {
	br.binding = binding
}

func (br *BlogRenderer) SlideToHtml(slide slides.Slide) (string, string) {
	data := RenderBlogSlide(slide)
	return string(data), hashIt(data)
}

func (br *BlogRenderer) SlideToPagePart(slide slides.Slide) (string, string) {
	data := RenderBlogSlide(slide)
	return string(data), hashIt(data)
}

func (br *BlogRenderer) Run() {
	pager := NewSlidePager(br.logger, br.binding)
	var page []slides.Slide

	for pager.Next() {
		page = pager.part
		for _, slide := range page {
			content, hash := br.SlideToHtml(slide)
			if slide.RenderHash.Blog != hash {
				br.logger.Info().
					Str("oldHash", slide.RenderHash.Blog).
					Str("newHash", hash).
					Str("GUIDHash", slide.GUIDHash).
					Msg("Will republish")
				slide.RenderHash.Blog = hash
				br.publisher.Add(
					fmt.Sprintf("blog/items/%s.html", slide.GUIDHash),
					content,
				)
			}
			br.binding.Out <- slide
		}
		br.publisher.Add(
			fmt.Sprintf("blog/page-%d.html", pager.page),
			string(RenderBlogSlideIndex(page)),
		)
	}
	br.logger.Info().Msg("End of blog publisher waiting")
	br.publisher.Wait()
	br.logger.Info().Msg("Waiting")
	close(br.binding.Out)
}
