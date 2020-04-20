package main

import (
	"gin-web-demo/db"
	_ "gin-web-demo/docs"
	router "gin-web-demo/routers"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func main() {
	//数据库初始化
	//todo 创建连接池
	db.InitDB()
	//defer 在正常执行后调用函数

	router := router.InitRouter()
	// 添加swagger页面
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(":8888")

}
