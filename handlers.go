package main

import (
	"fmt"
	"net/http"
)

func handleRequest() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path[1:]
		matches := usersPath.FindAllStringSubmatch(path, -1)

		// if we don't match anything render static content
		if len(matches) == 0 {
			fmt.Fprintf(w, "nickr")
			return
		}

		if len(matches) > 0 {
			submatches := matches[0]
			if len(submatches) == 3 {
				if r.Method == "GET" {
					handleGet(w, submatches)
				} else {
					if r.Method == "OPTIONS" {
						return
					}

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
			}
		}
	})
}

func handleGet(w http.ResponseWriter, submatches []string) {

	if submatches[2] != "" {
		// we have a name
		fmt.Printf("submatch: %s\n", submatches[2])

		usr, err := users.GetByName(submatches[2])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if usr != nil {
			handleUser(w, usr)
		} else {
			handleNotFound(w, "user not found")
		}
	} else {
		// render all users
		all, err := users.All()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := encodeJSON(all)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.Write(b)
		return
	}
}

func handleUser(w http.ResponseWriter, usr *User) {
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
			fmt.Printf("Update error: %v\n", err)
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
