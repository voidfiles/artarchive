package cmd

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/spf13/cobra"
)

func serve() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})

	router.Run(":" + port)
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
