package doers

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/voidfiles/artarchive/config"
	"github.com/voidfiles/artarchive/debug"
	"github.com/voidfiles/artarchive/logging"
	"github.com/voidfiles/artarchive/pipeline"
	"github.com/voidfiles/artarchive/scanner"
	"github.com/voidfiles/artarchive/storage"
)

func RunSlideSync() error {
	appConfig := config.NewAppConfig()
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	sss := s3.New(sess)
	logger := logging.NewLogger(false, os.Stdout)
	db, err := sqlx.Connect(appConfig.Database.Type, appConfig.Database.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	slideScanner := scanner.NewSlideScanner(sss, 0, appConfig.Bucket, appConfig.Version)

	itemStore := storage.MustNewItemStorage(db)
	slideStore := storage.NewDBStorageTransform(logger, itemStore)
	dropStore := storage.NewDBStorageDropTransform(logger, itemStore)
	debugConsumer := debug.NewDebugSlideConsumer(logger)
	// Dump things
	pipeline := pipeline.NewPipeline(slideScanner, debugConsumer, dropStore, slideStore)
	pipeline.Run()
	return nil
}
