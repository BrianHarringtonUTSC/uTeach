package main

import (
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"strconv"
)

func handleSubjects(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, "subjects.html", GetSubjects())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]
	err := RenderTemplate(w, "topics.html", GetTopics(subjectName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]
	topicName := vars["topicName"]
	err := RenderTemplate(w, "threads.html", GetThreads(subjectName, topicName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]
	topicName := vars["topicName"]
	threadID, err := strconv.Atoi(vars["threadID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = RenderTemplate(w, "thread.html", GetThread(subjectName, topicName, threadID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// allows user to be encoded so that it can be stored in a session
	gob.Register(&User{})

	LoadTemplates()

	authMiddleWare := alice.New(isAuth)

	router := mux.NewRouter()
	router.HandleFunc("/", handleSubjects)
	router.HandleFunc("/subjects", handleSubjects)
	router.HandleFunc("/topics/{subjectName}", handleTopics)
	router.HandleFunc("/threads/{subjectName}/{topicName}", handleThreads)
	router.HandleFunc("/thread/{subjectName}/{topicName}/{threadID}", handleThread)

	router.HandleFunc("/login/{utorid}", handleLogin)
	router.HandleFunc("/logout", handleLogout)
	router.Handle("/user", authMiddleWare.ThenFunc(handleUser))

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
