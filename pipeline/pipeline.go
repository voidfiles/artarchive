package pipeline

import (
	"log"
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

	bindings := make([]*slides.Binding, len(p.Steps)+2)
	var in chan slides.Slide
	var out chan slides.Slide
	for i := range bindings {
		if i == 0 {
			in = nil
		} else {
			in = bindings[i-1].Out
		}
		out = make(chan slides.Slide, 0)
		log.Printf("In: %v Out: %v", in, out)
		bindings[i] = &slides.Binding{
			In:  in,
			Out: out,
		}
	}

	p.Producer.Configure(*bindings[0])

	log.Printf("Running the producer")
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.Producer.Run()
	}()

	for i, step := range p.Steps {

		step.Configure(*bindings[i+1])

		wg.Add(1)
		go func(step SlideProcessor) {
			defer wg.Done()
			step.Run()
		}(step)
	}

	log.Printf("Running the consumer")
	p.Consumer.Configure(*bindings[len(p.Steps)+1])
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
