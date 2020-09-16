package processor

import (
	//"fmt"
	"gin-bookmall/modal"
	"gin-bookmall/tool"
	"net/http"

	"github.com/gin-gonic/gin"
)

//AddBookCar 添加图书到购物车
func PostCar(c *gin.Context) {
	book := &modal.Book{}
	if err := c.ShouldBind(book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := book.Query() //根据bookID查询book
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//获取用户session
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	if session.Car == nil { //创建一个购物车
		carid := tool.UniqueID()
		car := &modal.Car{
			ID:     carid,
			UserID: session.UserID,
			ShopID: book.ShopID,
		}
		cartItm := &modal.Cartitm{ //创建购物项
			CarID:    carid,
			BookID:   book.ID,
			Count:    1,
			BookName: book.Title,
			Price:    book.Price,
			Imgpath:  book.Imgpath,
		}
		cartItm.Amount = cartItm.GetAmout()
		car.Add() //将购物车保存到数据库
		cartItm.Add()
		session.Car = car
		session.Add()
	} else {
		cartitm := &modal.Cartitm{
			CarID:  session.Car.ID,
			BookID: book.ID,
		}
		err := cartitm.Query()
		if err != nil { //创建一个购物项
			cartitm.Count = 1
			cartitm.BookName = book.Title
			cartitm.Price = book.Price
			cartitm.Amount = cartitm.GetAmout()
			cartitm.Imgpath = book.Imgpath
			cartitm.Add() //将购物项加入到数据库中
		} else {
			cartitm.Count = cartitm.Count + 1
			cartitm.Amount = cartitm.GetAmout()
			cartitm.Update() //更新数据库
		}
	}
	c.JSON(http.StatusOK, book.Title)
}

//GetCar 获取购物车信息
func GetCar(c *gin.Context) {
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	car := modal.GetCar()
	car.Query(session.UserID) //通过用户名查购物车
	_, cartitms := modal.Querys(car.ID)
	var judge bool
	if len(cartitms) != 0 {
		car.CartItms = cartitms
		car.Totalcount = car.GetTotalCount()
		car.Totalamount = car.GetTotalAmount()
		judge = true
	}
	c.HTML(http.StatusOK, "car.html", gin.H{
		"judge": judge,
		"car":   car,
	})
}

//DeleteIDCar删除购物车中的购物项/清空购物车
func DeleteCar(c *gin.Context) {
	sess, _ := c.Get("sess")
	session := sess.(*modal.Session)
	bookid := c.Query("bookid")
	if bookid != "" {
		modal.GetCartitm().Delete(session.Car.ID, bookid)
	} else {
		modal.GetCar().Delete(session.Car.ID, session.UserID)
		session.Car = nil
		session.Add()
	}
	c.JSON(http.StatusOK, "成功")
}

//UpdateCartItmID 根据购物项的数量更新购物车
func PutCar(c *gin.Context) {
	cartitm := modal.GetCartitm()
	if err := c.ShouldBind(cartitm); err != nil {
		return
	}
	if cartitm.Count > 0 {
		i := &modal.Inventorie{
			ID: cartitm.BookID,
		}
		if err := i.Query(); err == nil && uint16(cartitm.Count) < i.Stock {
			cartitm.Amount = cartitm.GetAmout()
			//fmt.Println("cart amount ",)
			cartitm.Update()
			c.JSON(http.StatusOK, "成功")
			return
		}
	}
	c.JSON(http.StatusOK, "失败")
}
