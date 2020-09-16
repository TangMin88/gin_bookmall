package modal

import (
	"fmt"
	"testing"
)

func TestBook(t *testing.T) {
	// fmt.Println("book中的函数")
	// t.Run("查询一本书", testQuery)
}

func testQuery(t *testing.T) {
	book:=GetBook(2)
	err:=book.Query()

	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println(book.ID)
	fmt.Println(book)
}

