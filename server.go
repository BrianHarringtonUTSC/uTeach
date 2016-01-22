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
	subjects, err := GetSubjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderTemplate(w, r, "subjects.html", subjects)
}

func handleTopics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]

	topics, err := GetTopics(subjectName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderTemplate(w, r, "topics.html", topics)
}

func handleThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectName := vars["subjectName"]
	topicName := vars["topicName"]

	threads, err := GetThreads(subjectName, topicName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RenderTemplate(w, r, "threads.html", threads)
}

func handleThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threadID, err := strconv.Atoi(vars["threadID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	thread, err := GetThread(threadID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	RenderTemplate(w, r, "thread.html", thread)
}

func main() {
	LoadTemplates()

	InitDB()
	defer DB.Close()

	// allows user to be encoded so that it can be stored in a session
	gob.Register(&User{})

	authMiddleWare := alice.New(isAuth)

	router := mux.NewRouter()
	router.HandleFunc("/", handleSubjects)
	router.HandleFunc("/subjects", handleSubjects)
	router.HandleFunc("/topics/{subjectName}", handleTopics)
	router.HandleFunc("/threads/{subjectName}/{topicName}", handleThreads)
	router.HandleFunc("/thread/{subjectName}/{topicName}/{threadID}", handleThread)

	router.HandleFunc("/login/{username}", handleLogin)
	router.HandleFunc("/logout", handleLogout)
	router.Handle("/user", authMiddleWare.ThenFunc(handleUser))

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
