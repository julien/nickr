package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"regexp"
)

var (
	fbURL     = os.Getenv("NICKR_FIREBASE_URL")
	port      = flag.String("port", os.Getenv("PORT"), "http port")
	users     = NewUsers(fbURL)
	usersPath = regexp.MustCompile(`(users/?)(\w+)?`)
)

type response struct {
	Message string                 `json:"message"`
	Errors  map[string]interface{} `json:"errros,omitempty"`
}

func init() {
	if *port == "" {
		*port = "8080"
	}
	if fbURL == "" {
		log.Fatal("make sure the NICKR_FIREBASE_URL environment variable is set")
	}
}

func main() {
	flag.Parse()

	http.Handle("/", AddCORS(handleRequest(), "*", "X-Requested-With", "GET,POST,PUT,DELETE"))
	http.ListenAndServe(":"+*port, nil)
}
