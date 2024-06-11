package cache

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"time"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// Cache 提供缓存服务
type Cache struct {
	RedisClient *redis.Client
	DB          *gorm.DB
}

// NewCache  创建一个新的 CacheService
func NewCache(redisClient *redis.Client, db *gorm.DB) *Cache {
	return &Cache{
		RedisClient: redisClient,
		DB:          db,
	}
}

// GetOrSetCache 从缓存中获取数据，如果缓存中没有，则从数据库中获取并缓存
func (cs *Cache) GetOrSetCache(ctx context.Context, key string, fetchFromDB func() (interface{}, error), model interface{}, expire time.Duration) (interface{}, error) {
	// 尝试从 Redis 获取数据
	result, err := cs.RedisClient.Get(key).Result()
	if err == nil {
		// 创建一个新的指向正确类型的变量
		newModelPtr := reflect.New(reflect.TypeOf(model)).Interface()
		if err := json.Unmarshal([]byte(result), newModelPtr); err == nil {
			// 获取指针指向的值
			newModel := reflect.Indirect(reflect.ValueOf(newModelPtr)).Interface()
			return newModel, nil
		}
		logx.Error("从 Redis 解析数据失败:", err)
	} else if !errors.Is(err, redis.Nil) {
		logx.Error("从 Redis 获取数据失败:", err)
	}

	// 如果 Redis 中没有数据，则从数据库获取
	data, err := fetchFromDB()
	if err != nil {
		return nil, err
	}

	// 将数据存入 Redis
	dataBytes, err := json.Marshal(data)
	if err != nil {
		logx.Error("Failed to marshal data for Redis:", err)
		return data, nil
	}

	err = cs.RedisClient.Set(key, dataBytes, expire).Err()
	if err != nil {
		logx.Error("向 Redis 设置数据失败:", err)
	}

	return data, nil
}

// GetRandomExpire 获取一个随机的过期时间,防止缓存雪崩
func (cs *Cache) GetRandomExpire(times int) time.Duration {
	times = times + rand.Intn(10)
	return time.Duration(times)*time.Hour + time.Duration(rand.Intn(60))*time.Minute
}
