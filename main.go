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

	log.Fatal(http.ListenAndServe(":8080", r))
}

func queryDecode(response http.ResponseWriter, request *http.Request) {
	var source sources.Source

	vars := mux.Vars(request)
	page, _ := strconv.ParseInt(vars["page"], 10, 0)

	switch vars["source"] {
	case "org":
		source = sources.Org{}
		break
	default:
		log.Printf("Unhandled URL: %s\n", request.URL)
		http.NotFound(response, request)
		return
	}

	feed, err := source.CreateFeed(vars["id"], int(page))
	if err != nil {
		log.Printf("Error creating feed: %s\n", err)
		http.NotFound(response, request)
		return
	}

	response.Header().Set("Content-Type", "application/atom+xml")
	encoder := xml.NewEncoder(response)
	encoder.Indent("", "	")
	err = encoder.Encode(feed)
	if err != nil {
		log.Printf("Error encoding feed: %s\n", err)
		http.NotFound(response, request)
		return
	}
}
