package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const fbURL = "https://nickr.firebaseio.com/users/"

var (
	port  = flag.String("port", os.Getenv("PORT"), "http port")
	users = NewUsers(fbURL)
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

	fmt.Printf("Listening on port: %s\n", *port)
	http.Handle("/", AddCORS(collectionHandler(), "*", "X-Requested-With", "GET,POST,PUT,DELETE"))
	http.ListenAndServe(":"+*port, nil)
}

func encodeJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func decodeJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, &v)
}

func collectionHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			name := r.URL.Path[1:]
			if name == "" {
				// load all data in the "background"
				if _, err := users.All(); err != nil {
					fmt.Printf("Error loading data: %v\n", err)
				}

				w.Header().Set("Content-type", "text/html")
				w.Write([]byte("<h1>NickR</h1>"))
				return
			}

			usr, err := users.GetByName(name)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if usr != nil {
				handleGet(w, usr)
			} else {
				handleNotFound(w, "user not found")
			}

		} else {
			usr, err := bodyToUser(r.Body)
			if err != nil {
				fmt.Printf("Body error: %v\n", err)
			}

			switch r.Method {
			case "POST":
				handlePost(w, usr)
			case "PUT":
				handlePut(w, usr)
			case "DELETE":
				handleDelete(w, usr)
			default:
				handleNotFound(w, "user not found")
			}
		}
	})
}

func bodyToByte(body io.Reader) ([]byte, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func bodyToUser(body io.Reader) (*User, error) {
	b, err := bodyToByte(body)
	if err != nil {
		return nil, err
	}

	usr := &User{}
	if err := decodeJSON(b, usr); err != nil {
		return nil, err
	}

	return usr, nil
}

func handleGet(w http.ResponseWriter, usr *User) {
	res, err := encodeJSON(usr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func handlePost(w http.ResponseWriter, usr *User) {
	if err := users.Add(usr); err != nil {
		msg := &response{Message: fmt.Sprintf("%s", err)}
		res, err := encodeJSON(msg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return
	}

	res, err := encodeJSON(usr)
	if err != nil {
		fmt.Printf("Encode error: %v\n", err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func handlePut(w http.ResponseWriter, usr *User) {
	if id := users.GetUserID(usr.Name); id != "" {

		u, err := users.Update(id, usr)
		if err != nil {
			fmt.Printf("Set error: %v\n", err)
		}

		if u != nil {
			res, err := encodeJSON(u)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(res)
		}

	} else {
		handleNotFound(w, "user not found")
	}
}

func handleDelete(w http.ResponseWriter, usr *User) {
	if err := users.Delete(usr.Name); err != nil {
		handleNotFound(w, fmt.Sprintf("%s", err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleNotFound(w http.ResponseWriter, s string) {

	msg := &response{Message: s}
	res, err := encodeJSON(msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write(res)
}
