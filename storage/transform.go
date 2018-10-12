package storage

import (
	"github.com/rs/zerolog"
	"github.com/voidfiles/artarchive/slides"
)

type DBStorageProducer struct {
	binding slides.Binding
	i       *ItemStorage
	logger  zerolog.Logger
}

func NewDBStorageProducer(logger zerolog.Logger, i *ItemStorage) *DBStorageProducer {
	return &DBStorageProducer{
		i:      i,
		logger: logger,
	}
}

func (d *DBStorageProducer) Configure(binding slides.Binding) {
	d.binding = binding
}

func (d *DBStorageProducer) Run() {
	moreSlides := true
	next := int64(0)
	slides := make([]slides.Slide, 0)
	for {
		if !moreSlides {
			break
		}
		slides, next, _ = d.i.List(next)
		for _, slide := range slides {
			d.binding.Out <- slide
		}
		if next == int64(-1) {
			moreSlides = false
		}
	}
	close(d.binding.Out)
}

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

type DBUpdateTransform struct {
	binding slides.Binding
	i       *ItemStorage
	logger  zerolog.Logger
}

func NewDBUpdateTransform(logger zerolog.Logger, i *ItemStorage) *DBUpdateTransform {
	return &DBUpdateTransform{
		i:      i,
		logger: logger,
	}
}

func (d *DBUpdateTransform) Configure(binding slides.Binding) {
	d.binding = binding
}

func (d *DBUpdateTransform) Run() {
	for slide := range d.binding.In {
		d.i.UpdateByKey(slide.GUIDHash, slide)
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
