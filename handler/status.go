package handler

//HTTP407 异常抛出
var HTTP407 = []byte("HTTP/1.1 407 Proxy Authorization Required\r\nProxy-Authenticate: Basic realm=\"Toke is required\"\r\n\r\n")
