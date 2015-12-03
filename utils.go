package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// AddCORS adds Cross Origin Resource Sharing headers to responses.
func AddCORS(next http.Handler, origin, headers, methods string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if len(origin) == 0 {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)

		if len(headers) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", headers)
		}

		if len(methods) > 0 {
			w.Header().Set("Access-Control-Allow-Methods", methods)
		}
		next.ServeHTTP(w, r)
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

func encodeJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func decodeJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, &v)
}
