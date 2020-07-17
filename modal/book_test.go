package modal

import (
	"fmt"
	"testing"
)

func TestBook(t *testing.T) {
	//fmt.Println("book中的函数")
	//t.Run("testAddBook", testAddBook)
	//t.Run("testQueryBook", testQueryBook)
	//t.Run("QueryBookShopID", testQueryBookShopID)
	//t.Run("testTotalBookk", testTotalBook)
	//t.Run("testUpdateBook", testUpdateBook)
	//t.Run("testDeleteBookid", testDeleteBookid)
	//t.Run("testShopTotalBook", testShopTotalBook)
}

func testAddBook(t *testing.T) {
	book := &Book{
		Title:   "汤姆索亚历险记",
		Author:  "马克·吐温",
		Price:   33.3,
		Sales:   0,
		Stock:   100,
		Imgpath: "static/书籍图片/默认图片.jpeg",
		ShopID:  3,
	}

	book.Add()
}

func testQueryBook(t *testing.T) {
	book := &Book{
		ID: 3,
	}
	book.Query()
	fmt.Println(book)
}

// func testQueryBookShopID(t *testing.T) {
// 	books, _ := QueryBookShopID(2)
// 	for _, v := range books {
// 		fmt.Println(v)
// 	}
// }

func testUpdateBook(t *testing.T) {
	book := &Book{
		ID:    41,
		Sales: 50,
	}
	book.Update()
}

func testDeleteBookid(t *testing.T) {
	err := GetBook(40).Delete()
	if err != nil {
		fmt.Println(err)
	}
}
