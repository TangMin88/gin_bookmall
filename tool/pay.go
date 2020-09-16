package tool

import (
	"fmt"

	"github.com/smartwalle/alipay"

	//"net/http"
	"os/exec"
	"strings"
	//"time"
	//"github.com/gin-gonic/gin"
)

var (
	// appId
	appId = "2021000120617450"
	// 应用公钥
	aliPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqgNudVzbLI+56fQKZyBk1lM5122ixSTTU7Y3bXhymnZXvlo/LZC2ypRm4JqmcKkcx0NsHHAgGhI153WHZmTEcIiQPkSO7+5K9OoKnrfnZoprS9fEgxzkvT9lbKqHSXEhcZfAYIZQw/Gefk1UokqsodksEI+4yQye3zHVuU2Ed3Bt9i68lzeYKtZd9O17ktrIQSrmCkAkHQB2JKJyLRAWcBHqUN+M9LCWt6GyWR9hPoox9VK3a95MlPENVhm1K4NGBgiqOLCjbetyuQhuLPLSBFs6Dw9BC4XP6prT0RTRXv9QS0EPyNHjg3Cz9vGU6Ebzoe9kSef6IcYV14cljbenbQIDAQAB"
	// 应用私钥
	privateKey = "MIIEowIBAAKCAQEAggaujcNRwnN8TkWxo/kd6IKU6rJy6iJjywK5Nh7j5F9h0dR0vp3+1XMfCuqLJYoKJFWEM6InD6KnyK0u/zMstO6crhVQ/reCIGk0BfaOyXsNglOVhg+RUnknd+B5NfiiPjvlEGeYyLiSmtJ2fQn4YarSP5vaywU1qQTCrQCtKQ/DVUx5U+IvoeRV06tcbTI+srmTF/Na147/WHA4Qw/93U834YSiAqh5ao/ghnRKBbnQlkU3Bx00NkCGj1YYn9iNEOatIv3AuLI9gMHx9bqg3DkH6hZos/UFtAe8TxC9au6izViDVeWbcrnc9CraRk6P6BVdHLi3TbmjSeSujg6ruQIDAQABAoIBAA/nYVJVoZp3Ja0tOR0lS1M1JaHPUd7xdeNO4fiMrrMwN3bC1cS67oCNJC7hoUNmLvdivljSbJStAHi7NhRg3gcDaIaidNWy/GeadpKEJdLfCjf2oUNUhCHiT8GJ40mGr9GM4eevxDBI4yWsuHFy1r6bdjvxhEFw0a9qtaUTgBqVdZ5qhQre1PJRLFRgHRXffpkuzWwE+EO49RBg75xnitqTJ57CJQ8wx1QQoBPyvtRulBFRlUe2GBHMV23s3W/YSJrYU3Q2OXHLHo1+2Ym05H4DFFXMuy/s/2l4oikpZb8od6ZJJ0/m7/SI+XAW55Y6zAoI+xf/n4vcBSCtH76vwEECgYEA2vk3WMH3qgSYPpKcFuxTDd0kXxiV+ymmnh/9Adb1tyualhEPjmioue2K4+HwRYrwY9ZvXawj49nc+4i807UW1eyMp3SFGRRONxRz5FfPVSwmhiDtBiNaJDnR8nVOcy5uHMkTKUfKKJjBIMyy3NV0nDlk2UTNmC/jpEdziXJVDj8CgYEAmAMq/SPunZf99FSXJR7fjxaPXxXFefj7hC4J85GvwstRbqfHONw9Xj9ytwDpnQhpFuiIewngSHG1hA4gUMZLaNzI9XDYSFgb91RfM6q+YxmaWigrZbzMd3+S7Qz/pGf+AbaTm66cODvLUppthnl2PRvg5OXiMafxC+O2HRECuAcCgYBJQL3HL3xOoCLeK+WTtZNDPAuC1JK35wMaOtFE4Ehq8mdQdHyjw1dOe6zO9zKN0SECBSZUS9Xlz2ghrWid2iK9hdi33D76WNShkHIEnWt6rr5keHdSalkpbT5SwfNwjMPBVXLXiHCUjCVvd4sOXUoZSQn3tRCiLMMWneCExn61uQKBgAjRljrdTMiDus3j4mja46lPa73ea3hqA11ltloVB5dLaEv9G8emr0C6eZM4UFU12bLkhpZsukA5qIgisak99737oQTsKP/5bJXqpSNAMo9ZOuUkE3BxhYMhOYrbCDGnfsrmpqWPeayhe2gtYVE91qgw59kfpQGwdoF0EmqZRAkdAoGBAKf4KwS1KdmtXwZPSH25L6SKR7PiXTpEPTCseHB6cBRnkoIzNq57boG9Advn03MJ50BLc93juu1O/Cm1oWcZzpvwarT3XdMUS7RwnR6+WAsANf3sZ5oxTmdurC6U9rX69Ehl6uvQmZScKvjNEHkibimntA4JN92W1/4RrLXq8BbE"
	Client     = alipay.New(appId, aliPublicKey, privateKey, false)
)

