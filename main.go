package main

import (
	"flag"
	"net/http"
	"os"
	"regexp"
)

const fbURL = "https://nickr.firebaseio.com/users/"

var (
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
}

func main() {
	flag.Parse()

	http.Handle("/", AddCORS(handleRequest(), "*", "X-Requested-With", "GET,POST,PUT,DELETE"))
	http.ListenAndServe(":"+*port, nil)
}
