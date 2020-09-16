package modal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gin-bookmall/dao"
	"strconv"
	"time"
)

//Order 订单
type Order struct {
	ID            string       `json:"orderid" form:"orderid"` //订单号
	CreateTime    time.Time    `json:"createtime,string" `     //生成订单的时间
	TotalCount    uint         `json:"totalcount,string"`      //一份订单中图书的总数量
	TotalAmount   float64      `json:"totalamount,string"`     //一份订单中的总金额
	State         string       `json:"state"`                  //订单的状态 1 未发货，2 已发货， 3交易完成 WAIT_BUYER_PAY未付款 https://opendocs.alipay.com/open/59/103672
	UserID        uint16       `json:"user,string"`            //订单所属的用户
	OrderItem     []*Orderitem `gorm:"-" json:"orderitem,string"`
	ShopID        uint16       `json:"shopid"`               //订单所属商店
	PaymentTime   time.Time    `json:"paymenttime,string"`   //付款时间
	ConsignTime   time.Time    `json:"consigntime,string"`   //发货时间
	ReceivingTime time.Time    `json:"receivingtime,string"` //收货时间
	BuyerRate     string       `json:"buyerrate"`            //1 已评价，  2 未评价
	TradeNo       string       `json:"tradeno"`              //支付宝交易号
	SubCode       string       `json:"subcode"`              //交易状态
}

type Orders []*Order

//获取Order结构体
func GetOrder() *Order {
	order := &Order{}
	return order
}

//获取订单总数量
func (order *Order) GetOrderTotalCount() uint {
	var totalCount uint
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
	return order.State == "1"
}

//已发货
func (order *Order) SendComplate() bool {
	return order.State == "2"
}

//未付款
func (order *Order) NotPaying() bool {
	return order.State == "4"
}

//交易完成
func (order *Order) TheDeal() bool {
	return order.State == "3"
}

//AddOrder添加订单
func (order *Order) Add() error {
	if order.State != "4" {
		err := dao.Db.Create(order).Error
		if err != nil {
			return err
		}
	}
	ku := JointStr("orders", "user", strconv.FormatUint(uint64(order.UserID), 10), order.State)
	ks := JointStr("orders", "shop", strconv.FormatUint(uint64(order.ShopID), 10), order.State)
	return order.RAdd(ku, ks)
}

func (order *Order) RAdd(k ...string) error {
	pipel := dao.Rdb.Pipeline()
	str, _ := json.Marshal(order)
	for _, v := range k {
		pipel.HSet(v, order.ID, str)
	}
	if order.State == "4" {
		pipel.Set(order.ID, "未付款", 1800)
	}
	_, err := pipel.Exec()
	return err
}

//DeteleOrder 根据订单id删除订单
func (order *Order) Detele(state string) error {
	pipel := dao.Rdb.Pipeline()
	if state == "4" {
		pipel.Del(order.ID)
	}
	ku := JointStr("orders", "user", strconv.FormatUint(uint64(order.UserID), 10), state)
	ks := JointStr("orders", "shop", strconv.FormatUint(uint64(order.ShopID), 10), state)
	pipel.HDel(ku, order.ID)
	pipel.HDel(ks, order.ID)
	_, err := pipel.Exec()
	return err
}

//UpdateTheOrderState 根据订单id更新订单
func (order *Order) Update() error {
	err := dao.Db.Table("orders").Where("id=?", order.ID).Updates(order).Error
	if err != nil {
		return err
	}
	ku := JointStr("orders", "user", strconv.FormatUint(uint64(order.UserID), 10), order.State)
	ks := JointStr("orders", "shop", strconv.FormatUint(uint64(order.ShopID), 10), order.State)
	return order.RAdd(ku, ks)
}

//根据订单id查询店铺订单是否存在
func (order *Order) Querys() error {
	ks := JointStr("orders", "shop", strconv.FormatUint(uint64(order.ShopID), 10), "1")
	_, err := dao.Rdb.HGet(ks, order.ID).Result()
	if err != nil {
		return dao.Db.Where("id=? and shop_id = ?", order.ID, order.ShopID).First(order).Error

	}
	return nil
}

