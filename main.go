// Package main launches uTeach, a web platform for sharing educational material and resources.
package main

// Import necessary packages
import (
	"flag"
	"log"
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/config"
	"github.com/umairidris/uTeach/handlers"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to config file.")
	flag.Parse()

	if configPath == "" {
		log.Fatal("--config arg is missing.")
	}

	conf, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}
	app := application.New(*conf)
	router := handlers.Router(app)
	http.Handle("/", router)

	log.Println("Starting server at", app.Config.HTTPAddress)
	log.Fatal(http.ListenAndServe(app.Config.HTTPAddress, nil))
}
