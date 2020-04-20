package utils

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	UserId   int `gorm:"primary_key"`
	Phone    string
	WxopenId string
	Tcreate  *time.Time
	Tprocess *time.Time
	Balance  int
	Src      string
	Level    int
}

func main() {
	//连接数据库
	db, err := gorm.Open("mysql", "root:123456@tcp(192.168.66.162:3306)/go_test?charset=utf8")
	//一个坑，不设置这个参数，gorm会把表名转义后加个s，导致找不到数据库的表
	db.SingularTable(true)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	var user User
	if !db.HasTable(user) {
		db.CreateTable(user)
	}

	var phone = "12345678900"
	//条件查询
	err = db.Where("phone = ?", phone).Find(&user).Error
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(user.UserId)
	//把查询出来的一条数据删除
	err = db.Delete(&user).Error
	if err != nil {
		fmt.Println(err)
	}
}
