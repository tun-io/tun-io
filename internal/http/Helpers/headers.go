package Helpers

import "net/http"

func HeadersToMap(headers http.Header) map[string]string {
	headersMap := make(map[string]string)
	for name, values := range headers {
		if len(values) > 0 {
			headersMap[name] = values[0] // Use the first value for simplicity
		}
	}
	return headersMap
}

func HeadersFromMap(headers map[string]string) http.Header {
	httpHeaders := http.Header{}
	for name, value := range headers {
		httpHeaders.Add(name, value)
	}
	return httpHeaders
}
