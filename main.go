// Package uTeach implements a web app which is a platform for sharing educational material and resources.
package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"strconv"

	"github.com/umairidris/uTeach/app"
	"github.com/umairidris/uTeach/middleware"
	"github.com/umairidris/uTeach/routes"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config_path", "",
		"Path to JSON config file. See github.com/UmairIdris/uTeach/blob/master/sample_config.json for an example.")
	flag.Parse()

	if configPath == "" {
		fmt.Println("config_path arg not provided. ")
		return
	}

	app := app.New(configPath)

	middleware := middleware.Middleware{app}
	routeHandler := routes.RouteHandler{app}

	authMiddleWare := alice.New(middleware.AuthorizedHandler)

	router := mux.NewRouter()
	router.HandleFunc("/", routeHandler.GetSubjects)
	router.HandleFunc("/topics/{subjectName}", routeHandler.GetTopics)
	router.HandleFunc("/threads/{subjectName}/{topicName}", routeHandler.GetThreads)
	router.HandleFunc("/thread/{subjectName}/{topicName}/{threadID}", routeHandler.GetThread)

	router.HandleFunc("/user/{username}", routeHandler.GetUser)

	router.HandleFunc("/login/{username}", routeHandler.Login)
	router.HandleFunc("/logout", routeHandler.Logout)

	router.Handle("/upvote/{threadID}", authMiddleWare.ThenFunc(routeHandler.AddUpvote)).Methods("POST")
	router.Handle("/upvote/{threadID}", authMiddleWare.ThenFunc(routeHandler.RemoveUpvote)).Methods("DELETE")

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	address := app.Config.Host + ":" + strconv.Itoa(int(app.Config.Port))
	log.Fatal(http.ListenAndServe(address, nil))
}
