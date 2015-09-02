package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julien/nickr/utils"
)

const file string = "data.json"

var (
	port       = flag.String("port", os.Getenv("PORT"), "http port")
	collection = utils.Collection{}
)

func init() {
	if *port == "" {
		*port = "8080"
	}
}

func main() {
	flag.Parse()

	loadCollection()

	fmt.Printf("Listening on port: %s\n", *port)
	http.Handle("/", utils.AddCORSHeaders(characterHandler()))
	http.ListenAndServe(":"+*port, nil)
}

func loadCollection() error {
	_, err := collection.FromJSON(file)
	return err
}

func itemsToJSON(items []string) ([]byte, error) {
	return json.Marshal(items)
}

func itemsFromBody(body io.Reader) ([]string, error) {
	var err error
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var s []string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func addToCollection(key string, body io.Reader, w http.ResponseWriter) error {
	items, err := itemsFromBody(body)
	if err != nil {
		return err
	}

	collection.Add(key, items)

	if err := collection.Flush(file); err != nil {
		return err
	}

	loadCollection()
	return nil
}

func characterHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[1:]

		if name == "" {
			w.Header().Set("Content-type", "text/html")
			w.Write([]byte("<h1>NickR</h1>"))
			return
		}

		if it := collection.Get(name); it != nil {
			switch r.Method {

			case "DELETE":
				collection.Delete(name)
				if err := collection.Flush(file); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				loadCollection()
				w.WriteHeader(http.StatusNoContent)

			case "GET":
				res, err := itemsToJSON(it)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(res)

			case "PATCH":
				items, err := itemsFromBody(r.Body)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				collection.Set(name, items)
				if err := collection.Flush(file); err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				loadCollection()
				w.WriteHeader(http.StatusOK)

			case "PUT":
				if err := addToCollection(name, r.Body, w); err != nil {
					log.Fatal(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				w.WriteHeader(http.StatusOK)
			}
		} else if r.Method == "POST" {
			if err := addToCollection(name, r.Body, w); err != nil {
				log.Fatal(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	})
}
