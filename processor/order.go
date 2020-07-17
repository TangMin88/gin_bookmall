package processor

import (
	"gin-bookmall/modal"
	"gin-bookmall/tool"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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
	//查询购物车
	car := modal.GetCar()
	car.Query(session.UserID)
	c.HTML(http.StatusOK, "confirm.html", gin.H{
		"car":  car,
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
		orders, _ = modal.OrderQueryUs(session.UserID)
	} else {
		//查询不同状态的订单
		istate, _ := strconv.Atoi(state)
		orders, _ = modal.OrderQueryU(session.UserID, istate)
	}
	if len(orders) > 0 {
		judge = true
	}
	c.HTML(http.StatusOK, "myorder.html", gin.H{
		"orders": orders,
		"judge":  judge,
	})
}

//UpdateTheOrder 确认订单页面，根据是否付款，创建订单
func PostTheOrder(c *gin.Context) {
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	//获取状态
	state := c.Query("state")
	istate, _ := strconv.Atoi(state)
	car := modal.GetCar()
	car.Query(session.UserID)
	//购物车中一家一个订单编号
	for shopid, cartitm := range car.CartItms {
		//订单编号
		orderid := tool.UniqueID()
		order := modal.GetOrder()
		order.ID = orderid
		order.CreateTime = time.Now()
		order.State = istate
		toBeCharge := "0001-01-01 01:01:01"  //待转化为时间戳的字符串
		timeLayout := "2006-01-02 15:04:05"  //转化所需模板
		loc, _ := time.LoadLocation("Local") //获取时区
		theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc)
		order.ConsignTime = theTime
		order.ReceivingTime = theTime
		if istate == 1 {
			order.PaymentTime = time.Now()
		} else {
			order.PaymentTime = theTime
		}
		order.ShopID = shopid
		order.UserID = session.UserID
		order.Add()
		for _, v := range cartitm { //创建订单项，遍历购物车中的map中的购物项
			orderitem := &modal.Orderitem{
				Count:   v.Count,
				Amount:  v.Amount,
				Title:   v.Book.Title,
				Price:   v.Book.Price,
				Imgpath: v.Book.Imgpath,
				OrderID: orderid,
				BookID:  v.Book.ID,
			}
			book := v.Book //更新当前图书的销量与库存
			book.Sales = book.Sales + v.Count
			book.Stock = book.Stock - v.Count
			book.Update()
			orderitem.Add()
		}
	}
	//删除购物车
	car.Delete(car.ID)
	c.JSON(http.StatusOK, "成功")
}

//MyOrderState 我的订单，更新订单状态
func PutOrder(c *gin.Context) {
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	orderid := c.PostForm("orderid")
	state := c.PostForm("state")
	istate, _ := strconv.Atoi(state)
	order := modal.GetOrder()
	err := order.Query(orderid)
	if err == nil && order.UserID == session.UserID {
		order.State = istate
		switch istate {
		case 1: //付款时间
			order.PaymentTime = time.Now()
		case 3: //收货时间
			order.ReceivingTime = time.Now()
		}
		//更新订单状态
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
		book := modal.GetBook(v.BookID)
		book.Query()
		book.Sales = book.Sales - v.Count
		book.Stock = book.Stock + v.Count
		book.Update()
	}
}

//我的订单，取消订单
func DeleteOrder(c *gin.Context) {
	orderid := c.Query("orderid")
	//变更图书销量与数量
	ChangeBook(orderid)
	//删除订单
	modal.GetOrder().Detele(orderid)
	c.JSON(http.StatusOK, "成功")
}
