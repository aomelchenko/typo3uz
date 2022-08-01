package main

import (
	"log"
	"net/http"
	"typo3uz/database"
	"typo3uz/service"
)

func main() {
	// load config
	config, err := database.FromEnv()
	if err != nil {
		log.Fatalf(`failed to read config from env: %s`, err.Error())
	}
	if err = config.Validate(); err != nil {
		log.Fatalf("config validation failed: %s", err.Error())
	}

	db := database.InitDb(config)

	mux := http.NewServeMux()
	mux.Handle("/", &service.TGHandler{Session: db})
	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatalln(err)
	}
}
