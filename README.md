# BookMall
一款go语言编写的RESTful风格的书籍买卖服务端
## 内容
用户管理，购物车管理，订单管理，后台管理，价格区间查询书籍等
## 安装
go get -u github.com/TangMin88/gin_bookmall
## 环境依赖
* github.com/aliyun/alibaba-cloud-sdk-go v1.61.226  //短信验证
* github.com/gin-gonic/gin v1.6.3                   //web框架
* github.com/go-redis/redis v6.15.9+incompatible    //缓存
* github.com/go-sql-driver/mysql v1.5.0 // indirect 
* github.com/jinzhu/gorm v1.9.12                    //go orm
* github.com/onsi/ginkgo v1.14.1 // indirect
* github.com/onsi/gomega v1.10.2 // indirect
* github.com/smartwalle/alipay v1.0.2               //go 支付接口
## 版本控制
该项目使用go module(Go语言默认的依赖管理工具)进行版本管理,您可以在go mod文件参看当前可用版本.
