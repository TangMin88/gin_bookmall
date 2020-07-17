package modal

import (
	"gin-bookmall/dao"
	"time"
)

//Session 结构
type Session struct {
	ID       string //一个uuid生成的随机数
	UserName string //用户名
	UserID   int64  //用户id
	ShopName string
	ShopID   int64

	CreateTime time.Time //创建时间
}

//获取Session结构体
func GetSession(sessid string) *Session {
	session := &Session{
		ID: sessid,
	}
	return session
}

//AddSession 向数据库添加session
func (sess *Session) Add() error {
	return dao.Db.Create(sess).Error
}

//DeleteSession 向数据库删除session
func (sess *Session) Delete() error {
	return dao.Db.Where("id=?", sess.ID).Delete(Session{}).Error
}

//QuerySession 根据cookie的值查询session
func (sess *Session) Query() error {
	return dao.Db.Where("id=?", sess.ID).First(sess).Error
}

//根据用户的id查询session
func (sess *Session) QueryU(userid int64) error {
	return dao.Db.Where("user_id=?", userid).First(sess).Error
}
