package routers

import (
	co "gin-web-demo/controllers"
	"github.com/gin-gonic/gin"
)

/**
 *路由层，用于转发接口与方法调用
 */
func InitRouter() *gin.Engine {
	router := gin.Default()

	//Hello World
	router.GET("/hello", co.GetHelloInfo)

	router.POST("/user/add", co.InsertUser)

	router.GET("/user/getUserList", co.GetDataList)

	router.POST("/job/sub", co.SubmitJob)

	return router
}
