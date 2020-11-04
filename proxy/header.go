package proxy

import "net/http"

// CopyHeaders copy headers from source to destination.
// Nothing would be returned.
func CopyHeaders(dst, src http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

//ClearHeader clear headers.
func ClearHeader(headers http.Header) {
	for key := range headers {
		headers.Del(key)
	}
}
