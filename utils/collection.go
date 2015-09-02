package utils

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

// Collection is a map of items
type Collection map[string][]string

var lock sync.RWMutex

// ContainsString checks if a given array/slice of string has
// the specified value
func ContainsString(src []string, val string) bool {
	for _, i := range src {
		if i == val {
			return true
		}
	}
	return false
}

// FromJSON populates the Collection from a JSON file
func (c Collection) FromJSON(path string) (Collection, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

// Add adds a new key in the Collection with the given values
func (c Collection) Add(key string, values []string) Collection {
	items := c.Get(key)

	if c != nil {
		for _, i := range values {
			if !ContainsString(items, i) {
				items = append(items, i)
			}
		}
		c[key] = items
	} else {
		c[key] = values
	}

	return c
}

// Get retreives a key according to the specified name
func (c Collection) Get(key string) []string {
	return c[key]
}

// Set sets the value of the key according to the specified name and given values
func (c Collection) Set(key string, values []string) Collection {
	if ok := c.Get(key); ok != nil {
		c[key] = values
	}
	return c
}

// Delete removes a key from the Collection given the specified name
func (c Collection) Delete(key string) {
	delete(c, key)
}

// Flush writes the Collection as JSON to a file given the specified name
func (c Collection) Flush(path string) error {

	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	lock.Lock()
	out, err := json.MarshalIndent(c, "", "  ")
	lock.Unlock()

	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	_, err = w.WriteString(string(out))
	w.Flush()
	return nil
}
