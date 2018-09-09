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

	start := slides.Binding{
		Out: make(chan slides.Slide, 0),
	}
	p.Producer.Configure(start)
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.Producer.Run()
	}()

	var stepBinding = slides.Binding{}

	for i, step := range p.Steps {
		if i == 0 {
			stepBinding = slides.Binding{
				In:  start.Out,
				Out: make(chan slides.Slide, 0),
			}
		} else {
			stepBinding = slides.Binding{
				In:  stepBinding.Out,
				Out: make(chan slides.Slide, 0),
			}
		}
		wg.Add(1)
		step.Configure(stepBinding)
		go func(step SlideProcessor) {
			defer wg.Done()
			step.Run()
		}(step)
	}

	end := slides.Binding{
		In: stepBinding.Out,
	}
	p.Producer.Configure(end)
	wg.Add(1)
	go func() {
		defer wg.Done()
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
