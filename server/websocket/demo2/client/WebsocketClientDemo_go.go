package main

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var origin = "http://127.0.0.1:7777/"
var url = "ws://127.0.0.1:7777/ws"
var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func main() {
}
