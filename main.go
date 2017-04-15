package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"strconv"

	"github.com/Cj-Malone/feed-org/sources"
	"github.com/gorilla/mux"
	//"github.com/nenadl/atom"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/{source}/{id}/{page:[0-9]+}.atom", queryDecode)
	r.HandleFunc("/{source}/{page:[0-9]+}.atom", queryDecode)
	r.HandleFunc("/{source}.atom", queryDecode)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func queryDecode(responseWriter http.ResponseWriter, request *http.Request) {
	var source sources.Source

	vars := mux.Vars(request)
	page := 1
	if vars["page"] != "" {
		pageI, _ := strconv.ParseInt(vars["page"], 10, 0)
		page = int(pageI)
	}

	switch vars["source"] {
	case "org":
		source = sources.Org{}
		break
	case "imdb":
		source = sources.Imdb{}
		break
	case "iplayer":
		source = sources.Iplayer{}
		break
	default:
		log.Printf("Unhandled URL: %s\n", request.URL)
		http.Error(responseWriter, "Feed not found", http.StatusNotFound)
		return
	}

	feed, err := source.CreateFeed(vars["id"], int(page))
	if err != nil {
		log.Printf("Error creating feed: %s\n", err)
		http.Error(responseWriter, "Feed creation error", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/atom+xml")
	encoder := xml.NewEncoder(responseWriter)
	encoder.Indent("", "	")
	err = encoder.Encode(feed)
	if err != nil {
		log.Printf("Error encoding feed: %s\n", err)
		http.Error(responseWriter, "Feed encoding error", http.StatusInternalServerError)
		return
	}
}
