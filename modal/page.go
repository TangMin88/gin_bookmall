package modal

import (
	"encoding/json"
	"fmt"
	"gin-bookmall/dao"
	"strings"

	"github.com/go-redis/redis"
)

//Home 发送到主页
type Page struct {
	BooK        []*Book //每页查询出的图书存放的切片
	PageNo      uint64  `json:"pageNo" form:"pageNo"` //当前页
	TotalPage   uint64  //总页数，通过计算得到
	TotalRecord uint64  //总记录数，通过查询数据库得到
	Pricemin    float64 `json:"pricemin" form:"pricemin"`
	Pricemax    float64 `json:"pricemax" form:"pricemax"`
	State       string  `json:"state" form:"state"` //0 主页，1 价格区间主页
}

//设置每页显示的条数
var PageSize uint64 = 8

//totalRecord(总记录数)，pageNo(当前页)
func (p *Page) GetPage() {
	//p.PageSize = PageSize
	//总页数
	//var totalpage int64
	if (p.TotalRecord % PageSize) == 0 {
		p.TotalPage = p.TotalRecord / PageSize
	} else {
		p.TotalPage = p.TotalRecord/PageSize + 1
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
func (p *Page) PreviousPage() uint64 {
	return p.PageNo - 1
}

//获取下一页
func (p *Page) NextPage() uint64 {
	return p.PageNo + 1
}

//刷新

//QueryTotal获取带分页的图书
func (p *Page) QueryTotal() error {
	start := int64((p.PageNo - 1) * PageSize)
	//获取图书的总记录数
	var t1 uint64
	t2 := uint64(dao.Rdb.ZCard("stock").Val())
	err := dao.Db.Model(&Book{}).Count(&t1).Error
	if err != nil {
		return err
	}
	var books []*Book
	if t2 != t1 {
		rows, err1 := dao.Db.Table("books").Order("id Asc").Limit(PageSize).Offset(start).Rows()
		if err1 != nil {
			return err1
		}
		for rows.Next() {
			book := &Book{}
			err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Imgpath, &book.ShopID, &book.ShopName)
			if err != nil {
				fmt.Println("querytotal err ", err)
			}
			i := &Inventorie{}
			err1 := dao.Db.Where("id = ?", book.ID).First(i).Error
			if err1 == nil {
				i.UpdateSt(float64(i.Stock))
			}
			book.Inventory = i
			books = append(books, book)
		}
	} else {
		sstr, err := dao.Rdb.ZRevRange("stock", start, start+7).Result()
		if err != nil {
			return err
		}
		pipel := dao.Rdb.Pipeline()
		for _, v := range sstr {
			k := JointStr("books", v)
			pipel.Get(k).Result()
		}
		res, _ := pipel.Exec()
		for _, v1 := range res {
			book := &Book{}
			s, e := v1.(*redis.StringCmd).Result()
			if e != nil {
				b := v1.(*redis.StringCmd).Args()
				ss := strings.Split(b[1].(string), ":")
				err := dao.Db.Where("id=?", ss[1]).First(book).Error
				if err == nil {
					SAdd(b[1].(string), book)
				}
			} else {
				json.Unmarshal([]byte(s), book)
			}
			i := &Inventorie{
				ID: book.ID,
			}
			i.Query()
			book.Inventory = i
			books = append(books, book)
		}
	}
	p.BooK = books
	p.TotalRecord = t1
	p.GetPage() //获取图书总页数
	return nil
}

//TotalBookPrice 获取价格带分页的图书
func (p *Page) QueryTotalP() error {
	//获取图书的总记录数
	var totalRecord uint64
	dao.Db.Model(&Book{}).Where("price between ? and ?", p.Pricemin, p.Pricemax).Count(&totalRecord)
	p.TotalRecord = totalRecord
	p.GetPage()
	var books []*Book
	err := dao.Db.Where("price between ? and ?", p.Pricemin, p.Pricemax).Limit(PageSize).Offset((p.PageNo - 1) * PageSize).Find(&books).Error
	p.BooK = books
	return err
}

//属于店铺的所有图书
func (p *Page) QueryTotalS(shopid string) error {
	start := int64((p.PageNo - 1) * PageSize)
	//获取图书的总记录数
	var t1 uint64
	k := JointStr("shop", shopid, "book")
	t2 := uint64(dao.Rdb.ZCard(k).Val())
	err := dao.Db.Model(&Book{}).Where("shop_id= ?", shopid).Count(&t1).Error
	if err != nil {
		return err
	}
	var books []*Book
	if t2 != t1 {
		rows, err := dao.Db.Table("books").Where("shop_id =?", shopid).Limit(PageSize).Offset(start).Rows()
		if err != nil {
			fmt.Println("rows err", err)
			return err
		}
		for rows.Next() {
			book := &Book{}
			err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.Imgpath, &book.ShopID, &book.ShopName)
			if err != nil {
				fmt.Println("querytotal err ", err)
			}
			i := &Inventorie{}
			err1 := dao.Db.Where("id = ?", book.ID).First(i).Error
			book.Inventory = i
			if err1 == nil {
				book.SAdd()
			}
			books = append(books, book)
		}
	} else {
		k := JointStr("shop", shopid, "book")
		sstr, err := dao.Rdb.ZRevRange(k, start, start+7).Result()
		if err != nil {
			return err
		}
		pipel := dao.Rdb.Pipeline()
		for _, v := range sstr {
			k := JointStr("books", v)
			pipel.Get(k).Result()
		}
		res, _ := pipel.Exec()
		for _, v1 := range res {
			book := &Book{}
			s, e := v1.(*redis.StringCmd).Result()
			if e != nil {
				b := v1.(*redis.StringCmd).Args()
				ss := strings.Split(b[1].(string), ":")
				err := dao.Db.Where("id=?", ss[1]).First(book).Error
				if err == nil {
					SAdd(b[1].(string), book)
				}
			} else {
				json.Unmarshal([]byte(s), book)
			}
			i := &Inventorie{}
			dao.Db.Where("id = ?", book.ID).First(i)
			book.Inventory = i
			books = append(books, book)
		}
	}
	p.TotalRecord = t1
	p.GetPage()
	p.BooK = books
	return nil
}
