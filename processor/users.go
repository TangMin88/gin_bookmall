package processor

import (
	"gin-bookmall/modal"
	"gin-bookmall/tool"
	"net/http"

	//"time"

	"github.com/gin-gonic/gin"
)

//Login 去登录页面
func Login(c *gin.Context) {
	//change, _ := c.Get("change")
	c.HTML(http.StatusOK, "login.html", nil)
}

//HomeLogin 提交登录表单
func PostLogin(c *gin.Context) {
	//判断cookie是否存在，避免重复建立cookie
	_, err := c.Cookie("bookmall")
	if err == nil {
		Page(c)
	} else {
		user1 := modal.GetUser()
		if err := c.ShouldBind(user1); err == nil {
			user2 := modal.GetUser()
			err = user2.QueryU(user1.Username)
			if err == nil && user1.Password == user2.Password {
				sess := &modal.Session{
					ID:       tool.UniqueID(),
					UserName: user2.Username,
					UserID:   user2.ID,
				}
				car := modal.GetCar() //判断是否有对应的购物车
				if err1 := car.Query(user2.ID); err1 == nil {
					sess.Car = car
				}
				shop := modal.GetShop()
				if err := shop.QueryU(user2.ID); err == nil {
					sess.ShopName = shop.ShopName
					sess.ShopID = shop.ID
				}
				sess.Add()
				//将cookie发送给浏览器,第一个参数为 cookie 名；第二个参数为 cookie 值；第三个参数为 cookie 有效时长，；
				//第四个参数为 cookie 所在的目录；第五个为所在域，表示我们的 cookie 作用范围；第六个表示是否只能通过 https 访问；
				//第七个表示 cookie 是否可以通过 js代码进行操作
				c.SetCookie("bookmall", sess.ID, 900, "/", "127.0.0.1", false, true)
				c.HTML(http.StatusOK, "daozhuan.html", user2.Username)
			} else {
				c.HTML(http.StatusOK, "login.html", "密码错误")
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}
}

//登出
func Logout(c *gin.Context) {
	cookie, _ := c.Request.Cookie("bookmall")
	//modal.GetSession(cookie.Value).Delete()
	c.SetCookie(cookie.Name, cookie.Value, -1, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	c.HTML(http.StatusOK, "daozhuan.html", nil)
}

//Regist 注册页面
func Regist(c *gin.Context) {
	c.HTML(http.StatusOK, "regist.html", nil)
}

//HomeRegist 提交注册表单
func PostRegist(c *gin.Context) {
	user := modal.GetUser()
	if err := c.ShouldBind(user); err == nil {
		user.Add()
		c.HTML(http.StatusOK, "login.html", "注册成功，请重新登录")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

//QueryName 通过Ajax请求验证用户名是否可用
func PostName(c *gin.Context) {
	user := modal.GetUser()
	if err := c.ShouldBind(user); err == nil {
		err = user.QueryU(user.Username)
		if err == nil {
			c.JSON(http.StatusOK, "用户名已存在")
		} else {
			c.JSON(http.StatusOK, "<font styles='color:green'>用户名可用</font>")
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}

//AcquireNumber 获取验证码
func PostAcquire(c *gin.Context) {
	number := c.PostForm("number")
	//fmt.Println(number)
	err := tool.Verification(number)
	if err == nil {
		//fmt.Println(code)
		c.JSON(http.StatusOK, "已发送")
	} else {
		c.JSON(http.StatusBadRequest, "发送失败")
	}

}

//对比验证码
func PostContrast(c *gin.Context) {
	number := c.PostForm("number")
	num := c.PostForm("num")
	num2, err := modal.Get(number)
	if err == nil {
		if num2 == num {
			c.JSON(http.StatusOK, "验证码正确")
		} else {
			c.JSON(http.StatusOK, "验证码错误")
		}
	} else {
		c.JSON(http.StatusOK, "请重新发送验证码")
	}

}

//Shouj 获取手机号
func PostShouj(c *gin.Context) {
	username := c.PostForm("username")
	user := modal.GetUser()
	err := user.QueryU(username)
	if err == nil {
		c.JSON(http.StatusOK, user.Number)
	} else {
		c.JSON(http.StatusOK, "用户不存在")
	}

}

//PassWordBack 更新密码
func PutPassWordBack(c *gin.Context) {
	user := modal.GetUser()
	if err := c.ShouldBind(user); err == nil {
		user.Update()
		c.HTML(http.StatusOK, "login.html", "密码更新成功，请重新登录")
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
