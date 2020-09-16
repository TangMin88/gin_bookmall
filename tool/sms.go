package tool

import (
	"fmt"
	"gin-bookmall/modal"
	"math/rand"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

//GenValidateCode 生成随机数
func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

//Verification 发送短信验证码
func Verification(num string) error {
	shu := 5
	code := GenValidateCode(shu) // 生成5位数字验证码功能
	// 下列的accessKeyId以及accessSecret请按实际申请到的填写
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "你的accessKeyId", "你的accessSecret")

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = num  //接收短信的手机号码。
	request.SignName = "number" //短信签名名称
	request.TemplateCode = "短信模板ID"
	request.TemplateParam = "{\"code\":\"" + code + "\"}" //短信模板变量对应的实际值，JSON格式

	_, err = client.SendSms(request)
	if err != nil {
		return err
	}
	//fmt.Printf("response is %#v\n", response)

	err = modal.Set(num, code)
	if err != nil {
		return err
	}
	return nil
}
