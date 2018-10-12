package debug

import (
	"github.com/rs/zerolog"
	"github.com/voidfiles/artarchive/slides"
)

// DebugLimitSlideConsumer keeps track of how many slides its seen
type DebugLimitSlideConsumer struct {
	binding    slides.Binding
	slidesSeen int
	limit      int
	logger     zerolog.Logger
}

// NewDebugLimitSlideConsumer creates a new DebugSlideConsumer
func NewDebugLimitSlideConsumer(logger zerolog.Logger, limit int) *DebugLimitSlideConsumer {
	return &DebugLimitSlideConsumer{
		logger: logger,
		limit:  limit,
	}
}

// Configure adds a channel bindings to the transformer
func (dsc *DebugLimitSlideConsumer) Configure(binding slides.Binding) {
	dsc.binding = binding
}

// Run runs the NewDebugLimitSlideConsumer
func (dsc *DebugLimitSlideConsumer) Run() {
	for slide := range dsc.binding.In {
		dsc.slidesSeen++
		dsc.binding.Out <- slide
		if dsc.slidesSeen >= dsc.limit {
			dsc.logger.Info().
				Int("limit", dsc.limit).
				Msg("Hit limit")
			break
		}
	}
	close(dsc.binding.Out)
}
