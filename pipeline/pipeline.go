package pipeline

import (
	"sync"

	"github.com/voidfiles/artarchive/slides"
)

type SlideProcessor interface {
	Configure(slides.Binding)
	Run()
}

type Pipeline struct {
	Producer SlideProcessor
	Steps    []SlideProcessor
	Consumer SlideProcessor
}

func (p *Pipeline) Run() {
	var wg = sync.WaitGroup{}
	wg.Add(1)
	start := slides.Binding{
		In:  make(chan slides.Slide, 0),
		Out: make(chan slides.Slide, 0),
	}
	go func() {
		defer wg.Done()
		p.Producer.Configure(start)
		p.Producer.Run()
	}()

	stepBinding := slides.Binding{
		In:  start.Out,
		Out: make(chan slides.Slide, 0),
	}

	for _, step := range p.Steps {
		wg.Add(1)

		go func(step SlideProcessor, stepBinding slides.Binding) {
			defer wg.Done()
			step.Configure(stepBinding)
			step.Run()
		}(step, stepBinding)

		stepBinding = slides.Binding{
			In:  stepBinding.Out,
			Out: make(chan slides.Slide, 0),
		}
	}

	end := slides.Binding{
		In:  stepBinding.Out,
		Out: make(chan slides.Slide, 0),
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.Producer.Configure(end)
		p.Consumer.Run()
	}()

	wg.Wait()

}

func NewPipeline(producer SlideProcessor, consumer SlideProcessor, steps ...SlideProcessor) *Pipeline {
	return &Pipeline{
		Producer: producer,
		Steps:    steps,
		Consumer: consumer,
	}
}
