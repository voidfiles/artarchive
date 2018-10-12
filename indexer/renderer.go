package indexer

import "github.com/voidfiles/artarchive/slides"

type Renderer interface {
	Run()
	Configure(slides.Binding)
}
