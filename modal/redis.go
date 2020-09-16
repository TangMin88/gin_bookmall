package modal

import (
	"encoding/json"
	"gin-bookmall/dao"
	"time"
	//"fmt"
)

//
func Set(k, v string) error {
	return dao.Rdb.Set(k, v, 30*time.Second).Err()
}

func Get(k string) (string, error) {
	return dao.Rdb.Get(k).Result()
}

//将图书添加/更新到redis string
func SAdd(k string, v interface{}) error {
	str, _ := json.Marshal(v)
	return dao.Rdb.Set(k, str, 0).Err()
}
