package main

import (
	"log"
)

// DebugSlideConsumer keeps track of how many slides its seen
type DebugSlideConsumer struct {
	slideChan  chan Slide
	slidesSeen int
}

// NewDebugSlideConsumer creates a new DebugSlideConsumer
func NewDebugSlideConsumer(slideChan chan Slide) *DebugSlideConsumer {
	return &DebugSlideConsumer{
		slideChan: slideChan,
	}
}

// Run runs the DebugSlideConsumer
func (dsc *DebugSlideConsumer) Run() {
	for _ = range dsc.slideChan {
		dsc.slidesSeen++
		if dsc.slidesSeen%100 == 0 {
			log.Printf("Slides seen: %v", dsc.slidesSeen)
		}
	}
}
