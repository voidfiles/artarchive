package indexer

import (
	"github.com/rs/zerolog"
	"github.com/voidfiles/artarchive/slides"
)

type SlidePager struct {
	logger  zerolog.Logger
	binding slides.Binding
	page    int
	part    []slides.Slide
}

func NewSlidePager(logger zerolog.Logger, binding slides.Binding) *SlidePager {
	return &SlidePager{
		logger:  logger,
		binding: binding,
		page:    0,
		part:    make([]slides.Slide, 100),
	}
}

func (sp *SlidePager) Next() bool {
	sp.logger.Info().Int("page", sp.page).Msg("starting a page")
	sp.part = make([]slides.Slide, 100)
	counter := 0
	for slide := range sp.binding.In {
		sp.part[counter] = slide
		sp.binding.Out <- slide
		counter++
		if counter == len(sp.part) {
			sp.page++
			sp.logger.Info().Int("page", sp.page).Int("counter", counter).Msg("Reached page")
			return true
		}
	}
	sp.logger.Info().Int("page", sp.page).Int("counter", counter).Msg("No more pages")
	close(sp.binding.Out)
	return false
}
