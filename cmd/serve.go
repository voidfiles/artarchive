package cmd

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload" // Keeps track of go stats
	"github.com/spf13/cobra"
	"github.com/voidfiles/artarchive/config"
	"github.com/voidfiles/artarchive/logging"
)

func serve() {
	appConfig := config.NewAppConfig()
	logger := logging.NewLogger(false, os.Stdout)
	ginLogger := logging.MustNewGinLogger(logger, "gin")
	router := gin.New()

	router.Use(ginLogger.Logger())

	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": appConfig.AuthPassword,
	}))

	authorized.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})

	router.Run(":" + appConfig.Port)
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs an http server",
	Long:  `Runs an http server`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("cmd: serve")
		serve()
	},
}
