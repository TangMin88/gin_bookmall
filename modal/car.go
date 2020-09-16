package modal

import (
	//"errors"
	"gin-bookmall/dao"
	//"fmt"
	"encoding/json"
	"strconv"
)

//import "fmt"

type Car struct {
	ID          string     `json:"carid" form:"carid"`          //购物车的id
	CartItms    []*Cartitm `gorm:"-" json:"-"`                  //购物车的所有购物项
	Totalcount  uint       `gorm:"-" json:"-"`                  //图书的总数量，通过计算得到
	Totalamount float64    `gorm:"-" json:"-"`                  //图书的总金额，通过计算得到
	UserID      uint16     `json:"userid,string" form:"userid"` //用户id
	ShopID      uint16     `json:"shopid,string"`
}

//获取Car结构体
func GetCar() *Car {
	car := &Car{}
	return car
}

//获取图书总数量
func (car *Car) GetTotalCount() uint {
	var totalCount uint
	//遍历购物车中的购物项
	for _, v := range car.CartItms {
		totalCount = totalCount + v.Count
	}
	return totalCount
}

//获取图书总金额
func (car *Car) GetTotalAmount() float64 {
	var TotalAmount float64
	for _, v := range car.CartItms {
		TotalAmount = TotalAmount + v.Amount
	}
	return TotalAmount
}

//添加购物车
func (car *Car) Add() error {
	err := dao.Db.Create(car).Error
	if err != nil {
		return err
	}
	k := JointStr("cars", strconv.FormatUint(uint64(car.UserID), 10))
	return SAdd(k, car)
}

func (car *Car) Querys() error {
	return dao.Db.Where("id=?", car.ID).First(car).Error
}

//根据用户的id查对应的购物车
func (car *Car) Query(userid uint16) error {
	k := JointStr("cars", strconv.FormatUint(uint64(userid), 10))
	str, err := dao.Rdb.Get(k).Result()
	if err != nil {
		err := dao.Db.Where("user_id=?", userid).First(car).Error
		if err != nil {
			return err
		}
		SAdd(k, car)
	} else {
		json.Unmarshal([]byte(str), car)
	}
	return nil
}

//根据购物车id删除购物车
func (car *Car) Delete(carid string, userid uint16) error {
	//删除购物车前先删除对应的购物项
	err := GetCartitm().Deletes(carid)
	if err != nil {
		return err
	}
	err = dao.Db.Where("id = ?", carid).Delete(Car{}).Error
	if err != nil {
		return err
	}
	k := JointStr("cars", strconv.FormatUint(uint64(userid), 10))
	err = dao.Rdb.Del(k).Err()
	if err != nil {
		return err
	}
	return nil
}
