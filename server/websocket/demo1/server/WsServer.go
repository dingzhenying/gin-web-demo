package server

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

const (
	serverDefaultWSPath   = "/ws"
	serverDefaultPushPath = "/push"
)

//默认配置
var defaultUpgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//是否跨域
	CheckOrigin: func(*http.Request) bool {
		return true
	},
}

// websocket server.
type Server struct {
	//服务端口（:80099）
	Addr string

	// 订阅地址，默认"/ws"
	WSPath string

	//消息推送地址，默认"/ws"
	PushPath string

	// Upgrader：gorilla/websocket中创建websocket连接的初始化对象
	// Upgrader为nil时，使用默认配置：ReadBufferSize（读取大小） 、 WriteBufferSize（写入大小）大小为 1024，CheckOrigin(跨域)为true

	Upgrader *websocket.Upgrader

	// 检查令牌是否有效并返回用户ID。如果令牌有效，则为userID
	// 必须返回并且确定应该为真。否则，确定应该为假。(定义方法)
	AuthToken func(token string) (userID string, ok bool)

	//授权推送请求。如果返回true，将发送消息，否则将丢弃该请求。默认nil和push请求
	//将始终被接受。
	PushAuth func(r *http.Request) bool
	//http服务中websocket方法体
	wh *websocketHandler
	//http服务中push方法体
	ph *pushHandler
}

// 创建监听server
func (s *Server) ListenAndServe() error {
	//binner结构体：userID与eventConne对象关系
	b := &binder{
		userID2EventConnMap: make(map[string]*[]eventConn),
		connID2UserIDMap:    make(map[string]string),
	}

	//websocket方法体
	wh := websocketHandler{
		upgrader: defaultUpgrader,
		binder:   b,
	}
	if s.Upgrader != nil {
		wh.upgrader = s.Upgrader
	}
	if s.AuthToken != nil {
		wh.calcUserIDFunc = s.AuthToken
	}
	s.wh = &wh
	http.Handle(s.WSPath, s.wh)

	// push request handler
	ph := pushHandler{
		binder: b,
	}
	if s.PushAuth != nil {
		ph.authFunc = s.PushAuth
	}
	s.ph = &ph
	http.Handle(s.PushPath, s.ph)

	return http.ListenAndServe(s.Addr, nil)
}

// 按照userId和event进行消息推送
func (s *Server) Push(userID, event, message string) (int, error) {
	return s.ph.push(userID, event, message)
}

// 通过userId,event删除连接
func (s *Server) Drop(userID, event string) (int, error) {
	return s.wh.closeConns(userID, event)
}

// 检查服务器参数，失败则返回错误。
func (s Server) check() error {
	if !checkPath(s.WSPath) {
		return fmt.Errorf("WSPath: %s not illegal", s.WSPath)
	}
	if !checkPath(s.PushPath) {
		return fmt.Errorf("PushPath: %s not illegal", s.PushPath)
	}
	if s.WSPath == s.PushPath {
		return errors.New("WSPath is equal to PushPath")
	}

	return nil
}

// 创建一个Server随想
func NewServer(addr string) *Server {
	return &Server{
		Addr:     addr,
		WSPath:   serverDefaultWSPath,
		PushPath: serverDefaultPushPath,
	}
}

func checkPath(path string) bool {
	if path != "" && !strings.HasPrefix(path, "/") {
		return false
	}
	return true
}
