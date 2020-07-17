package modal

import (
	"gin-bookmall/dao"
	//"ginbookmall/dao"
)

type Book struct {
	ID       int64
	Title    string  `json:"title" form:"title"`   //书名
	Author   string  `json:"author" form:"author"` //作者
	Price    float64 `json:"price" form:"price"`   //价格
	Sales    int64   //销量
	Stock    int64   `json:"stock" form:"stock"` //库存
	Imgpath  string  //图片路径
	ShopID   int64   //图书所属店铺id
	ShopName string
}

//获取Book结构体
func GetBook(id int64) *Book {
	book := &Book{
		ID: id,
	}
	return book
}

//AddBook 添加图书
func (b *Book) Add() error {
	return dao.Db.Create(b).Error
}

//QueryBook 根据书的id查询一本书
func (b *Book) Query() error {
	return dao.Db.Where("id=?", b.ID).First(b).Error
}

//根据图书的id更新图书
func (b *Book) Update() error {
	return dao.Db.Model(b).Updates(b).Error
}

//DeleteBookid 根据图书的id删除对应的图书
func (b *Book) Delete() error {
	return dao.Db.Where("id = ?", b.ID).Delete(Book{}).Error
}

//QueryBookShopID 根据店铺的id查询店铺中的所有图书
// func QueryBookShopID(shopID int64) ([]*Book, error) {
// 	var books []*Book
// 	err := dao.Db.Where("shopid=?", shopID).Find(&books).Error
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}
// 	return books, nil
// }
