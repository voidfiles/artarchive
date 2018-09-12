package storage

import (
	"github.com/rs/zerolog"
	"github.com/voidfiles/artarchive/slides"
)

type DBStorageTransform struct {
	binding slides.Binding
	i       *ItemStorage
	logger  zerolog.Logger
}

func NewDBStorageTransform(logger zerolog.Logger, i *ItemStorage) *DBStorageTransform {
	return &DBStorageTransform{
		i:      i,
		logger: logger,
	}
}

func (d *DBStorageTransform) Configure(binding slides.Binding) {
	d.binding = binding
}

func (d *DBStorageTransform) Run() {
	for slide := range d.binding.In {
		d.i.Store(slide.GUIDHash, slide)
		d.binding.Out <- slide
	}
	close(d.binding.Out)
}

type DBStorageDropTransform struct {
	binding slides.Binding
	i       *ItemStorage
	logger  zerolog.Logger
}

func NewDBStorageDropTransform(logger zerolog.Logger, i *ItemStorage) *DBStorageDropTransform {
	return &DBStorageDropTransform{
		i:      i,
		logger: logger,
	}
}

func (d *DBStorageDropTransform) Configure(binding slides.Binding) {
	d.binding = binding
}

func (d *DBStorageDropTransform) Run() {
	for slide := range d.binding.In {
		if _, err := d.i.FindByKey(slide.GUIDHash); err != nil {
			d.binding.Out <- slide
		}

	}
	close(d.binding.Out)
}
