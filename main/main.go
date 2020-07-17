package main

import (
	"gin-bookmall/hook"
	"gin-bookmall/processor"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default() //创建一个默认的路由引擎

	r.Static("/static", "../static") //处理静态资源
	r.LoadHTMLGlob("../htm/*")       //处理模板文件
	r.GET("/page", processor.Page)   //首页
	r.POST("/page", processor.Page)  //带价格的首页

	user := r.Group("/user")
	{
		user.GET("/login", processor.Login)      //去登录页面
		user.POST("/login", processor.PostLogin) //登录表单提交
		user.GET("/logout", processor.Logout)    //登出

		user.GET("/regist", processor.Regist)      //去注册页面
		user.POST("/regist", processor.PostRegist) //注册表单提交

		user.POST("/username", processor.PostName)           //通过Ajax请求验证用户名是否可用
		user.POST("/acquirenumber", processor.PostAcquire)   //获取验证码
		user.POST("/shouj", processor.PostShouj)             //获取手机号
		user.PUT("/passwordback", processor.PutPassWordBack) //更新密码
	}

	car := r.Group("/car")
	car.Use(hook.WhetherLading())
	{
		car.GET("", processor.GetCar)       //获取购物车信息
		car.POST("", processor.PostCar)     //添加图书到购物车
		car.DELETE("", processor.DeleteCar) //清空购物车/删除购物车中的购物项
		car.PUT("", processor.PutCar)       //更新购物车中购物项的数量
		carorder := car.Group("/order")
		{
			carorder.GET("", processor.GetToCheckOut) //去结账
			carorder.POST("", processor.PostTheOrder) //确认订单页面，根据是否付款，更新订单
		}

	}

	order := r.Group("/orders", hook.WhetherLading())
	{
		order.GET("", processor.GetOrder)       //我的订单，查询不同状态的订单
		order.PUT("", processor.PutOrder)       //我的订单，更新订单状态
		order.DELETE("", processor.DeleteOrder) //我的订单，取消订单
	}

	r.GET("/owner", hook.WhetherLading(), processor.AsTheOwner)      //成为店主
	r.POST("/owner", hook.WhetherLading(), processor.PostAsTheOwner) //提交成为店主的表单

	shop := r.Group("/shop")
	shop.Use(hook.WhetherLading(), hook.WhetherShop())
	{
		shop.GET("", processor.MyBookShop) //我的店铺
		book := shop.Group("/book")
		{
			book.DELETE("", processor.DleteBook) //删除书籍
			book.GET("", processor.GetBook)      //去添加/修改书籍页面
			book.POST("", processor.PostBook)    //提交添加/修改书籍表单
		}
		invoicep := shop.Group("/invoicep")
		{
			invoicep.GET("", processor.GetInvoice)     //我的货单
			invoicep.PUT("", processor.PutTheDelivery) //发货
		}
	}

	r.Run(":8888")
}
