package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// “/ws” 的连接调用的handler方法
type websocketHandler struct {
	// upgrader：原生websocket中的配置对象(用于调用连接的websocket服务)
	upgrader *websocket.Upgrader

	//用的存储userID 和 eventConn关系的结构体（连接信息的数据库）、
	//todo 可以使用外部数据存储校验信息
	binder *binder

	//token是否有效
	calcUserIDFunc func(token string) (userID string, ok bool)
}

// RegisterMessage：连接到服务器校验信息的结构体
type RegisterMessage struct {
	Token string
	Event string
}

// 初始化websocketHandler，创建websocket连接
func (wh *websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsConn, err := wh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer wsConn.Close()

	// 处理websocket请求
	conn := NewConn(wsConn)
	conn.AfterReadFunc = func(messageType int, r io.Reader) {
		var rm RegisterMessage
		decoder := json.NewDecoder(r)
		if err := decoder.Decode(&rm); err != nil {
			return
		}

		// calculate userID by token
		userID := rm.Token
		if wh.calcUserIDFunc != nil {
			uID, ok := wh.calcUserIDFunc(rm.Token)
			if !ok {
				return
			}
			userID = uID
		}

		// 绑定
		wh.binder.Bind(userID, rm.Event, conn)
	}
	conn.BeforeCloseFunc = func() {
		// 解绑
		wh.binder.Unbind(conn)
	}
	//监听
	conn.Listen()
}

// 关闭连接信息的方法体（关闭连接信息，删除存储在binner中的token信息）
func (wh *websocketHandler) closeConns(userID, event string) (int, error) {
	conns, err := wh.binder.FilterConn(userID, event)
	if err != nil {
		return 0, err
	}

	cnt := 0
	for i := range conns {
		// unbind
		if err := wh.binder.Unbind(conns[i]); err != nil {
			log.Printf("conn unbind fail: %v", err)
			continue
		}

		// close
		if err := conns[i].Close(); err != nil {
			log.Printf("conn close fail: %v", err)
			continue
		}

		cnt++
	}

	return cnt, nil
}

// 连接失败范返回的信息
var ErrRequestIllegal = errors.New("request data illegal")

// “/push” 调用推送消息的接口的方法体
type pushHandler struct {
	// 授权信息状态，为true时“/push”才会发送信息
	authFunc func(r *http.Request) bool

	binder *binder
}

// HTTP服务
func (s *pushHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		//post方法校验
		w.WriteHeader(http.StatusMethodNotAllowed) //405
		return
	}

	// 授权校验
	if s.authFunc != nil {
		if ok := s.authFunc(r); !ok {
			w.WriteHeader(http.StatusUnauthorized) //401
			return
		}
	}

	// 读取数据
	var pm PushMessage
	decoder := json.NewDecoder(r.Body)
	//json 转存到 PushMessage结构体中
	if err := decoder.Decode(&pm); err != nil {
		w.WriteHeader(http.StatusBadRequest) //400
		w.Write([]byte(ErrRequestIllegal.Error()))
		return
	}

	// 连接信息校验
	if pm.UserID == "" || pm.Event == "" || pm.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrRequestIllegal.Error()))
		return
	}
	//推送消息给指定用户
	cnt, err := s.push(pm.UserID, pm.Event, pm.Message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result := strings.NewReader(fmt.Sprintf("message sent to %d clients", cnt))
	io.Copy(w, result)
}

func (s *pushHandler) push(userID, event, message string) (int, error) {
	if userID == "" || event == "" || message == "" {
		return 0, errors.New("parameters(userId, event, message) can't be empty")
	}

	// 获取指定userID+event获取conns连接
	conns, err := s.binder.FilterConn(userID, event)
	if err != nil {
		return 0, fmt.Errorf("filter conn fail: %v", err)
	}
	cnt := 0
	for i := range conns {
		_, err := conns[i].Write([]byte(message))
		if err != nil {
			s.binder.Unbind(conns[i])
			continue
		}
		cnt++
	}

	return cnt, nil
}

// 发送消息的结构体
type PushMessage struct {
	UserID  string `json:"userId"`
	Event   string
	Message string
}
