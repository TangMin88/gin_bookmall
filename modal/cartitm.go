package modal

import (
	"encoding/json"
	"fmt"
	"gin-bookmall/dao"
	"strconv"
)

//CartItm 购物车中的购物项
type Cartitm struct {
	ID       uint16  `json:"id,string" form:"cartitmid"` //购物项的id
	BookID   uint16  `json:"bookid,string" form:"bookid"`
	BookName string  `json:"bookname"`                  //书名
	Price    float64 `json:"price,string" form:"price"` //书 单价
	Count    uint    `json:"count,string" form:"count"` //购物项中的图书数量
	Amount   float64 `json:"amount,string"`             //购物项中的金额小计，通过计算得到
	CarID    string  `json:"carid" form:"carid"`        //当前购物项属于哪个购物车
	Imgpath  string  `json:"imgpath" form:"imgpath"`
}

//获取Cartitm结构体
func GetCartitm() *Cartitm {
	cartitm := &Cartitm{}
	return cartitm
}

//获取金额小计
func (cartItm *Cartitm) GetAmout() float64 {
	//获取当前购物项中图书的价格
	price := cartItm.Price
	return float64(cartItm.Count) * price
}

//AddCartitm 添加购物车中的购物项
func (cartItm *Cartitm) Add() error {
	err := dao.Db.Create(cartItm).Error
	if err != nil {
		return err
	}
	str, _ := json.Marshal(cartItm)
	return dao.Rdb.HSet(cartItm.CarID, strconv.FormatUint(uint64(cartItm.BookID), 10), str).Err()
}

//查找购物项是否存在
func (cartItm *Cartitm) Query() error {
	res, err := dao.Rdb.HGet(cartItm.CarID, strconv.FormatUint(uint64(cartItm.BookID), 10)).Result()
	if err != nil {
		err1 := dao.Db.Where("car_id =? and book_id = ? ", cartItm.CarID, cartItm.BookID).First(cartItm).Error
		if err1 != nil {
			return err1
		}

	}
	json.Unmarshal([]byte(res), cartItm)
	return nil
}

//GetCartitmCarID 根据car的id查询购物车对应的所有购物项
func Querys(carid string) (error, []*Cartitm) {
	var cartitms []*Cartitm
	err := dao.Rdb.HVals(carid).ScanSlice(cartitms)
	if err != nil {
		err := dao.Db.Table("cartitms").Where("car_id=?", carid).Find(&cartitms).Error
		if err != nil {
			return err, nil
		}
		pipel := dao.Rdb.Pipeline()
		for _, v := range cartitms {
			str, _ := json.Marshal(v)
			pipel.HSet(v.CarID, strconv.FormatUint(uint64(v.BookID), 10), str)
		}
		_, err = pipel.Exec()
		if err != nil {
			fmt.Println("pipel err", err)
		}
	}
	return nil, cartitms
}

//UpdateCartitm 更新购物项
func (cartItm *Cartitm) Update() error {
	err := dao.Db.Model(cartItm).Where("car_id=? and book_id =?", cartItm.CarID, cartItm.BookID).Updates(cartItm).Error
	if err != nil {
		return err
	}
	str, _ := json.Marshal(cartItm)
	return dao.Rdb.HSet(cartItm.CarID, strconv.FormatUint(uint64(cartItm.BookID), 10), str).Err()
}

//DeleteCartItm 根据购物车id删除对应的所有购物项
func (cartItm *Cartitm) Deletes(CarID string) error {
	err := dao.Db.Where("car_id = ?", CarID).Delete(Cartitm{}).Error
	if err != nil {
		return err
	}
	return dao.Rdb.Del(CarID).Err()
}

//DeleteIDCartItm 删除对应的购物项
func (cartItm *Cartitm) Delete(carid, bookid string) error {
	err := dao.Db.Where("car_id= ? and book_id = ?", carid, bookid).Delete(Cartitm{}).Error
	if err != nil {
		return err
	}
	return dao.Rdb.HDel(carid, bookid).Err()
}
