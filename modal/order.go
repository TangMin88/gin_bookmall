package modal

import (
	"database/sql"
	"fmt"

	//"fmt"
	"gin-bookmall/dao"
	"time"
)

//Order 订单
type Order struct {
	ID            string       `json:"orderid" form:"orderid"` //订单号
	CreateTime    time.Time    //生成订单的时间
	TotalCount    int64        `gorm:"-"` //一份订单中图书的总数量
	TotalAmount   float64      `gorm:"-"` //一份订单中的总金额
	State         int          //订单的状态 1 未发货，2 已发货， 3交易完成 4未付款
	UserID        int64        //订单所属的用户
	OrderItem     []*Orderitem `gorm:"-"`
	ShopID        int64        //订单所属商店
	PaymentTime   time.Time    //付款时间
	ConsignTime   time.Time    //发货时间
	ReceivingTime time.Time    //收货时间
	BuyerRate     int          //1 已评价，  2 未评价
}

type Orders []*Order

//获取Order结构体
func GetOrder() *Order {
	order := &Order{}
	return order
}

//获取订单总数量
func (order *Order) GetOrderTotalCount() int64 {
	var totalCount int64
	//遍历购物车中的购物项
	for _, v := range order.OrderItem {
		totalCount = totalCount + v.Count
	}
	return totalCount
}

//获取订单总金额
func (order *Order) GetOrderTotalAmount() float64 {
	var TotalAmount float64
	for _, v := range order.OrderItem {
		TotalAmount = TotalAmount + v.Amount
	}
	return TotalAmount
}

//未发货
func (order *Order) NoSend() bool {
	return order.State == 1
}

//已发货
func (order *Order) SendComplate() bool {
	return order.State == 2
}

//未付款
func (order *Order) NotPaying() bool {
	return order.State == 4
}

//交易完成
func (order *Order) TheDeal() bool {
	return order.State == 3
}

//AddOrder添加订单
func (order *Order) Add() error {
	return dao.Db.Create(order).Error
}

//DeteleOrder 根据订单id删除订单
func (order *Order) Detele(id string) error {
	//先删除订单项
	err := GetOrderitem().Delete(id)
	if err != nil {
		return err
	}
	return dao.Db.Where("id = ?", id).Delete(Order{}).Error
}

//UpdateTheOrderState 根据订单id更新订单状态
func (order *Order) Update() error {
	return dao.Db.Table("orders").Where("id=?", order.ID).Updates(order).Error
}

//根据订单id查询订单
func (order *Order) Query(orderid string) error {
	return dao.Db.Where("id=? ", orderid).First(order).Error
}

//QueryOrder 根据店铺id查询全部订单
func OrderQuerySs(shopid int64) ([]*Order, error) {
	rows, err := dao.Db.Table("orders").Where("shop_id=? ", shopid).Order("create_time").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := ss(rows)
	return orders, nil
}

//QueryOrder 根据店铺id和状态查询订单
func OrderQueryS(shopid int64, state int) ([]*Order, error) {
	rows, err := dao.Db.Table("orders").Where("shop_id=? and state=?", shopid, state).Order("create_time").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := ss(rows)
	return orders, nil
}

//QueryOrder 根据用户id和状态查询订单
func OrderQueryU(userid int64, state int) ([]*Order, error) {
	rows, err := dao.Db.Table("orders").Where("user_id=? and state=?", userid, state).Order("create_time").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := ss(rows)
	return orders, nil
}

//QueryOrder 根据用户id查询全部订单
func OrderQueryUs(userid int64) ([]*Order, error) {
	rows, err := dao.Db.Table("orders").Where("user_id=? ", userid).Order("create_time").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := ss(rows)
	return orders, nil
}

func ss(rows *sql.Rows) []*Order {
	var orders []*Order
	for rows.Next() {
		order := &Order{}
		err := rows.Scan(&order.ID, &order.State, &order.UserID, &order.ShopID, &order.CreateTime, &order.PaymentTime, &order.ConsignTime, &order.ReceivingTime, &order.BuyerRate)
		if err != nil {
			fmt.Println("order_ss", err)
		}
		orderitem := GetOrderitems()
		orderitem.Querys(order.ID)
		order.OrderItem = orderitem
		order.TotalCount = order.GetOrderTotalCount()
		order.TotalAmount = order.GetOrderTotalAmount()
		orders = append(orders, order)
	}
	return orders
}
