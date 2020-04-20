package db

import (
	"encoding/json"
	"fmt"
	model "gin-web-demo/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var db *gorm.DB
var user *model.User
var job *model.LSFJobReq

func InitDB() {
	fmt.Println("InitDB...")
	//初始化数据库
	db = createTables(NewConn())
}

func NewConn() *gorm.DB {
	db, err := gorm.Open("mysql", "root:123456@tcp(192.168.66.162:3306)/go_test?charset=utf8")
	if err != nil {
		log.Fatalln(err)
	}
	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(20)
	createTables(db)

	fmt.Println("new conn...")
	return db
}

func createTables(db *gorm.DB) *gorm.DB {

	if !db.HasTable(user) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(user).Error; err != nil {
			panic(err)
		}
	}
	if !db.HasTable(job) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(job).Error; err != nil {
			panic(err)
		}
	}
	return db

}

func InsertUser(u *model.User) {
	use, _ := json.Marshal(u)
	fmt.Println(string(use))
	fmt.Println("插入用户...")
	db.Create(&u)
}
func GetUserList() []*model.User {
	var users []*model.User
	fmt.Println("查询用户列表...")
	db.Find(&users)
	return users
}
func InsertJob(job *model.LSFJobReq) {
	jobinfo, _ := json.Marshal(job)

	fmt.Println(string(jobinfo))
	fmt.Println("插入job...")
	db.Create(&job)
}
