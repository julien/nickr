package utils

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
)

func Flush(path string, data interface{}, lock sync.RWMutex) error {

	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}

	lock.Lock()
	out, err := json.MarshalIndent(data, "", " ")
	lock.Unlock()

	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	_, err = w.WriteString(string(out))
	w.Flush()
	return nil
}
