package main

import (
	"github.com/robfig/cron"
	"log"
	"time"
)

func main() {
	log.Println("Starting...")

	c := cron.New()
	c.AddFunc("@every 1s", func() {
		log.Println("Running1...", time.Now().Format("2006-01-02 15:04:05"))
	})

	c.AddFunc("@every 2s", func() {
		log.Println("Running2...", time.Now().Format("2006-01-02 15:04:05"))
	})

	c.Start()

	select {}

}
