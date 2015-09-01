package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

type Entity struct {
	Name      string   `json:"name"`
	Nicknames []string `json:"nicknames"`
}

var (
	port = flag.String("port", os.Getenv("PORT"), "http port")
)

func init() {
	if *port == "" {
		*port = "8080"
	}
}

func main() {
	flag.Parse()

	fmt.Printf("listening on port: %s\n", *port)

	http.Handle("/", indexHandler())
	http.ListenAndServe(":"+*port, nil)
}

func indexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Nickr")
	})
}
