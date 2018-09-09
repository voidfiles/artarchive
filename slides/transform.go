package slides

import "log"

type SlideResolverTransform struct {
	binding Binding
	ss      *SlideStorage
}

func NewSlideResolverTransform(ss *SlideStorage) *SlideResolverTransform {
	return &SlideResolverTransform{
		ss: ss,
	}
}

func (sc *SlideResolverTransform) Configure(binding Binding) {
	sc.binding = binding
}

func (sc *SlideResolverTransform) Run() {
	for slide := range sc.binding.In {
		log.Printf("Transform: %v", slide.GUIDHash)
		sc.binding.Out <- sc.ss.Resolve(slide)
	}
	close(sc.binding.Out)
}

type SlideUploader struct {
	binding Binding
	ss      *SlideStorage
}

func NewSlideUploader(ss *SlideStorage) *SlideUploader {
	return &SlideUploader{
		ss: ss,
	}
}

func (sc *SlideUploader) Configure(binding Binding) {
	sc.binding = binding
}

func (sc *SlideUploader) Run() {
	for slide := range sc.binding.In {
		sc.binding.Out <- sc.ss.Upload(slide)
	}
	close(sc.binding.Out)
}
