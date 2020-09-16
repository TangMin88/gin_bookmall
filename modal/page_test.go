package modal

import (
	"fmt"
	"testing"
)

func TestPage(t *testing.T) {
	fmt.Println("page中的函数")
	t.Run("QueryTotal", testQueryTotal)
}

func testQueryTotal(t *testing.T) {
	page := &Page{
		PageNo: 1,
	}
	err := page.QueryTotal()
	if err != nil {
		fmt.Println("test", err)
	}
	for _, v := range page.BooK {
		fmt.Println(v)
	}
}
