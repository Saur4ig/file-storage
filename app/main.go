package main

import (
	"log"

	"github.com/saur4ig/file-storage/internal/config"
)

func main() {
	log.Println("Starting app")
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	log.Println(conf)
}
