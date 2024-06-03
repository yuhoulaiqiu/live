package core

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var rdb *redis.Client

func InitRedis(addr, pwd string, db int) (client *redis.Client) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd, // no password set
		DB:       db,  // use default DB
		PoolSize: 100,
	})
	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := rdb.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("连接redis数据库成功")
	}
	return rdb
}
