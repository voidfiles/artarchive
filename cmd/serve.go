package cmd

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload" // Keeps track of go stats
	"github.com/spf13/cobra"
	"github.com/voidfiles/artarchive/config"
)

func serve() {
	appConfig := config.NewAppConfig()

	router := gin.New()
	router.Use(gin.Logger())

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
