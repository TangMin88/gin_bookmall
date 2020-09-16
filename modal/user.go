package modal

import "gin-bookmall/dao"

//User 用户
type User struct {
	ID       uint16 //自增的数
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Email    string `json:"email" form:"email"`
	Number   string `json:"number" form:"number"`
	Address  string `json:"address" form:"address"` //默认地址
}

//获取Book结构体
func GetUser() *User {
	user := &User{}
	return user
}

//AddUser 向数据库中添加用户
func (user *User) Add() error {
	return dao.Db.Create(user).Error
}

//QueryUser 根据用户名向数据库中查询一条记录
func (user *User) QueryU(userame string) error {
	return dao.Db.Where("username=?", userame).First(user).Error
}

//QueryUser 根据用户ID向数据库中查询一条记录
func (user *User) QuerysI(id uint16) error {
	return dao.Db.Where("id=?", id).First(user).Error
}

//UpdateUser 更新密码
func (user *User) Update() error {
	return dao.Db.Model(user).Updates(user).Error
}
