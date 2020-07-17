package modal

import (
	//"fmt"
	"fmt"
	"testing"
)

func TestCar(t *testing.T) {
	fmt.Println("car中的函数")
	//t.Run("testAddCar", testAddCar)
	//t.Run("testGetCarUserID", testGetCarUserID)

	t.Run("testDeleteCar", testDeleteCar)
}

func testAddCar(t *testing.T) {
	book := GetBook(5)
	err := book.Query()
	if err != nil {
		fmt.Println("1", err)
		return
	}
	cartitm := &Cartitm{
		Book:   book,
		BookID: book.ID,
		Count:  2,
		CarID:  "59595959",
	}
	cartitm.Amount = cartitm.GetAmout()
	car := &Car{
		ID:     "59595959",
		UserID: 3,
	}
	car.GetMap(cartitm)
	car.Totalcount = car.GetTotalCount()
	car.Totalamount = car.GetTotalAmount()
	err = car.Add()
	if err != nil {
		fmt.Println("2", err)
		return
	}
	err = cartitm.Add()
	if err != nil {
		fmt.Println("3", err)
		return
	}
}

func testGetCarUserID(t *testing.T) {
	car := GetCar()
	err := car.Query(3)
	if err != nil {
		return
	}
	fmt.Println(car)
	for _, v := range car.CartItms {
		for _, v1 := range v {
			fmt.Println(v1)
			fmt.Println(v1.Book)
		}

	}
}

func testDeleteCar(t *testing.T) {
	car := GetCar()
	car.Delete("59595959")
}