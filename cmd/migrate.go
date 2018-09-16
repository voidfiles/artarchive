package cmd

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // For sqlx
	"github.com/spf13/cobra"
	"github.com/voidfiles/artarchive/config"
	"github.com/voidfiles/artarchive/storage"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Run database migrations`,
	Run: func(cmd *cobra.Command, args []string) {
		appConfig := config.NewAppConfig()

		db, err := sqlx.Connect(appConfig.Database.Type, appConfig.Database.DatabaseURL)
		if err != nil {
			log.Fatal(err)
		}
		store := storage.MustNewItemStorage(db)
		store.MigrateDB()
	},
}
