package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func Test_Server_1(t *testing.T) {
	port := 12345
	userID := uuid.New().String()
	event := "e1"

	count := 100
	//创建server对象
	server := NewServer(":" + strconv.Itoa(port))

	//启动写成注册和启动服务
	go runBasicWServer(server)
	time.Sleep(time.Millisecond * 300)

	// 注册消息的结构体
	registerCh := make(chan struct{})
	//发送消息的结构体
	pushCh := make(chan struct{})
	//创建写成发送消息
	go func() {
		//管道
		<-registerCh
		for i := 0; i < count; i++ {
			msg := fmt.Sprintf("hello -- %d", i)
			_, err := server.Push(userID, event, msg)
			if err != nil {
				t.Errorf("push message fail: %v", err)
			}

			// time.Sleep(time.Microsecond)
		}

		if _, err := server.Drop(userID, ""); err != nil {
			t.Errorf("drop connections fail: %v", err)
		}

		close(pushCh)
	}()

	// 注册用户信息
	url := fmt.Sprintf("ws://127.0.0.1:%d/ws", port)
	//默认值连接websocker 服务
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("dial websocket url fail: %v", err)
	}

	// register
	rm := RegisterMessage{
		Token: userID,
		Event: event,
	}
	if err := conn.WriteJSON(rm); err != nil {
		t.Fatalf("registe fail: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	close(registerCh)

	// read
	cnt := 0
	for {
		_, r, err := conn.NextReader()
		if err != nil {
			//t.Errorf(err.Error())
			break
		}
		b, _ := ioutil.ReadAll(r)
		t.Logf("msg: %s", string(b))

		cnt++
	}

	if cnt != count {
		t.Errorf("message received count: %d, expected %d", cnt, count)
	}
}

//创建server监听
func runBasicWServer(s *Server) {
	if err := s.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
