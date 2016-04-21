// Package main launches uTeach, a web platform for sharing educational material and resources.
package main

import (
	"flag"
	"fmt"
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

	httpAddress := fmt.Sprintf("%s:%d", app.Config.Address, app.Config.HttpPort)
	httpsAddress := fmt.Sprintf("%s:%d", app.Config.Address, app.Config.HttpsPort)

	log.Println("Serving http at", httpAddress, "https at", httpsAddress)

	redirectToHttps := func(w http.ResponseWriter, r *http.Request) {
		url := fmt.Sprintf("https://%s%s", httpsAddress, r.RequestURI)
		http.Redirect(w, r, url, http.StatusFound)
	}

	go func() {
		log.Fatal(http.ListenAndServe(httpAddress, http.HandlerFunc(redirectToHttps)))
	}()

	log.Fatal(http.ListenAndServeTLS(httpsAddress, app.Config.HttpsCertPath, app.Config.HttpsKeyPath, nil))
}
