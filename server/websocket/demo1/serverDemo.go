package main

import (
	"log"
	"net/http"

	wserver "gin-web-demo/server/websocket/demo1/server"
)

func main() {
	server := wserver.NewServer(":12345")

	// 设置连接的url
	server.WSPath = "/ws"
	// 设置推送消息的uri
	server.PushPath = "/push"
	// token校验
	server.AuthToken = func(token string) (userID string, ok bool) {

		// TODO: token校验
		//if token == "aaa" {
		//	return "jack", true
		//}
		//
		//return "", false

		return "jack", true
	}

	// 请求校验
	server.PushAuth = func(r *http.Request) bool {
		// TODO: check if request is valid
		log.Fatalln(r)
		return true
	}

	// Run server
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