//查询订单是否存在
func (order *Order) Query(state string) error {
	return dao.Db.Where("id=? and state = ?", order.ID, state).First(order).Error
}

//根据订单id查询用户订单是否存在
func (order *Order) Queryu(state string) error {
	ks := JointStr("orders", "user", strconv.FormatUint(uint64(order.UserID), 10), state)
	_, err := dao.Rdb.HGet(ks, order.ID).Result()
	if err != nil {
		return dao.Db.Where("id=? and shop_id = ?", order.ID, order.UserID).First(order).Error

	}
	return nil
}

//QueryOrder 查询全部订单  id:用户/shop  l:shop/user
func OrderQuerys(id uint16, l string) ([]*Order, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if l == "shop" {
		rows, err = dao.Db.Table("orders").Where("shop_id=? ", id).Order("create_time").Rows()
	} else {
		rows, err = dao.Db.Table("orders").Where("user_id=? ", id).Order("create_time").Rows()
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []*Order
	for rows.Next() {
		order := &Order{}
		err := rows.Scan(&order.ID, &order.State, &order.UserID, &order.ShopID, &order.CreateTime, &order.PaymentTime, &order.ConsignTime, &order.ReceivingTime, &order.BuyerRate, &order.TradeNo, &order.TotalCount, &order.TotalAmount, &order.SubCode)
		if err != nil {
			fmt.Println("order_ss", err)
		}
		orderitem := GetOrderitems()
		orderitem.Querys(order.ID)
		order.OrderItem = orderitem
		orders = append(orders, order)
	}
	return orders, nil
}

//QueryOrder 根据店铺id和状态查询订单
func OrderQueryS(shopid uint16, state string) ([]*Order, error) {
	ks := JointStr("orders", "shop", strconv.FormatUint(uint64(shopid), 10), state)
	sstr, err := dao.Rdb.HVals(ks).Result()
	var orders []*Order
	if err != nil {
		rows, err := dao.Db.Table("orders").Where("shop_id=? and state=?", shopid, state).Order("create_time").Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		orders = ss(rows, "shop")
	} else {
		for _, v := range sstr {
			var order *Order
			err := json.Unmarshal([]byte(v), order)
			if err != nil {
				fmt.Println("OrderQueryS unmarshal err ", err)
			}
			orders = append(orders, order)
		}
	}

	return orders, nil
}

//QueryOrder 根据用户id和状态查询订单
func OrderQueryU(userid uint16, state string) ([]*Order, error) {
	ku := JointStr("orders", "user", strconv.FormatUint(uint64(userid), 10), state)
	sstr, err := dao.Rdb.HVals(ku).Result()
	var orders []*Order
	if err != nil {
		rows, err := dao.Db.Table("orders").Where("user_id=? and state=?", userid, state).Order("create_time").Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		orders = ss(rows, "user")
	} else {
		for _, v := range sstr {
			var order *Order
			err := json.Unmarshal([]byte(v), order)
			if err != nil {
				fmt.Println("OrderQueryS unmarshal err ", err)
			}
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func ss(rows *sql.Rows, k string) []*Order {
	var orders []*Order
	pipel := dao.Rdb.Pipeline()
	var kk string
	for rows.Next() {
		order := &Order{}
		err := rows.Scan(&order.ID, &order.State, &order.UserID, &order.ShopID, &order.CreateTime, &order.PaymentTime, &order.ConsignTime, &order.ReceivingTime, &order.BuyerRate, &order.TradeNo, &order.TotalCount, &order.TotalAmount, &order.SubCode)
		if err != nil {
			fmt.Println("order_ss", err)
		}
		orderitem := GetOrderitems()
		orderitem.Querys(order.ID)
		order.OrderItem = orderitem
		str, _ := json.Marshal(order)
		if k == "user" {
			kk = JointStr("orders", "user", strconv.FormatUint(uint64(order.UserID), 10), order.State)
		} else {
			kk = JointStr("orders", "shop", strconv.FormatUint(uint64(order.ShopID), 10), order.State)
		}
		pipel.HSet(kk, order.ID, str)
		orders = append(orders, order)
	}
	_, err := pipel.Exec()
	if err != nil {
		fmt.Println("ss pipel err", err)
	}
	return orders
}
