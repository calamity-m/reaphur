package serr

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func EncodeJSON[T any](val T) (string, error) {

	w := bytes.Buffer{}
	if err := json.NewEncoder(&w).Encode(val); err != nil {
		return "", fmt.Errorf("failed encoding json: %w", err)
	}

	return w.String(), nil
}

func EncodeJSONArr[T any](arr []T) (string, error) {
	w := bytes.Buffer{}
	if err := json.NewEncoder(&w).Encode(arr); err != nil {
		return "", fmt.Errorf("failed encoding json: %w", err)
	}

	return w.String(), nil
}

func DecodeJSON[T any](barray []byte) (T, error) {
	var val T

	if err := json.NewDecoder(bytes.NewReader(barray)).Decode(&val); err != nil {
		return val, fmt.Errorf("failed decoding json: %w", err)
	}

	return val, nil
}

func DecodeJSONS[T any](jstr string) (T, error) {
	return DecodeJSON[T]([]byte(jstr))
}
