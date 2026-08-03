package cache

import (
	"context"
	"f_blog/backend/internal/config"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

var ctx = context.Background()

// Init 初始化 Redis 连接
func Init(cfg *config.RedisConfig) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connect failed: %v", err)
	}

	log.Println("Redis connected successfully")
	return nil
}

// Get 获取缓存
func Get(key string) (string, error) {
	return RDB.Get(ctx, key).Result()
}

// Set 设置缓存（带过期时间）
func Set(key string, value interface{}, expiration time.Duration) error {
	return RDB.Set(ctx, key, value, expiration).Err()
}

// Delete 删除缓存
func Delete(key string) error {
	return RDB.Del(ctx, key).Err()
}

// DeleteByPattern 按模式删除缓存（慎用，生产环境建议用精确 key）
func DeleteByPattern(pattern string) error {
	iter := RDB.Scan(ctx, 0, pattern, 100).Iterator()
	for iter.Next(ctx) {
		if err := RDB.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// Exists 检查 key 是否存在
func Exists(key string) bool {
	n, _ := RDB.Exists(ctx, key).Result()
	return n > 0
}

// Incr 原子递增
func Incr(key string) (int64, error) {
	return RDB.Incr(ctx, key).Result()
}

// Expire 设置过期时间
func Expire(key string, expiration time.Duration) error {
	return RDB.Expire(ctx, key, expiration).Err()
}
