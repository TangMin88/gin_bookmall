package dao

import (
	"github.com/go-redis/redis"
	//"fmt"
)

//
var Rdb *redis.Client

//测试用，一台电脑上开启两个redis服务端
func init(){
	Rdb = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "master",
		SentinelAddrs: []string{":26379", ":26479"},
		DB:       0,
	})
	if Rdb == nil {
		panic("Redis NewClient failed")
	}
	err := Rdb.Ping().Err()
	if err != nil {
		panic("rds.NewClient error : " + err.Error())
	}
}