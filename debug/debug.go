package debug

import (
	"github.com/rs/zerolog"
	"github.com/voidfiles/artarchive/slides"
)

// DebugSlideConsumer keeps track of how many slides its seen
type DebugSlideConsumer struct {
	binding    slides.Binding
	slidesSeen int
	logger     zerolog.Logger
}

// NewDebugSlideConsumer creates a new DebugSlideConsumer
func NewDebugSlideConsumer(logger zerolog.Logger) *DebugSlideConsumer {
	return &DebugSlideConsumer{
		logger: logger,
	}
}

// Run runs the DebugSlideConsumer
func (dsc *DebugSlideConsumer) Configure(binding slides.Binding) {
	dsc.binding = binding
}

// Run runs the DebugSlideConsumer
func (dsc *DebugSlideConsumer) Run() {
	for slide := range dsc.binding.In {
		dsc.logger.Info().
			Msgf("debugger: %v", slide.GUIDHash)
		dsc.slidesSeen++
		if dsc.slidesSeen%100 == 0 {
			dsc.logger.Info().
				Int("seen", dsc.slidesSeen).
				Msg("Slide checkpoint")
		}
	}
}
