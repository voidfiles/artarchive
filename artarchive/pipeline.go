package artarchive

import "sync"

type SlideProcessor interface {
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
	go func() {
		defer wg.Done()
		p.Producer.Run()
	}()

	for _, step := range p.Steps {
		wg.Add(1)
		go func(step SlideProcessor) {
			defer wg.Done()
			step.Run()
		}(step)
	}

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
