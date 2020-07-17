package processor

import (
	"gin-bookmall/modal"
	"gin-bookmall/tool"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//AddBookCar 添加图书到购物车
func PostCar(c *gin.Context) {
	book := &modal.Book{}
	if err := c.ShouldBind(book); err == nil {
		err := book.Query() //根据bookID查询book
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//获取用户session
		sess, _ := c.Get("sess")
		session := sess.(*modal.Session)
		car := modal.GetCar() //判断是否有对应的购物车
		err = car.Query(session.UserID)
		if err != nil { //创建一个购物车
			carid := tool.UniqueID()
			car.ID = carid
			//car.UserID = userid.(int64)
			car.UserID = session.UserID
			cartItm := &modal.Cartitm{ //创建购物项
				CarID:  carid,
				Book:   book,
				BookID: book.ID,
				Count:  1,
			}
			cartItm.Amount = cartItm.GetAmout()
			car.Add() //将购物车保存到数据库
			cartItm.Add()
		} else {
			val, ok := car.CartItms[book.ShopID]
			if ok {
				for _, v := range val {
					if v.BookID == book.ID { //如果购物车中的购物项找到了对应的图书，则将图书数量加一后退出，否则到最后都没有进入if语句则创建购物项
						v.Count = v.Count + 1
						v.Amount = v.GetAmout()
						v.Update() //更新数据库
						c.JSON(http.StatusOK, book.Title)
						return
					}
				}
			}
			cartItm := &modal.Cartitm{ //创建一个购物项
				CarID:  car.ID,
				Book:   book,
				BookID: book.ID,
				Count:  1,
			}
			cartItm.Amount = cartItm.GetAmout()
			cartItm.Add() //将购物项加入到数据库中
		}
		c.JSON(http.StatusOK, book.Title)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

//GetCar 获取购物车信息
func GetCar(c *gin.Context) {
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	car := modal.GetCar()
	err := car.Query(session.UserID) //通过用户名查购物车
	var judge bool
	if err == nil {
		judge = true
	}
	c.HTML(http.StatusOK, "car.html", gin.H{
		"judge": judge,
		"car":   car,
	})
}

//DeleteIDCar删除购物车中的购物项/清空购物车
func DeleteCar(c *gin.Context) {
	carid := c.Query("carid")
	if carid != "" { //清空购物车
		modal.GetCar().Delete(carid)
	} else {
		cartitmid := c.Query("cartitmid")
		icatr, _ := strconv.ParseInt(cartitmid, 10, 64)
		modal.GetCartitm().Delete(icatr)
	}

	c.JSON(http.StatusOK, "成功")
}

//UpdateCartItmID 根据购物项的数量更新购物车
func PutCar(c *gin.Context) {
	cartitmid := c.Query("cartitmid")
	count := c.Query("count")
	icatr, _ := strconv.ParseInt(cartitmid, 10, 64)
	icount, _ := strconv.ParseInt(count, 10, 64)
	if icount > 0 {
		cartitm := modal.GetCartitm()
		cartitm.Query(icatr)
		if icount < cartitm.Book.Stock {
			cartitm.Count = icount
			cartitm.Amount = cartitm.GetAmout() //数量改变，需重新计算金额小计
			cartitm.Update()
		}
	}
	c.JSON(http.StatusOK, "成功")
}
