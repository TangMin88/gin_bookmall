package modal

import (
	"fmt"
	"testing"
)

func TestOrder(t *testing.T) {
	fmt.Println("order中的函数")
	//t.Run("testDetele", testDetele)
	//t.Run("testQueryU", testQueryU)
	t.Run("testUpdate", testUpdate)
}

func testDetele(t *testing.T) {
	err := GetOrder().Detele("08fa065f745af1585be84981aced4e11")
	if err != nil {
		fmt.Println(err)
	}
}

func testQueryU(t *testing.T) {
	orders, err := OrderQueryU(3, 1)
	if err != nil {
		fmt.Println(err)
		return

	}
	for _, v := range orders {
		fmt.Println(v)
	}

}

func testUpdate(t *testing.T) {

	order := GetOrder()
	order.ID = "0c4a2158cac1b6eebbcd7514ec46c26a"
	order.State = 2
	err := order.Update()
	if err != nil {
		fmt.Println(err)
	}
}
