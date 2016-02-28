package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth/gothic"
)

func corsHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
		next(w, r)
	}
}

func authCallbackHandler(privateKey []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		token := jwt.New(jwt.GetSigningMethod("RS256"))
		str, _ := token.SignedString(privateKey)

		var data = struct {
			Token string `json:"token"`
		}{
			str,
		}
		b, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("error: %s\n", err)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func authCheckHandler(publicKey []byte, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseFromRequest(r, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})
		if err == nil && token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
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
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	usr, err := users.GetByName(r.URL.Query().Get(":name"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if usr != nil {
		handleUser(w, usr)
	} else {
		handleNotFound(w, "user not found")
	}
}

func postUserHandler(w http.ResponseWriter, r *http.Request) {
	usr, err := bodyToUser(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	handlePost(w, usr)
}

func putUserHandler(w http.ResponseWriter, r *http.Request) {
	usr, err := bodyToUser(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	handlePut(w, usr)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	usr, err := bodyToUser(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	handleDelete(w, usr)
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
