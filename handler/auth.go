package handler

import (
	"errors"
	"net/http"
)

//Auth 验证授权 这块可以用JWT做或者OAUTH
func (server *Server) Auth(rw http.ResponseWriter, req *http.Request) (string, bool) {
	if needAUTH() {
		msg, err := auth(rw, req)
		if err != nil {
			return "", false
		}
		return msg, true
	}
	return "", true
}

//needAUTH 判断当前请求是否开启授权
func needAUTH() bool {
	return true
}

func auth(rw http.ResponseWriter, req *http.Request) (path string, err error) {
	token := req.Header.Get("token")
	// TODO 验证授权密码这块可以根据自己的业务复杂度处理
	if token == "cc" {
		return "push", nil
	}
	return "", errors.New("TOKEN IS ERROR")
}