//func init() {
//client.LoadAliPayPublicKey("aliPublicKey")
// client.LoadAppPublicCert("应用公钥证书")
// client.LoadAliPayPublicCert("支付宝公钥证书")
// client.LoadAliPayRootCert("支付宝根证书")
//}

//网站扫码支付
func WebPageAlipay(s, o string) {
	fmt.Println("1111")
	pay := alipay.AliPayTradePagePay{}
	// 支付宝回调地址（需要在支付宝后台配置）
	// 支付成功后，支付宝会发送一个POST消息到该地址
	pay.NotifyURL = "http://192.168.2.104:8888/alipay"
	// 支付成功之后，浏览器将会重定向到该 URL
	pay.ReturnURL = "http://192.168.2.104:8888/return"
	//支付标题
	pay.Subject = "BookMall"
	//订单号，一个订单号只能支付一次
	pay.OutTradeNo = o
	//销售产品码，与支付宝签约的产品码名称,目前仅支持FAST_INSTANT_TRADE_PAY
	pay.ProductCode = "FAST_INSTANT_TRADE_PAY"
	//金额
	pay.TotalAmount = s
	url, err := Client.TradePagePay(pay)
	if err != nil {
		fmt.Println(err)
	}
	payURL := url.String()
	//这个 payURL 即是用于支付的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	//fmt.Println(payURL)

	//打开默认浏览器
	payURL = strings.Replace(payURL, "&", "^&", -1)
	exec.Command("cmd", "/c", "start", payURL).Start()
}

//手机客户端支付
func WapAlipay(s, o string) {
	fmt.Println("2222")
	pay := alipay.AliPayTradeWapPay{}
	pay.NotifyURL = "http://192.168.2.104:8888/alipay"
	// 支付成功之后，支付宝将会重定向到该 URL
	pay.ReturnURL = "http://192.168.2.104:8888/return"
	//支付标题
	pay.Subject = "BookMall"
	//订单号，一个订单号只能支付一次
	pay.OutTradeNo = o
	//商品code
	//pay.ProductCode = time.Now().String()
	//金额
	pay.TotalAmount = s
	//用户付款中途退出返回商户网站的地址
	pay.QuitURL = "http://location:8888/orders"
	url, err := Client.TradeWapPay(pay)
	if err != nil {
		fmt.Println(err)
	}
	payURL := url.String()
	//这个 payURL 即是用于支付的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	fmt.Println(payURL)
	//打开默认浏览器
	payURL = strings.Replace(payURL, "&", "^&", -1)
	exec.Command("cmd", "/c", "start", payURL).Start()
}

//交易关闭
func TradeClose(orderid string) {
	alipay := alipay.AliPayTradeClose{}
	alipay.OutTradeNo = orderid
	_, err := Client.TradeClose(alipay)
	if err != nil {
		fmt.Println("TradeClose", err)
	}
}
