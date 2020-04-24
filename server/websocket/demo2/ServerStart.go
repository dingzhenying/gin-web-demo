package main

import (
	"fmt"
	impl "gin-web-demo/server/websocket/demo2/server"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte("hello"))
	var (
		wsConn *websocket.Conn
		err    error
		conn   *impl.Connection
		data   []byte
	)
	// 完成ws协议的握手操作
	// Upgrade:websocket（允许跨域的websocket）
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	if conn, err = impl.InitConnection(wsConn); err != nil {
		goto ERR
	}

	// 启动线程，不断发消息
	//go func(){
	//	var (err error)
	//	var i =1
	//	for{
	//		if err = conn.WriteMessage([]byte("heartbeat"+strconv.Itoa(i)));err != nil{
	//			return
	//		}
	//		time.Sleep(1*time.Second)
	//		i++
	//	}
	//}()
	//轮询读取数据(读取到数据处理)
	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		fmt.Println(string(data))
		//todo 获取到数据处理
		if err = conn.WriteMessage([]byte("read Mess: " + string(data))); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()

}

func main() {
	//创建接口服务
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe("0.0.0.0:7777", nil)
}
