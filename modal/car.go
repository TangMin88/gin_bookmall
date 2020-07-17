package modal

import (
	"errors"
	"gin-bookmall/dao"
)

//import "fmt"

type Car struct {
	ID          string               `json:"carid" form:"carid"`   //购物车的id
	CartItms    map[int64][]*Cartitm `gorm:"-"`                    //购物车的所有购物项
	Totalcount  int64                `gorm:"-"`                    //图书的总数量，通过计算得到
	Totalamount float64              `gorm:"-"`                    //图书的总金额，通过计算得到
	UserID      int64                `json:"userid" form:"userid"` //用户id
}

//获取Car结构体
func GetCar() *Car {
	car := &Car{}
	return car
}

//GetMap 判断
func (car *Car) GetMap(c *Cartitm) {
	if car.CartItms == nil {
		car.CartItms = make(map[int64][]*Cartitm, 1)
	}
	value, ok := car.CartItms[c.Book.ShopID]
	if !ok {
		value = make([]*Cartitm, 0, 1)
	}
	value = append(value, c)
	car.CartItms[c.Book.ShopID] = value
	//fmt.Println(car.CartItms)
}

//获取图书总数量
func (car *Car) GetTotalCount() int64 {
	var totalCount int64
	//遍历购物车中的购物项
	for _, v := range car.CartItms {
		for _, v1 := range v {
			totalCount = totalCount + v1.Count
		}
	}
	return totalCount
}

//获取图书总金额
func (car *Car) GetTotalAmount() float64 {
	var TotalAmount float64
	for _, v := range car.CartItms {
		for _, v1 := range v {
			TotalAmount = TotalAmount + v1.Amount
		}
	}
	return TotalAmount
}

//AddCar 添加购物车
func (car *Car) Add() error {
	return dao.Db.Create(car).Error
}

//GetCarUserID 根据用户的id查对应的购物车
func (car *Car) Query(userid int64) error {
	err := dao.Db.Where("user_id=?", userid).First(car).Error
	if err != nil {
		return err
	}
	//查询购物项
	_, cartitms := Querys(car.ID)
	if len(cartitms) == 0 { //当有购物车却没有购物项时，删除购物车，避免购物时重复增加购物车
		car.Delete(car.ID)
		return errors.New("购物项为空")
	} else {
		for _, v := range cartitms {
			book := GetBook(v.BookID)
			book.Query()
			v.Book = book
			car.GetMap(v)
		}
	}
	car.Totalcount = car.GetTotalCount()
	car.Totalamount = car.GetTotalAmount()
	return nil
}

//UpdateCar 更新购物车
func (car *Car) Update() error {
	return dao.Db.Model(car).Updates(car).Error
}

//DeleteCar 根据购物车id删除购物车
func (car *Car) Delete(id string) error {
	//删除购物车前先删除对应的购物项
	err := GetCartitm().Deletes(id)
	if err != nil {
		return err
	}
	return dao.Db.Where("id = ?", id).Delete(Car{}).Error
}
