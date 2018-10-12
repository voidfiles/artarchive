package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/voidfiles/artarchive/indexer"
	"github.com/voidfiles/artarchive/slides"
	"github.com/voidfiles/artarchive/storage"
	"gopkg.in/go-playground/validator.v9"
)

type RequestContext interface {
	Param(string) string
	AbortWithStatusJSON(int, interface{})
	JSON(int, interface{})
	Data(code int, contentType string, data []byte)
	ShouldBindJSON(obj interface{}) error
	DefaultQuery(key, defaultValue string) string
}

type RemoteSlideStore interface {
	Resolve(slides.Slide) slides.Slide
	Upload(slides.Slide) slides.Slide
}

type LocalSlideStore interface {
	FindByKey(string) (slides.Slide, error)
	UpdateByKey(string, interface{}) error
	List(after int64) ([]slides.Slide, int64, error)
	FindSites(query string, after int64) ([]slides.Site, int64, error)
}

type ServerHandlers struct {
	slideS3Storage RemoteSlideStore
	slideDbSTorage LocalSlideStore
	logger         zerolog.Logger
}

func MustNewServerHandlers(logger zerolog.Logger, slideS3Storage RemoteSlideStore, slideDbSTorage LocalSlideStore) *ServerHandlers {
	return &ServerHandlers{
		slideDbSTorage: slideDbSTorage,
		slideS3Storage: slideS3Storage,
		logger:         logger,
	}
}

func (sh *ServerHandlers) GetSlide(c RequestContext) {
	key := c.Param("key")
	sh.logger.Info().Str("key", key).Msgf("Loooking for slide by key: %s", key)
	if key == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{
			"error": "missing-param",
		})
		return
	}
	slide, err := sh.slideDbSTorage.FindByKey(key)
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

	c.JSON(http.StatusOK, slide)
}

type JSONError struct {
	Field      string `json:"field"`
	Validation string `json:"validation"`
}

func errorsToInfo(errors validator.ValidationErrors) []JSONError {
	returnErrors := make([]JSONError, len(errors))
	for i, err := range errors {
		returnErrors[i] = JSONError{
			Field:      err.Namespace(),
			Validation: err.Tag(),
		}
	}
	return returnErrors
}

func (sh *ServerHandlers) UpdateSlide(c RequestContext) {
	slide := slides.Slide{}
	if err := c.ShouldBindJSON(&slide); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{
				"error":   "server-error",
				"details": "invalid validation",
			})
			return
		}
		errors := err.(validator.ValidationErrors)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "validation-errors",
			"details": errorsToInfo(errors),
		})
		return
	}

	key := c.Param("key")
	sh.logger.Info().Str("key", key).Msgf("Updating slide by key: %s", key)
	if key == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{
			"error": "missing-param",
		})
		return
	}
	slide.Edited = time.Now().UTC()
	if err := sh.slideDbSTorage.UpdateByKey(key, slide); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
			"error": "server-error",
		})
		return

	}
	sh.slideS3Storage.Upload(slide)

	c.JSON(http.StatusOK, slide)
}

func (sh *ServerHandlers) ListSlides(c RequestContext) {
	after, err := strconv.ParseInt(c.DefaultQuery("after", "0"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{
			"error": "invalid-query",
		})
		return
	}
	slides, next, err := sh.slideDbSTorage.List(after)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
			"error": "server-error",
		})
		return

	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"next":   next,
		"slides": slides,
	})
}

func (sh *ServerHandlers) ListSites(c RequestContext) {
	after, err := strconv.ParseInt(c.DefaultQuery("after", "0"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{
			"error": "invalid-query",
		})
		return
	}
	query := c.DefaultQuery("query", "")
	sites, next, err := sh.slideDbSTorage.FindSites(query, after)
	if err != nil {
		sh.logger.Error().AnErr("error", err).Msg("An error occured")
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]string{
			"error": "server-error",
		})
		return

	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"next":  next,
		"sites": sites,
	})
}

func (sh *ServerHandlers) RenderSlide(c RequestContext) {
	key := c.Param("key")
	sh.logger.Info().Str("key", key).Msgf("Render slide: %s", key)

	slide, err := sh.slideDbSTorage.FindByKey(key)
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

	c.Data(http.StatusOK, "text/html", indexer.RenderBlogSlide(slide))
}
