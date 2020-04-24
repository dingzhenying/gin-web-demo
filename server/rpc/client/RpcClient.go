package main

import (
	"flag"
	"fmt"
	functionModel "gin-web-demo/server/rpc"
	"log"
	"net/rpc"
	"strconv"
)

type Conf struct {
	serverAddress string
	i1            string
	i2            string
}

var conf = Conf{}

func SetConfiguration() {
	flag.StringVar(&conf.serverAddress, "server", "127.0.0.1:8999", "The address of the rpc")
	flag.StringVar(&conf.i1, "i1", "100", "100")
	flag.StringVar(&conf.i2, "i2", "2", "2")
}

func main() {
	SetConfiguration()
	flag.Parse()
	fmt.Println("severAddress = ", conf.serverAddress)

	// DelayHTTP在指定的网络地址连接到HTTP RPC服务器
	// 在默认HTTP RPC路径上监听。
	client, err := rpc.DialHTTP("tcp", conf.serverAddress)
	fmt.Println("client ====", client)
	if err != nil {
		log.Fatal("发生错误了 在这里地方  DialHTTP", err)
	}

	i1_, _ := strconv.Atoi(conf.i1)
	i2_, _ := strconv.Atoi(conf.i2)
	//指定出参类型（也可自定义）

	args := functionModel.Args{A: i1_, B: i2_}
	var reply int

	//调用调用命名函数，等待它完成，并返回其错误状态。
	err = client.Call("MathUtil.RpcFunction1", args, &reply)
	if err != nil {
		log.Fatal("调用RpcFunction1方法失败:", err)
	}
	fmt.Printf("Arith 乘法: %d*%d=%d\n", args.A, args.B, reply)
	//指定出参类型（也可自定义）
	var quot functionModel.Quotient
	//调用命名函数，等待它完成，并返回其错误状态。
	err = client.Call("MathUtil.RpcFunction2", args, &quot)
	if err != nil {
		log.Fatal("调用RpcFunction2方法失败:", err)
	}
	fmt.Printf("RpcFunction2 除法取整数: %d/%d=%d 余数 %d\n", args.A, args.B, quot.Quo, quot.Rem)
}
