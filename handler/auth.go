package handler

import (
	"errors"
	"net/http"
)

//Auth 验证授权 这块可以用JWT做或者OAUTH
func (server *Server) Auth(rw http.ResponseWriter, req *http.Request) bool {
	if needAUTH() {
		if err := auth(rw, req); err != nil {
			return false
		}
		return true
	}
	return true
}

//needAUTH 判断当前请求是否开启授权
func needAUTH() bool {
	return true
}

func auth(rw http.ResponseWriter, req *http.Request) (err error) {
	token := req.Header.Get("token")
	// TODO 验证授权密码这块可以根据自己的业务复杂度处理

	if token == "cc" {
		return nil
	}
	return errors.New("TOKEN IS ERROR")
}
