package doers

import (
	"log"

	"github.com/jasonlvhit/gocron"
)

func RunCron() {
	gocron.Every(5).Minutes().Do(func() {
		log.Printf("Running feed runner")
		FeedRunner()
	})
	<-gocron.Start()
}
