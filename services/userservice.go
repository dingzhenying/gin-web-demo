package services

import (
	db "gin-web-demo/db"
	model "gin-web-demo/models"
)

func InsertUser(id int, name string, age int) model.User {
	u := model.User{Id: id, Name: name, Age: age}
	db.InsertUser(&u)
	return u
}
func GetUserList() []*model.User {
	users := db.GetUserList()
	return users
}
