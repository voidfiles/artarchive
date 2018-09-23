package server

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/artarchive/slides"
	"github.com/voidfiles/artarchive/storage"
)

type MockSlideStore struct {
	data  slides.Slide
	slide slides.Slide
	err   error
}

func (m *MockSlideStore) FindByKey(key string) (slides.Slide, error) {
	return m.data, m.err
}

func (m *MockSlideStore) UpdateByKey(key string, data interface{}) error {
	return nil
}

func (m *MockSlideStore) Resolve(s slides.Slide) slides.Slide {
	return m.slide
}

func (m *MockSlideStore) Upload(s slides.Slide) slides.Slide {
	return m.slide
}

type MockRequestContext struct {
	param   string
	aborted bool
	code    int
	data    interface{}
}

func (m *MockRequestContext) Param(key string) string {
	return m.param
}

func (m *MockRequestContext) AbortWithStatusJSON(code int, data interface{}) {
	m.aborted = true
	m.code = code
	m.data = data
}

func (m *MockRequestContext) JSON(code int, data interface{}) {
	m.aborted = false
	m.code = code
	m.data = data
}

func (m *MockRequestContext) ShouldBindJSON(obj interface{}) error {
	return nil
}

func TestMustNewServerHandlers(t *testing.T) {
	mockStore := &MockSlideStore{}
	subject := MustNewServerHandlers(zerolog.New(os.Stdout), mockStore, mockStore)
	assert.IsType(t, &ServerHandlers{}, subject)
}

func TestGetSlide(t *testing.T) {
	testTable := []struct {
		key         string
		expectAbort bool
		code        int
		data        interface{}
		mockData    slides.Slide
		mockSlide   slides.Slide
		mockError   error
	}{
		{"a key", false, 200, slides.Slide{}, slides.Slide{}, slides.Slide{}, nil},
		{"", true, 400, map[string]string{"error": "missing-param"}, slides.Slide{}, slides.Slide{}, nil},
		{"123", true, 404, map[string]string{"error": "object-missing"}, slides.Slide{}, slides.Slide{}, storage.ErrMissingSlide},
	}
	for _, test := range testTable {
		mockStore := &MockSlideStore{
			data:  test.mockData,
			slide: test.mockSlide,
			err:   test.mockError,
		}
		handlers := MustNewServerHandlers(zerolog.New(os.Stderr), mockStore, mockStore)
		c := &MockRequestContext{
			param: test.key,
		}

		handlers.GetSlide(c)
		assert.Equal(t, test.key, c.param)
		assert.Equal(t, test.expectAbort, c.aborted)
		assert.Equal(t, test.code, c.code)
		assert.Equal(t, test.data, c.data)
	}
}

func TestUpdateSlide(t *testing.T) {
	testTable := []struct {
		key         string
		expectAbort bool
		code        int
		data        interface{}
		mockData    slides.Slide
		mockSlide   slides.Slide
		mockError   error
	}{
		{"a key", false, 200, slides.Slide{}, slides.Slide{}, slides.Slide{}, nil},
	}
	for _, test := range testTable {
		mockStore := &MockSlideStore{
			data:  test.mockData,
			slide: test.mockSlide,
			err:   test.mockError,
		}
		handlers := MustNewServerHandlers(zerolog.New(os.Stderr), mockStore, mockStore)
		c := &MockRequestContext{
			param: test.key,
		}

		handlers.UpdateSlide(c)
		assert.Equal(t, test.key, c.param)
		assert.Equal(t, test.expectAbort, c.aborted)
		assert.Equal(t, test.code, c.code)
	}
}
