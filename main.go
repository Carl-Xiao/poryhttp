package main

import (
	"fmt"
	"io"
	"net/http"
	"proxyHTTP/proxy"
	"time"

	"github.com/google/logger"
)

// Server is a server of proxy.
type Server struct {
	// User records user's name
	Tr   *http.Transport
	User string
}

func (server *Server) reverseHandler(req *http.Request) {
	req.Host = "127.0.0.1:8080"
	req.URL.Host = req.Host
	req.URL.Scheme = "http"
	logger.Infof("%v", req.RequestURI)
}

// ServeHTTP will be automatically called by system.
// ProxyServer implements the Handler interface which need ServeHTTP.
func (server *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, fmt.Sprintln(err))
		}
	}()

	server.reverseHandler(req)
	if req.Method == "CONNECT" {
		// proxy.HTTPSHandler(rw, req)
	} else {
		server.HTTPHandler(rw, req)
	}
}

// NewServer returns a new proxyserver.
func NewServer() *http.Server {
	return &http.Server{
		Addr:           ":8081",
		Handler:        &Server{Tr: &http.Transport{Proxy: http.ProxyFromEnvironment, DisableKeepAlives: true}},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

// HTTPHandler handles http connections.
// 处理普通的http请求
func (server *Server) HTTPHandler(rw http.ResponseWriter, req *http.Request) {
	logger.Infof("%v is sending request %v %v \n", server.User, req.Method, req.URL.Host)
	// RmProxyHeaders(req)
	resp, err := server.Tr.RoundTrip(req)
	if err != nil {
		logger.Errorf("%v", err)
		http.Error(rw, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	proxy.ClearHeader(rw.Header())
	proxy.CopyHeaders(rw.Header(), resp.Header)

	rw.WriteHeader(resp.StatusCode) //写入响应状态

	nr, err := ioCopy(rw, resp.Body)
	if err != nil && err != io.EOF {
		logger.Errorf("%v got an error when copy remote response to client. %v\n", server.User, err)
		return
	}
	logger.Infof("%v copied %v bytes from %v.\n", server.User, nr, req.URL.Host)
}

func ioCopy(dst io.Writer, src io.ReadCloser) (nr int, err error) {
	buf := make([]byte, 4096)
	defer func() {
		src.Close()
	}()

	var rerr, werr error
	var n int
	for {
		n, rerr = src.Read(buf)
		if n > 0 {
			_, werr = dst.Write(buf[:n])
			if flusher, ok := dst.(http.Flusher); ok {
				flusher.Flush()
			}
			nr += n
		}

		err = rerr
		if werr != nil {
			err = werr
		}
		if err != nil {
			return
		}
	}
}

func init() {
	http.DefaultServeMux.Handle("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, client"))
	}))
}

func main() {
	// web := NewWebServer()
	// go http.ListenAndServe(web.Port, web)
	// log.Println("Begin proxy")
	proxy := NewServer()
	logger.Fatal(proxy.ListenAndServe())
}
