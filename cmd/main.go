package main

import (
	"ChiliOverFlow/pkg/api"
	"ChiliOverFlow/pkg/db"
	"log"
)

func main() {
	mod, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	go api.SetupAPI(&mod)
	select {}
}
