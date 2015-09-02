package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/julien/nickr/utils"
)

const file string = "data.json"

var (
	port       = flag.String("port", os.Getenv("PORT"), "http port")
	characters map[string][]string
	fileLock   sync.RWMutex
)

func init() {
	if *port == "" {
		*port = "8080"
	}
}

func main() {
	flag.Parse()

	var err error
	characters, err = loadData()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Printf("Listening on port: %s\n", *port)
	http.Handle("/", utils.Cors(characterHandler()))
	http.ListenAndServe(":"+*port, nil)
}

func indexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Nickr")
	})
}

func loadData() (map[string][]string, error) {
	c := make(map[string][]string)

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return c, err
}

// TODO: GET       character
// TODO: POST      character [nick1, nick2, nicnN]
// TODO: PUT/PATCH character [nick1, nick2, nickN]
// TODO: DELETE    character
func characterHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// if r.Method == "GET" {
		// 	http.Error(w, "Shit ...", http.StatusMethodNotAllowed)
		// }

		var err error

		name := r.URL.Path[1:]
		if ok := characters[name]; ok != nil {
			fmt.Printf("Found it: %v\n", ok)

			switch r.Method {
			case "GET":
				b, err := json.Marshal(ok)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(b)

			case "PATCH":
				s, err := utils.ReadFromBody(r.Body)
				if err != nil {
					log.Fatal(err)
				}
				characters[name] = s
				b, err := json.Marshal(characters[name])
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				err = utils.Flush(file, characters, fileLock)
				if err != nil {
					log.Fatal("write error: %s\n", err)
				}
				characters, err = loadData()
				if err != nil {
					log.Fatal(err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(b)

			case "PUT":
				s, err := utils.ReadFromBody(r.Body)
				if err != nil {
					log.Fatal(err)
				}
				characters[name] = utils.AddNewNicknames(ok, s)
				b, err := json.Marshal(characters[name])
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				err = utils.Flush(file, characters, fileLock)
				if err != nil {
					log.Fatal("write error: %s\n", err)
				}
				characters, err = loadData()
				if err != nil {
					log.Fatal(err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(b)

			case "DELETE":
				delete(characters, name)
				err = utils.Flush(file, characters, fileLock)
				if err != nil {
					log.Fatalf("write error: %s\n", err)
				}
				characters, err = loadData()
				if err != nil {
					log.Fatalf("Error: %v\n", err)
				}
				w.WriteHeader(http.StatusNoContent)
			}
		} else {
			if r.Method != "POST" {
				http.Error(w, "Sorry we didn't find what you were looking for", http.StatusNotFound)
			} else {
				s, err := utils.ReadFromBody(r.Body)
				if err != nil {
					log.Fatal(err)
				}
				characters[name] = s
				b, err := json.Marshal(characters[name])
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				err = utils.Flush(file, characters, fileLock)
				if err != nil {
					log.Fatal("write error: %s\n", err)
				}
				characters, err = loadData()
				if err != nil {
					log.Fatal(err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(b)
			}
		}

	})
}
