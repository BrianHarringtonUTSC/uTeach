// Package main runs a web app, uTeach, which is a platform for sharing educational material and resources.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/umairidris/uTeach/application"
	"github.com/umairidris/uTeach/handlers"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "",
		"Path to JSON config file. See github.com/umairidris/uTeach/blob/master/sample/config.json for an example.")
	flag.Parse()

	if configPath == "" {
		fmt.Println("config arg not provided.")
		return
	}

	app := application.New(configPath)
	http.Handle("/", handlers.Router(app))
	log.Fatal(http.ListenAndServe(app.Config.HTTPAddress, nil))
}
