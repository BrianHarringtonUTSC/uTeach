// Package main runs a web app, uTeach, which is a platform for sharing educational material and resources.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/handlers"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to JSON config file.")
	flag.Parse()

	if configPath == "" {
		fmt.Println("ERROR: --config arg is missing.",
			"See https://raw.githubusercontent.com/umairidris/uTeach/master/sample/config.json for an example.")
		os.Exit(1)
	}

	app := application.New(configPath)
	router := handlers.Router(app)
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(app.Config.HTTPAddress, nil))
}
