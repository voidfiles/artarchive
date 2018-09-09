package pipeline

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/artarchive/slides"
)

type MockProducer struct {
	name    string
	binding slides.Binding
	ran     bool
	counter int
}

func (m *MockProducer) Configure(binding slides.Binding) {
	log.Printf("%s.Configure(%v)", m.name, binding)
	m.binding = binding
}

func (m *MockProducer) Run() {
	m.ran = true
	log.Printf("%s.Run() - producing", m.name)
	slides := []slides.Slide{
		slides.Slide{},
		slides.Slide{},
		slides.Slide{},
	}
	for _, slide := range slides {
		m.counter++
		log.Printf("%s.Run sending slide: %v", m.name, slide)
		m.binding.Out <- slide
	}
	close(m.binding.Out)
}

type MockConsumer struct {
	name    string
	binding slides.Binding
	ran     bool
	counter int
}

func (m *MockConsumer) Configure(binding slides.Binding) {
	log.Printf("%s.Configure(%v)", m.name, binding)
	m.binding = binding
}
func (m *MockConsumer) Run() {
	m.ran = true
	log.Printf("%s.Run() - consuming", m.name)
	for slide := range m.binding.In {
		m.counter++
		log.Printf("%s.Run recieving slide: %v", m.name, slide)
	}
}

func TestNewPipeline(t *testing.T) {
	producer := &MockProducer{name: "producer", ran: false}
	consumer := &MockConsumer{name: "consumer", ran: false}
	pipeline := NewPipeline(producer, consumer)
	assert.IsType(t, &Pipeline{}, pipeline)
	pipeline.Run()
	assert.Equal(t, consumer.ran, true)
	assert.Equal(t, producer.ran, true)
	assert.Equal(t, producer.counter, 3)
	assert.Equal(t, consumer.counter, 3)
}
