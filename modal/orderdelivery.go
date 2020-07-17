package modal

type Orderdelevery struct {
	ID             string
	ReceiverName   string //收货人
	ReceiverMobile string //电话号码
	// ReceiverState    string //省份
	// ReceiverCity     string //城市
	// ReceiverDistrict string //区/县
	ReceiverAddress string //收货具体地址
	ReceiverZip     string //邮政编码
}
