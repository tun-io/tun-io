package Helpers

import (
	"io"
	"net/http"
)

func GetRequestBody(r *http.Request) []byte {
	if r.Body == nil {
		return nil
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	return body
}
