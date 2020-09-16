package modal

import "gin-bookmall/dao"

//OrderItem 订单中的订单项
type Orderitem struct {
	ID      uint16  `json:"Orderitemid,string"` //订单项的id
	Count   uint    `json:"count,string"`       //订单项中图书的数量
	Amount  float64 `json:"amount,string"`      //订单项中图书的金额小计
	Title   string  `json:"title"`              //订单项中图书的书名
	Price   float64 `json:"price,string"`       //订单项中图书的价格
	Imgpath string  `json:"imgpath"`            //订单项中图书的封面
	OrderID string  `json:"orderid"`            //订单项所属的订单
	BookID  uint16  `json:"bookid,string"`
}

type OrderitemS []*Orderitem

//获取Orderitem结构体
func GetOrderitem() *Orderitem {
	orderitem := &Orderitem{}
	return orderitem
}

//AddOrderItem添加订单项
func (orderitem *Orderitem) Add() error {
	return dao.Db.Create(orderitem).Error
}

//根据订单id删除订单项
func (orderitem *Orderitem) Delete(OrderID string) error {
	return dao.Db.Where("order_id = ?", OrderID).Delete(Orderitem{}).Error
}

//获取Orderitem结构体切片
func GetOrderitems() OrderitemS {
	var orderitems OrderitemS
	return orderitems
}

//根据订单号查询订单项
func (orderitems *OrderitemS) Querys(OrderID string) error {
	return dao.Db.Where("order_id=?", OrderID).Find(&orderitems).Error
}
