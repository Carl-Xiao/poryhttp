package handler

import (
	"net/http"

	"github.com/lafikl/liblb/r2"
)

var r2LB *r2.R2

var serverNodes []string

func init() {
	//serverNodes后端服务器的节点
	serverNodes = append(serverNodes, "127.0.0.1:9070")
	r2LB = r2.New(serverNodes...)
}

//LoadBalancing handles request for reverse proxy.
func (ps *Server) LoadBalancing(req *http.Request) {
	ps.loadBalancing(req)
}

//Done handles request for reverse proxy.
func (ps *Server) Done(req *http.Request) {

}

func (ps *Server) loadBalancing(req *http.Request) {
	var proxyHost string
	// Selects a back-end server base on polling algorithm which supports weight.
	proxyHost, _ = r2LB.Balance()
	req.Host = proxyHost
	req.URL.Host = proxyHost
	req.URL.Scheme = "http"
}
