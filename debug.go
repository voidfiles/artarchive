package main

import (
	"fmt"
	"log"
)

type DebugSlideConsumer struct {
	slideChan chan Slide
}

func NewDebugSlideConsumer(slideChan chan Slide) *DebugSlideConsumer {
	return &DebugSlideConsumer{
		slideChan: slideChan,
	}
}

func (dsc *DebugSlideConsumer) Run() {
	for slide := range dsc.slideChan {
		log.Printf("Site: %s", slide.Site.Title)
		log.Printf("Page: %s", slide.Page.URL)
		log.Printf("Data Blob: %s", fmt.Sprintf("%s/slides/%s/data.json", "v2", slide.GUIDHash))
		log.Printf("Image URL: %s", slide.SourceImageURL)
		log.Printf("Image Info: %v", slide.ArchivedImage)
		println("")
	}
}
