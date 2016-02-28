package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

var (
	fbURL                 = os.Getenv("NICKR_FIREBASE_URL")
	port                  = flag.String("port", os.Getenv("PORT"), "http port")
	users                 = NewUsers(fbURL)
	usersPath             = regexp.MustCompile(`(users/?)(\w+)?`)
	privateKey, publicKey []byte
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

	var err error
	privateKey, err = ioutil.ReadFile("keys/key.rsa")
	if err != nil {
		log.Fatal(err)
	}
	publicKey, err = ioutil.ReadFile("keys/key.rsa.pub")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	goth.UseProviders(github.New(
		os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"),
		os.Getenv("NICKR_OAUTH_CALLBACK_URL")),
	)

	p := pat.New()
	p.Get("/auth/{provider}/callback", authCallbackHandler(privateKey))
	p.Get("/login/{provider}", gothic.BeginAuthHandler)
	p.Get("/users", authCheckHandler(publicKey, corsHandler(getUsersHandler)))
	p.Get("/user/{name}", authCheckHandler(publicKey, corsHandler(getUserHandler)))
	p.Post("/user", authCheckHandler(publicKey, corsHandler(postUserHandler)))
	p.Put("/user", authCheckHandler(publicKey, corsHandler(putUserHandler)))
	p.Delete("/user", authCheckHandler(publicKey, corsHandler(deleteUserHandler)))

	http.ListenAndServe(":"+*port, p)
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

func encodeJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func decodeJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, &v)
}
