package modal

import "gin-bookmall/dao"

//OrderItem 订单中的订单项
type Orderitem struct {
	ID      int64   //订单项的id
	Count   int64   //订单项中图书的数量
	Amount  float64 //订单项中图书的金额小计
	Title   string  //订单项中图书的书名
	Price   float64 //订单项中图书的价格
	Imgpath string  //订单项中图书的封面
	OrderID string  //订单项所属的订单
	BookID  int64
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
