package processor

import (
	"fmt"
	"gin-bookmall/modal"
	"gin-bookmall/tool"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay"
)

//确认订单页面，提交订单
func Pay(c *gin.Context) {
	hander := c.Request.Header
	hs := hander.Get("user-agent")
	b := strings.Contains(hs, "Windows")
	// totalamount := c.PostForm("totalamount")
	// orderid := c.PostForm("orderid")
	totalamount := c.Query("totalamount")
	orderid := c.Query("orderid")
	if b { //电脑
		tool.WebPageAlipay(totalamount, orderid)
	} else { //手机
		tool.WapAlipay(totalamount, orderid)
	}
}

//ToCheckOut 去结账
func GetToCheckOut(c *gin.Context) {
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	//查询用户信息
	user := modal.GetUser()
	user.QuerysI(session.UserID)
	Orderdelevery := &modal.Orderdelevery{
		ReceiverName:    user.Username,
		ReceiverMobile:  user.Number,
		ReceiverAddress: user.Address,
	}
	_, cartitms := modal.Querys(session.Car.ID)
	session.Car.CartItms = cartitms
	session.Car.Totalcount = session.Car.GetTotalCount()
	session.Car.Totalamount = session.Car.GetTotalAmount()
	c.HTML(http.StatusOK, "confirm.html", gin.H{
		"car":  session.Car,
		"user": Orderdelevery,
	})
}

//QueryOrderStatus 我的订单，查询不同状态的订单
func GetOrder(c *gin.Context) {
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	//获取状态
	state := c.Query("state")
	var orders []*modal.Order
	var judge bool
	if state == "" {
		//查询全部订单
		orders, _ = modal.OrderQuerys(session.UserID, "user")
	} else {
		//查询不同状态的订单
		orders, _ = modal.OrderQueryU(session.UserID, state)
	}
	if len(orders) > 0 {
		judge = true
	}
	c.HTML(http.StatusOK, "myorder.html", gin.H{
		"orders": orders,
		"judge":  judge,
	})
}

//UpdateTheOrder 支付成功之后,创建/更新订单
func PostTheOrder(c *gin.Context) {
	fmt.Println("3333")
	var noti, _ = tool.Client.GetTradeNotification(c.Request)
	if noti != nil {
		fmt.Println("交易状态为:", noti.TradeStatus)
		// //确认订单是否在
		order := &modal.Order{
			ID: noti.OutTradeNo,
		}
		err := order.Query("4")
		if noti.TradeStatus == "TRADE_SUCCESS" || noti.TradeStatus == "TRADE_FINISHED" {
			order.State = "1"
			order.PaymentTime = time.Now()
			order.TradeNo = noti.TradeNo
			order.TotalAmount, _ = strconv.ParseFloat(noti.TotalAmount, 64)
		} else {
			order.PaymentTime = tool.Time()
		}
		if noti.TradeStatus == "WAIT_BUYER_PAY" {
			order.State = "4"
		}
		order.SubCode = noti.TradeStatus
		if err != nil {
			//判断是否是登录状态
			sess, ok := c.Get("sess")
			session := sess.(*modal.Session)
			if !ok {
				car := &modal.Car{
					ID: noti.OutTradeNo,
				}
				car.Querys()
				session.Car = car
			}
			//订单编号
			order.ID = noti.OutTradeNo
			order.CreateTime = time.Now()
			order.ConsignTime = tool.Time()
			order.ReceivingTime = tool.Time()
			order.UserID = session.Car.UserID
			order.ShopID = session.Car.ShopID
			_, cartItms := modal.Querys(order.ID)
			for _, v := range cartItms { //创建订单项
				orderitem := &modal.Orderitem{
					Count:   v.Count,
					Amount:  v.Amount,
					Title:   v.BookName,
					Price:   v.Price,
					Imgpath: v.Imgpath,
					OrderID: order.ID,
					BookID:  v.BookID,
				}
				in := &modal.Inventorie{ //更新当前图书的销量与库存
					ID: v.BookID,
				}
				sa, _ := in.UpdateSa(float64(v.Count))
				st, _ := in.UpdateSt(float64(-v.Count))
				in.Sales = uint16(sa)
				in.Stock = uint16(st)
				in.Update()
				orderitem.Add()
			}
			//删除购物车
			session.Car.Delete(order.ID, order.UserID)
			if ok {
				session.Car = nil
				session.Add()
			}
			order.Add()
		} else {
			order.Detele("4")
			err := order.Update()
			if err != nil {
				fmt.Println("通知页面", err)
			}
		}
		c.Next()
	}
	alipay.AckNotification(c.Writer) // 确认收到通知消息
}

//MyOrderState 订单页面，收货
func PutOrder(c *gin.Context) {
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	orderid := c.PostForm("orderid")
	order := &modal.Order{
		ID:     orderid,
		UserID: session.UserID,
	}
	err := order.Queryu("2")
	if err == nil && order.UserID == session.UserID {
		order.State = "3"
		//收货时间
		order.ReceivingTime = time.Now()
		//更新订单
		order.Update()
		c.JSON(http.StatusOK, "成功")
		return
	}
}

//变更图书
func ChangeBook(orderID string) {
	//订单取消后将图书的数量与销量更改回
	orderitems := modal.GetOrderitems()
	orderitems.Querys(orderID)
	for _, v := range orderitems {
		in := &modal.Inventorie{ //更新当前图书的销量与库存
			ID: v.BookID,
		}
		sa, _ := in.UpdateSa(float64(-v.Count))
		st, _ := in.UpdateSt(float64(v.Count))
		in.Sales = uint16(sa)
		in.Stock = uint16(st)
		in.Update()
	}
}

//我的订单，取消订单
func DeleteOrder(c *gin.Context) {
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	orderid := c.Query("orderid")
	order := &modal.Order{
		ID:     orderid,
		UserID: session.UserID,
	}
	order.Queryu("4")
	if order.State == "4" {
		//变更图书销量与数量
		ChangeBook(orderid)
		//删除订单
		modal.GetOrder().Detele("4")
		tool.TradeClose(order.ID) //关闭交易
		c.JSON(http.StatusOK, "成功")
	}
}

//退款
