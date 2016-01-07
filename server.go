package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type topic struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

const dataFile = "./temp.json" // TODO: replace with mongodb
var dataMutex = new(sync.Mutex)

func handleTopics(w http.ResponseWriter, r *http.Request) {
	dataMutex.Lock()
	defer dataMutex.Unlock()

	// Read the topics from the file.
	topicsData, err := ioutil.ReadFile(dataFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read the data file (%s): %s", dataFile, err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	// stream the contents of the file to the response
	io.Copy(w, bytes.NewReader(topicsData))

}

func main() {
	http.HandleFunc("/api/topics", handleTopics)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
