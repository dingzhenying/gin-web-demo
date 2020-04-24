package main

import (
	"fmt"
	model "gin-web-demo/server/rpc"
	"log"
	//model "gin-web-demo/server/rpc"
	"net"
	"net/http"
	"net/rpc"
)

func rpcServer() {
	//1、初始化（服务对象）指针数据类型
	mathUtil := new(model.MathUtil) //初始化指针数据类型
	fmt.Println("mathUtil=====", mathUtil)

	//2、调用net/rpc包的功能将服务对象进行注册
	err := rpc.Register(mathUtil)
	if err != nil {
		log.Fatalln(err.Error())
	}

	//3、通过该函数把mathUtil中提供的服务注册到HTTP协议上，方便调用者可以利用http的方式进行数据传递
	rpc.HandleHTTP()

	//4、在特定的端口进行监听
	listen, err := net.Listen("tcp", ":8999")
	if err != nil {
		log.Fatalln(err.Error())
		//panic(err.Error())
	}
	//service接受侦听器l上传入的HTTP连接，
	http.Serve(listen, nil)
}

func main() {
	rpcServer()

}
