package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func ReadFromBody(body io.Reader) ([]string, error) {
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
