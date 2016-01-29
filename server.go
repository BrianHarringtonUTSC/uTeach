package main

import (
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
)

func main() {
	LoadTemplates()

	err := InitDB()
	if err != nil {
		panic(err)
	}

	// allows user to be encoded so that it can be stored in a session
	gob.Register(&User{})

	authMiddleWare := alice.New(AuthorizedHandler)

	router := mux.NewRouter()
	router.HandleFunc("/", handleGetSubjects)
	router.HandleFunc("/topics/{subjectName}", handleGetTopics)
	router.HandleFunc("/threads/{subjectName}/{topicName}", handleGetThreads)
	router.HandleFunc("/thread/{subjectName}/{topicName}/{threadID}", handleGetThread)

	router.HandleFunc("/login/{username}", handleLogin)
	router.HandleFunc("/logout", handleLogout)

	router.Handle("/upvote/{threadID}", authMiddleWare.ThenFunc(handleAddUpvote)).Methods("POST")
	router.Handle("/upvote/{threadID}", authMiddleWare.ThenFunc(handleRemoveUpvote)).Methods("DELETE")

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
