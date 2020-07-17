package modal

import "gin-bookmall/dao"

//Home 发送到主页
type Page struct {
	BooK        []*Book //每页查询出的图书存放的切片
	PageNo      int64   `json:"pageNo" form:"pageNo"` //当前页
	PageSize    int64   //每页显示的条数
	TotalPage   int64   //总页数，通过计算得到
	TotalRecord int64   //总记录数，通过查询数据库得到
	Pricemin    float64 `json:"pricemin" form:"pricemin"`
	Pricemax    float64 `json:"pricemax" form:"pricemax"`
	State       int     `json:"state" form:"state"` //0 主页，1 价格区间主页
}

//设置每页显示的条数
var PageSize int64 = 4

//totalRecord(总记录数)，pageNo(当前页)
func (p *Page) GetPage() {
	p.PageSize = PageSize
	//总页数
	//var totalpage int64
	if (p.TotalRecord % p.PageSize) == 0 {
		p.TotalPage = p.TotalRecord / p.PageSize
	} else {
		p.TotalPage = p.TotalRecord/p.PageSize + 1
	}
	//p.TotalPage=totalpage
}

//判断是否有上一页
func (p *Page) IsHasPrev() bool {
	return p.PageNo > 1
}

//判断是否有下一页
func (p *Page) IsHasNext() bool {
	return p.PageNo < p.TotalPage
}

//获取上一页
func (p *Page) PreviousPage() int64 {
	return p.PageNo - 1
}

//获取下一页
func (p *Page) NextPage() int64 {
	return p.PageNo + 1
}

//QueryTotal获取带分页的图书
func (p *Page) QueryTotal() error {
	//获取图书的总记录数
	var totalRecord int64
	dao.Db.Model(&Book{}).Count(&totalRecord)
	p.TotalRecord = totalRecord
	//获取图书总页数，每页显示的条数
	p.GetPage()
	var books []*Book
	err := dao.Db.Limit(p.PageSize).Offset((p.PageNo - 1) * p.PageSize).Find(&books).Error
	p.BooK = books
	return err
}

//TotalBookPrice 获取价格带分页的图书
func (p *Page) QueryTotalP() error {
	//获取图书的总记录数
	var totalRecord int64
	dao.Db.Model(&Book{}).Where("price between ? and ?", p.Pricemin, p.Pricemax).Count(&totalRecord)
	p.TotalRecord = totalRecord
	p.GetPage()
	var books []*Book
	err := dao.Db.Where("price between ? and ?", p.Pricemin, p.Pricemax).Limit(p.PageSize).Offset((p.PageNo - 1) * p.PageSize).Find(&books).Error
	p.BooK = books
	return err
}

//属于店铺的所有图书
func (p *Page) QueryTotalS(shopid int64) error {
	//接收总记录数
	var totalRecord int64
	//获取图书的总记录数
	dao.Db.Model(&Book{}).Where("shop_id= ?", shopid).Count(&totalRecord)
	p.TotalRecord = totalRecord
	p.GetPage()
	var books []*Book
	err := dao.Db.Where("shop_id =?", shopid).Limit(p.PageSize).Offset((p.PageNo - 1) * p.PageSize).Find(&books).Error
	p.BooK = books
	return err
}
