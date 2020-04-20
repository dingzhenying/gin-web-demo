package controllers

import (
	"encoding/json"
	service "gin-web-demo/services"
	mesg "gin-web-demo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @接口测试
// @Description get data
// @Accept  json
// @Produce json
// @Param dataInfo query string true "string valid"
// @Success 200 {string} string "hello"
// @Router /hello/ [get]
func GetHelloInfo(c *gin.Context) {
	dataInfo := c.Request.FormValue("dataInfo")
	outData := service.HelloService(dataInfo)
	//返回结果
	c.JSON(http.StatusOK, mesg.GetSuccessMsg(outData))
}

// @增加用户
// @Description add user
// @Accept  json
// @Produce json
// @Param Id query int false "int valid"
// @Param Name query string false "string valid"
// @Param Age query int false "int valid"
// @Router /user/add [post]
func InsertUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("Id"))
	name := c.Request.FormValue("Name")
	age, _ := strconv.Atoi(c.Request.FormValue("Age"))
	data := service.InsertUser(id, name, age)
	jobinfo, _ := json.Marshal(data)

	//返回结果
	c.JSON(http.StatusOK, gin.H{
		"data": string(jobinfo),
	})
}

// @获取用户列表
// @Description get data
// @Accept  json
// @Produce json
// @Success 200 {string} string "userList"
// @Router /user/getUserList [get]
func GetDataList(c *gin.Context) {
	outData := service.GetUserList()
	//返回结果
	c.JSON(http.StatusOK, mesg.GetSuccessMsg(outData))
}

// @提交job任务
// @Description submit lsf job
// @Accept  json
// @Produce json
// @Param jobName query string false "string valid"
// @Param jobinfo query string false "string valid"
// @Param jobData query string false "string valid"
// @Router /job/sub [post]
func SubmitJob(c *gin.Context) {
	jobName := c.Request.FormValue("jobName")
	jobinfo := c.Request.FormValue("jobinfo")
	jobData := c.Request.FormValue("jobData")
	data, _ := service.SubmitJob(jobName, jobinfo, jobData)
	//返回结果
	c.JSON(http.StatusOK, mesg.GetSuccessMsg(data))

}
