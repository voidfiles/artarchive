package server

import (
	"encoding/json"
	"net/http"

	"github.com/voidfiles/artarchive/slides"
	"github.com/voidfiles/artarchive/storage"
)

type RequestContext interface {
	Param(string) string
	AbortWithStatusJSON(int, interface{})
	JSON(int, interface{})
}

type RemoteSlideStore interface {
	Resolve(slides.Slide) slides.Slide
}

type LocalSlideStore interface {
	FindByKey(string) ([]byte, error)
}

type ServerHandlers struct {
	slideS3Storage RemoteSlideStore
	slideDbSTorage LocalSlideStore
}

func MustNewServerHandlers(slideS3Storage RemoteSlideStore, slideDbSTorage LocalSlideStore) *ServerHandlers {
	return &ServerHandlers{
		slideDbSTorage: slideDbSTorage,
		slideS3Storage: slideS3Storage,
	}
}

func (sh *ServerHandlers) GetSlide(c RequestContext) {
	key := c.Param("key")
	if key == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{
			"error": "missing-param",
		})
		return
	}
	data, err := sh.slideDbSTorage.FindByKey(key)
	if err != nil {
		if err == storage.ErrMissingSlide {
			c.AbortWithStatusJSON(http.StatusNotFound, map[string]string{
				"error": "object-missing",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
			"error": "server-error",
		})
		return

	}
	slide := slides.Slide{}
	err = json.Unmarshal(data, &slide)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
			"error": "server-error",
		})
		return

	}
	c.JSON(http.StatusOK, slide)
}
