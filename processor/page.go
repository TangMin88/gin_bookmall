package processor

import (
	//"fmt"
	"gin-bookmall/modal"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Page 主页
func Page(c *gin.Context) {
	page := &modal.Page{}
	if err := c.ShouldBind(&page); err == nil {
		if page.PageNo == 0 { //当前页页码
			page.PageNo = 1
		}
		switch page.State {
		case "1":
			page.QueryTotalP() //带价格查询的主页
		default:
			page.QueryTotal() //主页
		}
		cookie, err := c.Cookie("bookmall")
		var judge bool
		sess := modal.GetSession(cookie)
		if err == nil {
			err = sess.Query()
			if err == nil {
				judge = true
			}
		}
		c.HTML(http.StatusOK, "page.html", gin.H{
			"page":    page,
			"judge":   judge,
			"session": sess,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
