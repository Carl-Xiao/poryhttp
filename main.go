package main

import (
	"log"
	"os"
	"proxyHTTP/handler"

	"github.com/google/logger"
)

const logPath = "./example.log"

func main() {
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()
	defer logger.Init("LoggerExample", true, true, lf).Close()
	logger.SetFlags(log.Lshortfile)
	proxy := handler.NewServer()
	logger.Fatal(proxy.ListenAndServe())

	// backend := ":8888"
	// proxy := ":8002"

	// backendURL := &url.URL{Scheme: "ws://", Host: backend}
	// p := wsutil.NewSingleHostReverseProxy(backendURL)

	// go http.ListenAndServe(proxy, p)
	// time.Sleep(100 * time.Second)

}
