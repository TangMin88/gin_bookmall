package modal

import (
	"encoding/json"
	"gin-bookmall/dao"
	"strconv"

	"github.com/go-redis/redis"
)

type Book struct {
	ID        uint16      `json:"bookid,string" form:"bookid"`
	Title     string      `json:"title" form:"title"`        //书名
	Author    string      `json:"author" form:"author"`      //作者
	Price     float64     `json:"price,string" form:"price"` //价格
	Imgpath   string      `json:"imgpath"`                   //图片路径
	ShopID    uint16      `json:"shopid,string"`             //图书所属店铺id
	ShopName  string      `json:"shopname"`
	Inventory *Inventorie `gorm:"-" json:"-"`
}
type Inventorie struct {
	ID    uint16 `json:"bookid,string" form:"bookid"` //图书的id
	Sales uint16 `json:"sales,string" form:"sales"`   //销量
	Stock uint16 `json:"stock,string" form:"stock"`   //库存
}

//获取Book结构体
func GetBook(id uint16) *Book {
	book := &Book{
		ID: id,
	}

	return book
}

//AddBook 添加图书
func (b *Book) Add() error {
	err := dao.Db.Table("books").Create(b).Error
	if err != nil {
		return err
	}
	k := JointStr("books", strconv.FormatUint(uint64(b.ID), 10))
	err = SAdd(k, b)
	if err != nil {
		return err
	}
	b.SAdd()
	b.Inventory.ID = b.ID
	return b.Inventory.Add()

}

//shop  book集合 加 bookid
func (b *Book) SAdd() error {
	k := JointStr("shop", strconv.FormatUint(uint64(b.ShopID), 10), "book")
	st := redis.Z{ //库存
		Score:  float64(b.Inventory.Stock),
		Member: b.ID,
	}
	return dao.Rdb.ZAdd(k, st).Err()
}

//添加库存
func (i *Inventorie) Add() error {
	err := dao.Db.Create(i).Error
	if err != nil {
		return err
	}
	_, err = i.UpdateSt(float64(i.Stock))
	return err

}

//更新库存
func (i *Inventorie) Update() error {
	return dao.Db.Model(i).Updates(i).Error
}

//更新销量sales
func (i *Inventorie) UpdateSa(s float64) (float64, error) {
	return dao.Rdb.ZIncrBy("sales", s, strconv.FormatUint(uint64(i.ID), 10)).Result()
}

//更新库存stock
func (i *Inventorie) UpdateSt(s float64) (float64, error) {
	return dao.Rdb.ZIncrBy("stock", s, strconv.FormatUint(uint64(i.ID), 10)).Result()
}

//添加库存到redis
// func (i *Inventorie) RAdd() error {
// 	// sa := redis.Z{ //销量
// 	// 	Score:  float64(i.Sales),
// 	// 	Member: i.ID,
// 	// }
// 	st := redis.Z{ //库存
// 		Score:  float64(i.Stock),
// 		Member: i.ID,
// 	}
// 	pipel := dao.Rdb.Pipeline()
// 	pipel.ZAdd("sales", sa)
// 	pipel.ZAdd("stock", st)
// 	_, err := pipel.Exec()
// 	return err
// }

//QueryBook 根据书的id查询一本书
func (b *Book) Query() error {
	k := JointStr("books", strconv.FormatUint(uint64(b.ID), 10))
	str, err := dao.Rdb.Get(k).Result()
	if err != nil {
		err := dao.Db.Where("id=?", b.ID).First(b).Error
		if err == nil {
			k := JointStr("books", strconv.FormatUint(uint64(b.ID), 10))
			err = SAdd(k, b)
			return err
		}
		return err
	}
	return json.Unmarshal([]byte(str), b)
}

//查库存
func (i *Inventorie) Query() error {
	f1, err := dao.Rdb.ZScore("stock", strconv.FormatUint(uint64(i.ID), 10)).Result()
	if err != nil {
		err1 := dao.Db.Where("id = ?", i.ID).First(i).Error
		if err1 == nil {
			i.UpdateSt(float64(i.Stock))
			return nil
		}
		return err
	}
	i.Stock = uint16(f1)
	return nil
}

//根据图书的id更新图书
func (b *Book) Update() error {
	err := dao.Db.Model(b).Updates(b).Error
	if err != nil {
		return err
	}
	k := JointStr("books", strconv.FormatUint(uint64(b.ID), 10))
	return SAdd(k, b)
}

//DeleteBookid 根据图书的id删除对应的图书
func (b *Book) Delete() error {
	err := dao.Db.Where("id = ?", b.ID).Delete(Book{}).Error
	if err != nil {
		return err
	}
	err = dao.Db.Where("id = ?", b.ID).Delete(Inventorie{}).Error
	if err != nil {
		return err
	}
	return b.RDelete()
}

func (b *Book) RDelete() error {
	k := JointStr("books", strconv.FormatUint(uint64(b.ID), 10))
	pipli := dao.Rdb.Pipeline()
	pipli.Del(k).Err()
	pipli.ZRem("sales", strconv.FormatUint(uint64(b.ID), 10)).Err()
	pipli.ZRem("stock", strconv.FormatUint(uint64(b.ID), 10)).Err()
	ks := JointStr("shop", strconv.FormatUint(uint64(b.ShopID), 10), "book")
	pipli.ZRem(ks, strconv.FormatUint(uint64(b.ID), 10)).Err()
	_, err := pipli.Exec()
	return err
}
