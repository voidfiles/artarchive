package artarchive

import (
	"github.com/rs/zerolog"
)

// DebugSlideConsumer keeps track of how many slides its seen
type DebugSlideConsumer struct {
	slideChan  chan Slide
	slidesSeen int
	logger     zerolog.Logger
}

// NewDebugSlideConsumer creates a new DebugSlideConsumer
func NewDebugSlideConsumer(logger zerolog.Logger, slideChan chan Slide) *DebugSlideConsumer {
	return &DebugSlideConsumer{
		slideChan: slideChan,
		logger:    logger,
	}
}

// Run runs the DebugSlideConsumer
func (dsc *DebugSlideConsumer) Run() {
	for _ = range dsc.slideChan {
		dsc.slidesSeen++
		if dsc.slidesSeen%100 == 0 {
			dsc.logger.Info().
				Int("seen", dsc.slidesSeen).
				Msg("Slide checkpoint")
		}
	}
}
