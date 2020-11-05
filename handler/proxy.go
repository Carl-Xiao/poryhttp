package handler

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/logger"
)

// Server is a server of proxy.
type Server struct {
	Tr *http.Transport
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

	if !server.Auth(rw, req) {
		rw.WriteHeader(407)
		rw.Write(HTTP407)
		return
	}
	//TODO 依据token值选择后台负载均衡的方式
	req.URL.Path = "/push" + req.URL.Path
	server.LoadBalancing(req)
	defer server.Done(req)

	server.HTTPHandler(rw, req)
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
	logger.Infof("sending request %v %v \n", req.Method, req.URL.Host)
	RmProxyHeaders(req)
	resp, err := server.Tr.RoundTrip(req)
	if err != nil {
		logger.Errorf("%v", err)
		http.Error(rw, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	ClearHeader(rw.Header())
	CopyHeaders(rw.Header(), resp.Header)

	rw.WriteHeader(resp.StatusCode) //写入响应状态

	nr, err := io.Copy(rw, resp.Body)
	if err != nil && err != io.EOF {
		logger.Errorf("got an error when copy remote response to client. %v\n", err)
		return
	}
	logger.Infof("copied %v bytes from %v.\n", nr, req.URL.Host)
}
