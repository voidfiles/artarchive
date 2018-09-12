package cmd

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // For sqlx
	"github.com/spf13/cobra"
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
		log.Print("cmd: database migration")
		db, err := sqlx.Connect("postgres", "user=stitchfix_owner password=stitchfix_owner dbname=artarchive sslmode=disable")
		if err != nil {
			log.Fatalln(err)
		}
		store := storage.MustNewItemStorage(db)
		store.MigrateDB()
	},
}
