package modal

import (
	"encoding/json"
	"fmt"
	"gin-bookmall/dao"
	"time"
)

//Session 结构
type Session struct {
	ID       string `json:"id"`            //一个uuid生成的随机数
	UserName string `json:"username"`      //用户名
	UserID   uint16 `json:"userid,string"` //用户id
	ShopName string `json:"shopname"`
	ShopID   uint16 `json:"shopid,string"`
	Car      *Car   `json:"carid"`
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
	b, _ := json.Marshal(sess)
	err := dao.Rdb.Set(sess.ID, b, 15*60*time.Second).Err()
	if err != nil {
		fmt.Println("set err", err)
	}
	return err
}

//DeleteSession 向数据库删除session
func (sess *Session) Delete() error {
	return dao.Rdb.Del(sess.ID).Err()
}

//QuerySession 根据cookie的值查询session
func (sess *Session) Query() error {
	str, err := dao.Rdb.Get(sess.ID).Result()
	if err != nil {
		fmt.Printf("Query() Get err:%v\n", err)
		return err
	}
	return json.Unmarshal([]byte(str), sess)
}
