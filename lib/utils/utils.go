package utils

import (
	"io"
	"encoding/json"
	"bytes"
)

func CreateBodyFromMap(data map[string]interface{}) (io.Reader, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return bytes.NewReader([]byte{}), err
	}

	return bytes.NewReader(b), nil
}
