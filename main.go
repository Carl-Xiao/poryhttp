package main

import (
	"proxyHTTP/handler"

	"github.com/google/logger"
)

func main() {
	proxy := handler.NewServer()
	logger.Fatal(proxy.ListenAndServe())
}
