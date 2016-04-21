// Package main launches uTeach, a web platform for sharing educational material and resources.
package main

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
	defer app.DB.Close()

	router := handlers.Router(app)
	http.Handle("/", router)

	log.Println("Serving at", app.Config.HttpAddress)
	log.Fatal(http.ListenAndServe(app.Config.HttpAddress, nil))
}
