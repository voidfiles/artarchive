package server

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/voidfiles/artarchive/config"
	"github.com/voidfiles/artarchive/logging"
	"github.com/voidfiles/artarchive/slides"
	"github.com/voidfiles/artarchive/storage"
)

type ContextFunc func(c RequestContext)

func bind(f ContextFunc) func(*gin.Context) {
	return func(c *gin.Context) {
		f(c)
	}
}

func serve() {
	appConfig := config.NewAppConfig()
	logger := logging.NewLogger(false, os.Stdout)

	db, err := sqlx.Connect(appConfig.Database.Type, appConfig.Database.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	sss := s3.New(sess)

	ginLogger := logging.MustNewGinLogger(logger, "gin")
	router := gin.New()

	router.Use(ginLogger.Logger())

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": appConfig.AuthPassword,
	}))

	slidesDBStorage := storage.MustNewItemStorage(db)
	slideS3Storage := slides.NewSlideStorage(sss, appConfig.Bucket, appConfig.Version)
	handlers := MustNewServerHandlers(slideS3Storage, slidesDBStorage)

	authorized.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
	authorized.GET("/slides/*key", bind(handlers.GetSlide))

	router.Run(":" + appConfig.Port)
}
