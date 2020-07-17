package hook

import (
	"gin-bookmall/modal"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StatCost 是一个统计耗时请求耗时的中间件
// func StatCost() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		start := time.Now()
// 		c.Next()
// 		cost := time.Since(start)
// 		fmt.Println(cost)
// 	}
// }

//LadingCookie 判断是不是登录状态
func WhetherLading() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("bookmall")
		if err == nil {
			sess := modal.GetSession(cookie)
			err = sess.Query()
			if sess.ID != "" {
				c.SetCookie("bookmall", cookie, 1000, "/", "127.0.0.1", false, true)
				c.Set("sess", sess) // 在请求上下文中设置值，后续的处理函数能够取到该值
				c.Next()            //调用后面的函数
				return
			}
		}
		c.Abort() //阻止调用后面的函数
		s := c.Request.Method
		if s == "GET" {
			c.HTML(http.StatusOK, "login.html", "请先登录")
		} else {
			c.JSON(http.StatusOK, "请先登录")
		}
	}
}

//判断是否有店铺
func WhetherShop() gin.HandlerFunc {
	return func(c *gin.Context) {
		//shopid, b := c.Get("shopid") //从上下文中取值（跨中间件取值）
		sess, _ := c.Get("sess")
		session := sess.(*modal.Session)
		if session.ShopID > 0 {
			c.Next() //调用后面的函数
		} else {
			c.Abort() //阻止调用后面的函数
			s := c.Request.Method
			if s == "GET" {
				c.HTML(http.StatusOK, "owner.html", "请先成为会员")
			} else {
				c.JSON(http.StatusOK, "请先成为会员")
			}
		}
	}
}
