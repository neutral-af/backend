package utils

import (
	"bytes"
	"encoding/json"
	"io"
)

func CreateBodyFromMap(data map[string]interface{}) (io.Reader, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return bytes.NewReader([]byte{}), err
	}

	return bytes.NewReader(b), nil
}
