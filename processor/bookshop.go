package processor

import (
	"fmt"
	"gin-bookmall/modal"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//AsTheOwner 成为店主
func AsTheOwner(c *gin.Context) {
	c.HTML(http.StatusOK, "owner.html", nil)
}

//PostAsTheOwner 提交成为店主表单
func PostAsTheOwner(c *gin.Context) {
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	if session.ShopID == 0 { //查询用户是否有店铺，避免重复
		shopName := c.PostForm("shopname")
		shop := modal.GetShop()
		shop.ShopName = shopName
		shop.UserID = session.UserID
		shop.Add()
	}
	MyBookShop(c)
}

//MyBookShop 我的店铺
func MyBookShop(c *gin.Context) {
	//获取当前页页码
	page := &modal.Page{}
	if err := c.ShouldBind(page); err == nil {
		if page.PageNo == 0 {
			page.PageNo = 1
		}
		sess, _ := c.Get("sess")
		session := sess.(*modal.Session)
		//获取当前页图书
		err := page.QueryTotalS(session.ShopID)
		judge := false
		if err == nil {
			judge = true
		}
		c.HTML(http.StatusOK, "myshop.html", gin.H{
			"page":     page,
			"judge":    judge,
			"shopname": session.ShopName,
			//"error":  err.Error(),
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

//MyInvoice 我的货单
func GetInvoice(c *gin.Context) {
	state := c.Query("state")
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	var orders []*modal.Order
	var judge bool
	if state == "" {
		orders, _ = modal.OrderQuerySs(session.ShopID)
	} else { //查询不同状态的订单
		istate, _ := strconv.Atoi(state)
		orders, _ = modal.OrderQueryS(session.ShopID, istate)
	}
	if len(orders) > 0 {
		judge = true
	}
	c.HTML(http.StatusOK, "myinvoice.html", gin.H{
		"orders": orders,
		"judge":  judge,
	})
}

//TheDelivery 发货
func PutTheDelivery(c *gin.Context) {
	//order := modal.GetOrder()
	orderid := c.Query("orderid")
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	if orderid != "" {
		//查询订单id是否属于登录的店铺
		order := modal.GetOrder()
		err := order.Query(orderid)
		if err == nil && order.ShopID == session.ShopID {
			order.ConsignTime = time.Now()
			order.State = 2
			//更新订单状态
			order.Update()
		}
	}
	c.JSON(http.StatusOK, "成功")
}

//CheckTheDetails 查看详情
// func CheckTheDetails(c *gin.Context) {
// 	orderid := c.Query("orderid")
// 	shop := c.Query("shopname")
// 	if orderid != "" {
// 		//查询订单的订单项
// 		orderitems := modal.GetOrderitems()
// 		err := orderitems.Querys(orderid)
// 		if err == nil {
// 			//查询用户信息
// 			userid, _ := c.Get("userid")
// 			user := modal.GetUser()
// 			user.QuerysI(userid.(int64))
// 			c.HTML(http.StatusOK, "checkThedetails.html", gin.H{
// 				"shop":       shop,
// 				"user":       user,
// 				"orderitems": orderitems,
// 			})
// 		} else {
// 			//返回我的订单页面
// 			MyInvoice(c)
// 		}
// 	} else {
// 		//返回我的订单页面
// 		MyInvoice(c)
// 	}
// }

//DleteShopBook 删除书籍
func DleteBook(c *gin.Context) {
	id := c.Query("bookid")
	iid, _ := strconv.ParseInt(id, 10, 64)
	//查询要删除的书籍是否属于登录的书店
	book := modal.GetBook(iid)
	err := book.Query()
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	if err == nil && book.ShopID == session.ShopID {
		book.Delete()
		c.JSON(http.StatusOK, "成功")
	} else {
		c.JSON(http.StatusOK, "删除书籍失败")
	}
}

//添加/修改书籍
func GetBook(c *gin.Context) {
	var j bool
	bookid := c.Query("bookid")
	if bookid != "" {
		//修改
		ibookid, _ := strconv.ParseInt(bookid, 10, 64)
		book := modal.GetBook(ibookid)
		err := book.Query()
		sess, _ := c.Get("sess")
		session := sess.(*modal.Session)
		if err == nil && book.ShopID == session.ShopID {
			j = true
			c.HTML(http.StatusOK, "updatebook.html", gin.H{
				"j":    j,
				"book": book,
			})
		}
	} else {
		//添加
		c.HTML(http.StatusOK, "updatebook.html", gin.H{
			"j": j,
		})
	}
}

//添加/修改书籍表单
func PostBook(c *gin.Context) {
	file, err := c.FormFile("f")
	var filename string
	bookid := c.Query("bookid")
	if err != nil {
		if bookid == "" { //添加
			filename = "默认图片.jpeg"
		} else { //修改
			filename = c.PostForm("f1")
		}

	} else {
		filename = file.Filename
		osfile, err2 := os.Create("../static/书籍图片/" + filename)
		if err2 != nil {
			fmt.Println("err2", err2)
		}
		defer osfile.Close()                    //关闭文件句柄
		c.SaveUploadedFile(file, osfile.Name()) // 上传文件到指定的目录
	}
	book := &modal.Book{}
	if err := c.ShouldBind(book); err == nil {
		book.Imgpath = filename
		if bookid == "" { //添加
			sess, _ := c.Get("sess")
			session := sess.(*modal.Session)
			book.ShopID = session.ShopID
			book.ShopName = session.ShopName
			book.Sales = 0
			book.Add()
		} else { //修改
			ibookid, _ := strconv.ParseInt(bookid, 10, 64)
			book.ID = ibookid
			book.Update()
		}
		MyBookShop(c)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
