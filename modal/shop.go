package modal

import "gin-bookmall/dao"

type Shop struct {
	ID       int64   `json:"shopid" form:"shopid"`
	ShopName string  `json:"shopname" form:"shopname"` //店名
	Books    []*Book //属于店铺的图书切片
	UserID   int64   `json:"userid" form:"userid"` //店铺所属用户id
}

//获取Shop结构体
func GetShop() *Shop {
	shop := &Shop{}
	return shop
}

//AddShop 添加店铺
func (shop *Shop) Add() error {
	return dao.Db.Create(shop).Error
}

//DeleteShop 删除店铺
func (shop *Shop) Delete(shopid int64) error {
	return dao.Db.Where("id=?", shopid).Delete(Shop{}).Error
}

//QueryShop 查询店铺
func (shop *Shop) QueryU(userid int64) error {
	return dao.Db.Where("user_id=?", userid).First(shop).Error
}

//QueryShopID 根据店铺id查询店铺
func (shop *Shop) QueryS(shopid int64) error {
	return dao.Db.Where("id=?", shopid).First(shop).Error
}

//UpdateShopName 更新店铺
func (shop *Shop) Update() error {
	return dao.Db.Model(shop).Updates(shop).Error
}
