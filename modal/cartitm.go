package modal

import (
	"gin-bookmall/dao"
)

//CartItm 购物车中的购物项
type Cartitm struct {
	ID     int64 `json:"cartitmid" form:"cartitmid"` //购物项的id
	Book   *Book `gorm:"-"`                          //购物项中的图书信息
	BookID int64
	Count  int64   `json:"count" form:"count"` //购物项中的图书数量
	Amount float64 //购物项中的金额小计，通过计算得到
	CarID  string  `json:"carid" form:"carid"` //当前购物项属于哪个购物车
}

//获取Cartitm结构体
func GetCartitm() *Cartitm {
	cartitm := &Cartitm{}
	return cartitm
}

//获取金额小计
func (cartItm *Cartitm) GetAmout() float64 {
	//获取当前购物项中图书的价格
	price := cartItm.Book.Price
	return float64(cartItm.Count) * price
}

//AddCartitm 添加购物车中的购物项
func (cartItm *Cartitm) Add() error {
	return dao.Db.Create(cartItm).Error
}

//根据id查询对应的购物项
func (cartItm *Cartitm) Query(id int64) error {
	err := dao.Db.Where("id =? ", id).First(cartItm).Error
	if err != nil {
		return err
	}
	//根据在数据库中的bookID查到对应的book，加入到购物项中
	book := GetBook(cartItm.BookID)
	err = book.Query()
	if err != nil {
		return err
	}
	cartItm.Book = book
	return nil
}

//GetCartitmCarID 根据car的id查询购物车对应的所有购物项
func Querys(carid string) (error, []*Cartitm) {
	var cartitms []*Cartitm
	err := dao.Db.Where("car_id=?", carid).Find(&cartitms).Error
	return err, cartitms
}

//UpdateCartitm 更新购物项
func (cartItm *Cartitm) Update() error {
	return dao.Db.Model(cartItm).Updates(cartItm).Error
}

//GetCartitmBookID 根据book的id与购物车的id查询对应的购物车是否有对应的购物项
// func (cartItm *Cartitm) Query() error {
// 	err := dao.Db.Where("book_id =? and car_id = ?", cartItm.BookID, cartItm.CarID).First(cartItm).Error
// 	if err != nil {
// 		return err
// 	}
// 	//根据在数据库中的bookID查到对应的book，加入到购物项中
// 	book := GetBook()
// 	book.ID = cartItm.BookID
// 	err = book.Query()
// 	if err != nil {
// 		return err
// 	}
// 	cartItm.Book = book
// 	return nil
// }

//DeleteCartItm 根据购物车id删除对应的所有购物项
func (cartItm *Cartitm) Deletes(CarID string) error {
	return dao.Db.Where("car_id = ?", CarID).Delete(Cartitm{}).Error
}

//DeleteIDCartItm 根据购物项id删除对应的购物项
func (cartItm *Cartitm) Delete(id int64) error {
	return dao.Db.Where("id=?", id).Delete(Cartitm{}).Error
}
